const axios = require("axios");
const {Step} = require("prosemirror-transform");
const {Article} = require("./article.js");
const {schema} = require("./schema.js");
const {getRoomId, findRoomIdByArticleId} = require("./rooms.js");
const {logger} = require("./logger.js");

const BACKEND_HOST = process.env.EMVI_WIKI_BACKEND_HOST || "https://api.emvi.com";

// List of all active article "rooms".
// The room ID is used as key, while the value is an Article object.
let docs = new Map();

// List of all users connected.
let users = new Map();

// Accepts a new connection and makes sure the user is authorized.
module.exports.acceptConnection =  function(client) {
    let access_token = client.handshake.query.access_token;
    let organization = client.handshake.query.organization;
    let config = {
        headers: {
            "Authorization": "Bearer "+access_token,
            "Organization": organization
        }
    };

    return new Promise((resolve, reject) => {
        axios.get(`${BACKEND_HOST}/api/v1/user/member`, config)
        .then(r => {
            r.organization = organization;
            r.read_only = r.data.read_only;
            users.set(r.data.user_id, true);
            resolve(r);
        })
        .catch(e => {
            reject(e);
        });
    });
};

// Socket.io message handling.
module.exports.handleConnection = function(io, apiClient, client, organization, user_id, read_only) {
    logger.debug("Client connected", {
        organization,
        user_id,
        read_only
    });

    if(!user_id) {
        throw new Error("User ID not set");
    }

    client.user = {id: user_id, organization, read_only};

    // called when opening new/existing article
    client.on("open_article", data => {
        errorHandler(client, data, data => {
            // do not accept connections from read only members
            if(client.user.read_only) {
                throw new Error("Read only member");
            }

            determineRoomID(client, organization, data.article_id, data.lang_id, data.room_id);

            client.join(client.roomID, () => {
                joinOrCreateRoom(apiClient, client, data.article_id, data.lang_id)
                .then(() => {
                    io.to(client.roomID).emit("author_connected", {user_id});
                })
                .catch(e => {
                    if(e && e.message === "Client connected to room already") {
                        client.emit("connection_err", {connection_error: "connected_already"});
                    }
                    else if(e && e.message === "Maximum clients reached") {
                        client.emit("connection_err", {connection_error: "max_clients"});
                    }

                    client.emit("connection_err", {connection_error: "article_not_found"});
                    client.connectionError = e;
                    client.disconnect();
                });
            });
        });
    });

    // called when closing article
    client.on("close_article", data => {
        errorHandler(client, data, () => {
            io.to(client.roomID).emit("author_disconnected", {user_id});
            clientClosingArticle(client);
        });
    });

    client.on("state_update", data => {
        errorHandler(client, data, data => {
            let steps = data.steps.map(s => Step.fromJSON(schema, s));
            let accepted = docs.get(client.roomID).update(data.version, steps, data.clientID);

            io.to(client.roomID).emit("state_updated", {
                accepted,
                version: data.version,
                steps: data.steps,
                clientID: data.clientID
            });
        });
    });

    client.on("get_state", data => {
        errorHandler(client, data, data => {
            let state = docs.get(client.roomID).getState(data.version);
            client.emit("state_changes", {
                title: docs.get(client.roomID).getTitle(),
                steps: state.steps,
                clientIDs: state.clientIDs
            });
        });
    });

    client.on("title_update", data => {
        errorHandler(client, data, data => {
            docs.get(client.roomID).setTitle(data.title);
            data.user = client.user.id;
            client.to(client.roomID).emit("title_updated", data);
        });
    });

    // used to drop unsaved changes
    client.on("leave", data => {
        errorHandler(client, data, () => {
            logger.debug("Dropping changes");
            docs.get(client.roomID).setSaved(true);
        });
    });

    client.on("save", data => {
        errorHandler(client, data, data => {
            docs.get(client.roomID).save(client, data, true)
            .then(r => {
                if(r) {
                    io.to(client.roomID).emit("saved", {
                        id: r.data.id,
                        user_id: client.user.id,
                        message: data.message,
                        wip: data.wip
                    });
                }
            })
            .catch(e => {
                io.to(client.roomID).emit("save_err", e);
            });
        });
    });

    client.on("add_tag", data => {
        errorHandler(client, data, data => {
            docs.get(client.roomID).addTag(client, data.tag)
            .then(added => {
                if(added) {
                    io.to(client.roomID).emit("tag_added", data);
                }
            })
            .catch(() => {
                client.disconnect();
            });
        });
    });

    client.on("remove_tag", data => {
        errorHandler(client, data, data => {
            if (docs.get(client.roomID).removeTag(data.tag)) {
                io.to(client.roomID).emit("tag_removed", data);
            }
        });
    });

    client.on("lang_update", data => {
        errorHandler(client, data, data => {
            docs.get(client.roomID).setLanguage(data.lang);
            io.to(client.roomID).emit("lang_updated", data);
        });
    });

    client.on("set_access", data => {
        errorHandler(client, data, data => {
            docs.get(client.roomID).setAccess(data);
            io.to(client.roomID).emit("access_set", data);
        });
    });

    client.on("remove_access", data => {
        errorHandler(client, data, data => {
            if (docs.get(client.roomID).removeAccess(data)) {
                io.to(client.roomID).emit("access_removed", data);
            }
        });
    });

    client.on("access_mode_update", data => {
        errorHandler(client, data, data => {
            docs.get(client.roomID).setAccessMode(data.mode);
            io.to(client.roomID).emit("access_mode_updated", data);
        });
    });

    client.on("client_access_update", data => {
        errorHandler(client, data, data => {
            docs.get(client.roomID).setClientAccess(data.access);
            io.to(client.roomID).emit("client_access_updated", data);
        });
    });

    client.on("set_rtl", data => {
        errorHandler(client, data, data => {
            docs.get(client.roomID).setRTL(data.rtl);
            io.to(client.roomID).emit("rtl_updated", data);
        });
    });

    client.on("update_cursor", data => {
        errorHandler(client, data, data => {
            data.user_id = client.user.id;
            client.to(client.roomID).emit("cursor_updated", data);
        });
    });

    client.on("disconnect", reason => {
        try {
            logger.debug("Client disconnect", {reason});
            users.delete(client.user.id);

            // don't remove user from room on connection error (because he couldn't connect properly in the first place)
            // don't remove user from room if no article was opened
            if(!client.connectionError && docs.get(client.roomID)) {
                io.to(client.roomID).emit("author_disconnected", {user_id});
                clientClosingArticle(client);
            }
        }
        catch (e) {
            logger.debug(e.message);
        }
    });

    client.on("error", err => {
        logger.error("Connection error", {err});
    });
};

