import {EditorState, TextSelection} from "prosemirror-state";
import {EditorView} from "prosemirror-view";
import {DOMParser} from "prosemirror-model";
import {Step} from "prosemirror-transform";
import {baseKeymap} from "prosemirror-commands";
import {keymap} from "prosemirror-keymap";
import {collab, sendableSteps, receiveTransaction, getVersion} from "prosemirror-collab";
import {history} from "prosemirror-history";
import {dropCursor} from "prosemirror-dropcursor";
import {gapCursor} from "prosemirror-gapcursor";
import {tableEditing, columnResizing} from "prosemirror-tables";

import {editorSchema} from "./schema.js";
import {menuCommands, inlineMenuCommands, tableMenuCommands} from "./commands.js";
import {tocPlugin} from "./plugins/toc.js";
import {keymapPlugin} from "./keymap.js";
import {inlineMenuPlugin} from "./plugins/inlinemenu.js";
import {uploadPlugin} from "./plugins/upload.js";
import {mentionsPlugin} from "./plugins/mentions.js";
import {menuPlugin} from "./plugins/menu.js";
import {tableMenuPlugin} from "./plugins/tablemenu.js";
import {multicursorPlugin} from "./plugins/multicursor.js";
import {placeholderPlugin} from "./plugins/placeholder.js";
import {embedPlugin} from "./plugins/embed.js";
import {markdownPlugin} from "./plugins/markdown.js";
import {CodeBlockView} from "./nodeviews/codeblock.js";
import {InfoboxView} from "./nodeviews/infobox.js";
import {ImageView} from "./nodeviews/image.js";
import {handleClickOn} from "./click.js";
import {debounce} from "../util";

const MAX_AUTHOR_COLORS = 9;
const PLACEHOLDER_TEXT = "Once upon a time...";

export class Editor {
    // Creates a new ProseMirror editor.
    // The callback functions are called in the following conditions:
    //
    // onUpdate is called every time state is updated, passing an object with the key being the object updated.
    // onSave is called when the article got successfully saved.
    // onLeave is called after the article was left on server and the client disconnected.
    // onDisconnect is called when the client disconnects (with or without error).
    // onError is called in case a connection error occurs.
    constructor(socket, articleId, langId, userId, editorDOM, content, onUpdate, onSave, onLeave, onDisconnect, onError) {
        // collaboration data and functions
        this.socket = socket;
        this.userId = userId;
        this.onUpdate = onUpdate;
        this.onSave = onSave;
        this.onLeave = onLeave;
        this.onDisconnect = onDisconnect;
        this.onError = onError;

        // collaboration
        this.authors = [];
        this.authorColorNumber = 1;

        // other fields
        this.cursorPos = 0;

        // plugins
        this.uploadPlugin = uploadPlugin(articleId, langId, onError);
        this.menuPlugin = menuPlugin(menuCommands(editorSchema, this));
        this.inlineMenuPlugin = inlineMenuPlugin(inlineMenuCommands(editorSchema));
        this.tableMenuPlugin = tableMenuPlugin(tableMenuCommands());
        this.tocPlugin = tocPlugin(onUpdate);
        this.mentionsPlugin = mentionsPlugin();
        this.multicursorPlugin = multicursorPlugin(this.updateCursor.bind(this));
        this.placeholderPlugin = placeholderPlugin(PLACEHOLDER_TEXT);
        this.embedPlugin = embedPlugin(editorSchema);
        this.markdownPlugin = markdownPlugin(editorSchema);

        // setup view using editor DOM element and given content
        this._setupView(editorDOM, content);

        // add debounced method to change the title
        this.sendUpdateTitle = debounce(function(title) {
            this.socket.emit("title_update", {title});
        }, 300);
    }

    // Establishes the websocket connection to the collaboration server and initializes all event functions.
    connect(articleId, langId, roomId) {
        this.socket.on("joined", this._joined.bind(this));
        this.socket.on("disconnect", this._disconnected.bind(this));
        this.socket.on("error", this._error.bind(this));
        this.socket.on("connection_err", this._connectionErr.bind(this));
        this.socket.on("state_updated", this._stateUpdated.bind(this));
        this.socket.on("state_changes", this._stateChanges.bind(this));
        this.socket.on("title_updated", this._titleUpdated.bind(this));
        this.socket.on("saved", this._saved.bind(this));
        this.socket.on("save_err", this._saveErr.bind(this));
        this.socket.on("tag_added", this._tagAdded.bind(this));
        this.socket.on("tag_removed", this._tagRemoved.bind(this));
        this.socket.on("lang_updated", this._langUpdated.bind(this));
        this.socket.on("access_set", this._accessSet.bind(this));
        this.socket.on("access_removed", this._accessRemoved.bind(this));
        this.socket.on("access_mode_updated", this._accessModeUpdated.bind(this));
        this.socket.on("client_access_updated", this._clientAccessUpdated.bind(this));
        this.socket.on("rtl_updated", this._rtlUpdated.bind(this));
        this.socket.on("author_connected", this._authorConnected.bind(this));
        this.socket.on("author_disconnected", this._authorDisconnected.bind(this));
        this.socket.on("cursor_updated", this._cursorUpdated.bind(this));
        this.socket.emit("open_article", {
            article_id: articleId,
            lang_id: langId,
            room_id: roomId
        });
    }