function determineRoomID(client, organization, article_id, lang_id, room_id) {
    if(article_id){
        client.roomID = findRoomIdByArticleId(docs, article_id, lang_id);
    }

    if(!client.roomID){
        client.roomID = getRoomId(organization, room_id);
    }
    else{
        logger.debug("Found room by article id", {article_id: article_id, room_id: client.roomID});
    }
}

function joinOrCreateRoom(apiClient, client, article_id, lang_id) {
    return new Promise((resolve, reject) => {
        try{
            if(!docs.get(client.roomID)){
                joinNewRoom(apiClient, client, article_id, lang_id)
                .then(() => {
                    resolve();
                })
                .catch(e => {
                    reject(e);
                });
            }
            else{
                joinExistingRoom(client);
                resolve();
            }
        }
        catch(e){
            reject(e);
        }
    });
}

function joinNewRoom(apiClient, client, article_id, lang_id) {
    return new Promise((resolve, reject) => {
        logger.info("Creating new room", {article_id: article_id, lang_id: lang_id, room_id: client.roomID});
        docs.set(client.roomID, new Article(apiClient));
        docs.get(client.roomID).addAuthor(client.user.id);
        docs.get(client.roomID).init(client)
        .then(() => {
            if (article_id) {
                // load and setup room
                docs.get(client.roomID).load(client, article_id, lang_id)
                .then(r => {
                    sendJoinedRoom(client, r.doc);
                    resolve();
                })
                .catch(e => {
                    reject(e);
                });
            } else {
                // setup an empty room for new article
                sendJoinedRoom(client);
                resolve();
            }
        })
        .catch(e => {
            reject(e);
        });
    });
}

function joinExistingRoom(client) {
    // prevent connecting twice to the same room
    if(!docs.get(client.roomID).addAuthor(client.user.id)){
        throw new Error("Client connected to room already");
    }

    // limit active connections to article depending on organization entry/expert
    if(docs.get(client.roomID).maxClientsReached()) {
        throw new Error("Maximum clients reached");
    }

    let doc = JSON.stringify(docs.get(client.roomID).getDoc());
    sendJoinedRoom(client, doc);
}

function sendJoinedRoom(client, doc) {
    let article = docs.get(client.roomID);
    client.emit("joined", {
        room_id: client.roomID,
        title: article.getTitle(),
        doc,
        tags: article.getTags(),
        lang: article.getLanguageID(),
        rtl: article.getRTL(),
        accessMode: article.getAccessMode(),
        access: article.getAccess(),
        clientAccess: article.getClientAccess(),
        version: article.getVersion(),
        authors: article.getAuthors()
    });
    logger.debug("Client joined room", {room_id: client.roomID});
}

function clientClosingArticle(client) {
    logger.debug("Client closing article");
    client.leaveAll();
    let room = docs.get(client.roomID);

    if(room) {
        room.removeAuthor(client.user.id);

        if (!room.hasAuthors()) {
            logger.debug("Last client closed article, closing room");

            room.save(client)
            .then(() => {
                docs.delete(client.roomID);
            })
            .catch(() => {
                docs.delete(client.roomID);
            });
        }
    }
}

// Saves all articles and closes connections.
module.exports.saveArticles = function(io) {
    for(let roomID of docs.keys()) {
        saveArticle(io, roomID);
    }
};

function saveArticle(io, roomID) {
    try {
        let room = docs.get(roomID);
        let authorId = room.getAuthors()[0]; // chose author who opened the room
        room.save({user: {id: authorId}, roomID});
        closeConnections(io, roomID);
    }
    catch(e) {
        logger.error("Error saving article on shutdown", {e});
        // continue with the next one...
    }
}

function closeConnections(io, roomID) {
    io.to(roomID).clients((error, socketIds) => {
        if(!error) {
            socketIds.forEach(socketId => {
                io.sockets.sockets[socketId].close();
            });
        }
    });
}

// Central error handling function which disconnects the client on exceptions.
function errorHandler(client, data, func) {
    try{
        func(data);
    }
    catch (e) {
        logger.error(e.message);
        client.disconnect();
    }
}

module.exports.getRoomCount = function() {
    return docs.size;
};

module.exports.getConnectionCount = function() {
    return users.size;
};