    destroy() {
        this.menuPlugin.destroy();
        this.inlineMenuPlugin.destroy();
        this.uploadPlugin.destroy();
        this.mentionsPlugin.destroy();
    }

    save(message, wip) {
        this.socket.emit("save", {message, wip});
    }

    closeArticle() {
        this.uploadPlugin.spec.fileUploader.cancelPendingUploads();
        this.socket.emit("close_article");
        this.socket.removeAllListeners();
        this.destroy();
    }

    leave() {
        this.socket.emit("leave");
        this.closeArticle();
        this.onLeave();
    }

    setTitle(title) {
        if(title !== this.title) {
            this.title = title;
            this.sendUpdateTitle(title);
        }
    }

    setLanguage(lang) {
        this.socket.emit("lang_update", {lang});
    }

    setAccessMode(mode) {
        this.socket.emit("access_mode_update", {mode});
    }

    setClientAccess(access) {
        this.socket.emit("client_access_update", {access});
    }

    addTag(name) {
        this.socket.emit("add_tag", {tag: name});
    }

    removeTag(name) {
        this.socket.emit("remove_tag", {tag: name});
    }

    setAccess(data) {
        let user = data.user || null;
        let group = data.group || null;
        let write = data.write !== undefined ? data.write : true;
        let req = {user, group, write};
        this.socket.emit("set_access", req);
    }

    removeAccess(data) {
        this.socket.emit("remove_access", data);
    }

    setRTL(rtl) {
        this.socket.emit("set_rtl", {rtl});
    }

    updateCursor(from, to) {
        this.socket.emit("update_cursor", {from, to});
        this.cursorPos = to;
    }

    getCursorPos() {
        return this.cursorPos;
    }

    setCursorPos(pos) {
        this.view.dispatch(this.view.state.tr.setSelection(TextSelection.near(this.view.state.doc.resolve(pos))));
        this.view.focus();
    }

    _setupView(editorDOM, content) {
        editorDOM.addEventListener("keydown", e => {
            if(e.keyCode === 9) { // tab
                e.preventDefault();
            }
        });

        let self = this;
        this.view = new EditorView(editorDOM, {
            state: EditorState.create({
                doc: DOMParser.fromSchema(editorSchema).parse(content),
                plugins: this._getPlugins()
            }),
            nodeViews: {
                code_block(node, view, getPos) {
                    return new CodeBlockView(node, view, getPos);
                },
                infobox(node, view, getPos) {
                    return new InfoboxView(node, view, getPos);
                },
                image(node, view, getPos) {
                    return new ImageView(node, view, getPos);
                }
            },
            dispatchTransaction(tx) {
                self._dispatchTransaction(tx);
            },
            handleClickOn
        });

        // disable browser build-in features to edit pictures and tables
        document.execCommand("enableObjectResizing", false, false);
        document.execCommand("enableInlineTableEditing", false, false);
    }

    _joined(data) {
        if(data.doc){
            let newState = EditorState.create({
                doc: editorSchema.nodeFromJSON(JSON.parse(data.doc)),
                plugins: this._getPlugins(data.version)
            });
            this.view.updateState(newState);
        }

        this.room = data.room_id;
        this.title = data.title;
        this.language = data.lang;
        this.tags = data.tags;
        this.accessMode = data.accessMode;
        this.access = data.access;
        this.clientAccess = data.clientAccess;
        this.rtl = data.rtl;
        this._getState();
        this.uploadPlugin.spec.fileUploader.setRoom(this.room);

        this.onUpdate({
            room: this.room,
            initialTitle: data.title,
            language: this.language,
            tags: this.tags,
            accessMode: this.accessMode,
            access: this.access,
            clientAccess: this.clientAccess,
            rtl: this.rtl
        });

        for(let i = 0; i < data.authors.length; i++) {
            this._authorConnected({user_id: data.authors[i]});
        }
    }

    _getPlugins(version) {
        let collabConfig;

        if(version) {
            collabConfig = {version};
        }

        return [
            this.mentionsPlugin,
            keymapPlugin(editorSchema),
            keymap(baseKeymap),
            history(),
            collab(collabConfig),
            dropCursor({width: 2, color: "#797C80"}),
            gapCursor(),
            this.menuPlugin,
            this.inlineMenuPlugin,
            this.tableMenuPlugin,
            this.uploadPlugin,
            this.tocPlugin,
            this.multicursorPlugin,
            this.placeholderPlugin,
            this.embedPlugin,
            this.markdownPlugin,
            columnResizing(),
            tableEditing()
        ];
    }

    _disconnected() {
        this.onDisconnect();
    }

    _error() {
        this.onError();
    }

    _connectionErr(data) {
        this.onError(data);
    }

    _saveErr(data) {
        this.onError(data);
    }

    _dispatchTransaction(tx) {
        let state = this.view.state.apply(tx);
        let sendable = sendableSteps(state);

        if(sendable) {
            this._sendState(sendable);
        }

        this.view.updateState(state);
    }

    _sendState(sendable) {
        this.socket.emit("state_update", {
            version: sendable.version,
            steps: sendable.steps,
            clientID: sendable.clientID
        });
    }

    _stateUpdated(data) {
        if(data && data.accepted && data.steps && data.steps.length) {
            let steps = data.steps.map(step => Step.fromJSON(editorSchema, step));
            let clientIDs = repeat(data.clientID, steps.length);
            this.view.dispatch(receiveTransaction(this.view.state, steps, clientIDs));
            this.onUpdate({stateUpdated: true});
        }
    }

    _getState() {
        this.socket.emit("get_state", {version: getVersion(this.view.state)});
    }

    _stateChanges(data) {
        if(data.title !== this.title) {
            this.onUpdate({
                title: data.title
            });
            this.title = data.title;
        }

        if(data.steps && data.steps.length) {
            data.steps = data.steps.map(step => Step.fromJSON(editorSchema, step));
            this.view.dispatch(receiveTransaction(this.view.state, data.steps, data.clientIDs));
        }
    }

    _titleUpdated(data) {
        if(data.user !== this.userId && this.title !== data.title) {
            this.title = data.title;

            this.onUpdate({
                title: data.title
            });
        }
    }

    _saved(data) {
        this.onSave(data.id, data.user_id, data.message, data.wip);
    }

    _tagAdded(data) {
        this.tags.push(data.tag);

        this.onUpdate({
            tags: this.tags
        });
    }

    _tagRemoved(data) {
        let tag = data.tag.toLowerCase();

        for(let i = 0; i < this.tags.length; i++){
            if(this.tags[i].toLowerCase() === tag){
                this.tags.splice(i, 1);
                break;
            }
        }

        this.onUpdate({
            tags: this.tags
        });
    }

    _langUpdated(data) {
        if(this.language !== data.lang) {
            this.language = data.lang;

            this.onUpdate({
                language: data.lang
            });
        }
    }

    _accessModeUpdated(data) {
        if(this.accessMode !== data.mode) {
            this.accessMode = data.mode;

            this.onUpdate({
                accessMode: data.mode
            });
        }
    }

    _clientAccessUpdated(data) {
        if(this.clientAccess !== data.access) {
            this.clientAccess = data.access;

            this.onUpdate({
                clientAccess: data.access
            });
        }
    }

    _rtlUpdated(data) {
        if(this.rtl !== data.rtl) {
            this.rtl = data.rtl;

            this.onUpdate({
                rtl: this.rtl
            });
        }
    }

    _accessSet(data) {
        let index = this._accessExists(data);

        if(index === -1){
            this.access.push(data);
        }
        else{
            this.access[index] = data;
        }

        this.onUpdate({
            access: this.access
        });
    }

    _accessRemoved(data) {
        let index = this._accessExists(data);

        if(index !== -1){
            this.access.splice(index, 1);
        }

        this.onUpdate({
            access: this.access
        });
    }

    _accessExists(access) {
        for(let i = 0; i < this.access.length; i++){
            if(access.user_id !== null && this.access[i].user_id === access.user_id ||
                access.user_group_id !== null && this.access[i].user_group_id === access.user_group_id){
                return i;
            }
        }

        return -1;
    }

    _authorConnected(data) {
        if(data.user_id === this.userId) {
            return;
        }

        let color = "color-"+(this.authorColorNumber%MAX_AUTHOR_COLORS+1);
        this.authorColorNumber++;
        this.authors.push({user_id: data.user_id, color});
        this.onUpdate({author_connected: data.user_id, color});

        // send cursor so other users can see it
        this.updateCursor(this.view.state.selection.from, this.view.state.selection.to);
    }

    _authorDisconnected(data) {
        for(let i = 0; i < this.authors.length; i++) {
            if(this.authors[i].user_id === data.user_id) {
                this.multicursorPlugin.removeCursor(data.user_id);
                this.authors.splice(i, 1);
                break;
            }
        }

        this.onUpdate({author_disconnected: data.user_id});
    }

    _cursorUpdated(data) {
        if(data.user_id !== this.userId) {
            let color = "color-1";

            for(let i = 0; i < this.authors.length; i++) {
                if(this.authors[i].user_id === data.user_id) {
                    color = this.authors[i].color;
                    break;
                }
            }

            this.multicursorPlugin.setCursor(data.user_id, color, data.from, data.to);
        }
    }
}

function repeat(value, n) {
    let array = [];

    for(let i = 0; i < n; i++) {
        array.push(value);
    }

    return array;
}
