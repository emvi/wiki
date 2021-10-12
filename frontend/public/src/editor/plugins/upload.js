import {Plugin, Selection} from "prosemirror-state";
import {Decoration, DecorationSet} from "prosemirror-view";
import {editorSchema} from "../schema.js";
import axios from "axios";
import {CancelToken} from "axios";
import {getSizeFromBytes} from "../../util";

const MAX_SIZE_BYTES = 50000000; // 50 MB

class FileUpload {
    constructor(articleId, langId, onError) {
        this.view = null;
        this.onError = onError;
        this.articleId = articleId;
        this.room = "";
        this.langId = langId;
        this.placeholderId = 0;
        this.placeholder = DecorationSet.empty;
        this.pendingUploads = [];
    }

    init(view) {
        if(!this.view) {
            this.view = view;
        }

        this.dropEventHandler = e => {
            e.stopPropagation();
            e.preventDefault();
            this.handleDrop(e);
            return false;
        };
        this.view.dom.addEventListener("drop", this.dropEventHandler);
    }

    destroy() {
        this.view.dom.removeEventListener("drop", this.dropEventHandler);
    }

    handleDrop(e) {
        if(e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files.length) {
            let transfer = e.dataTransfer.files;
            let pos = this.view.posAtCoords({left: e.clientX, top: e.clientY}).pos;

            // FIXME upload all files, not just the first one
            /*for (let i = 0; i < transfer.length; i++) {
                this.uploadFile(transfer[i], pos);
            }*/
            this.uploadFile(transfer[0], pos);
        }
    }

    uploadFile(file, pos) {
        if(file.size > MAX_SIZE_BYTES) {
            this.onError({upload_error: "file_size"});
            return;
        }

        this.placeholderId++;
        let id = this.placeholderId;
        this.createPlaceholder(pos, id, file.name, file.size, 0);
        this.sendFile(file, id);
    }

    sendFile(file, id) {
        let form = new FormData();
        form.append("file", file);

        let cancelToken = CancelToken.source();
        this.pendingUploads.push(cancelToken);
        let config = {
            headers: {
                "Content-Type": "multipart/form-data"
            },
            onUploadProgress: e => {
                let percentCompleted = Math.round((e.loaded*100)/e.total);
                this.updatePlaceholder(id, percentCompleted);
            },
            cancelToken: cancelToken.token
        };
        let articleOrRoom = this.articleId ? `?article=${this.articleId}` : `?room=${this.room}`;

        axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/content${articleOrRoom}&language=${this.langId}`, form, config)
        .then(r => {
            this.removePendingUpload(cancelToken);

            if(file.type.includes("image")) {
                this.insertImage(id, r.data.unique_name);
            }
            else if(file.type.includes("pdf")) {
                this.insertPDF(id, r.data.unique_name);
            }
            else{
                this.insertFile(id, r.data.unique_name, file.name, file.size);
            }

            this.removePlaceholder(id);
        })
        .catch(e => {
            console.error(e);
            this.removePendingUpload(cancelToken);
            this.removePlaceholder(id);

            if(e.errors && e.errors.length && e.errors[0].message && e.errors[0].message === "File too large") {
                this.onError({upload_error: "file_size"});
            }
            else {
                this.onError({upload_error: true});
            }
        });
    }

    createPlaceholder(pos, id, filename, filesize, progress) {
        let tr = this.view.state.tr;

        if(!pos) {
            pos = tr.selection.from;
        }

        let toDOM = () => {
            filesize = getSizeFromBytes(filesize);
            let dom = document.createElement("template");

            // upload placeholder HTML
            dom.innerHTML = `<span class="file upload no-select" contenteditable="false">\
<span class="progress no-select">\
<svg viewBox="0 0 36 36" class="file--progress--svg">\
<path class="stroke"\
stroke-dasharray="${progress}, 100"\
d="M18 2.0845\
a 15.9155 15.9155 0 0 1 0 31.831\
a 15.9155 15.9155 0 0 1 0 -31.831"></path>\
<path class="background"\
stroke-dasharray="100, 100"\
d="M18 2.0845\
a 15.9155 15.9155 0 0 1 0 31.831\
a 15.9155 15.9155 0 0 1 0 -31.831"></path>\
</svg>\
</span>\
<span class="name no-select">${filename}</span> \
<span class="size no-select">(${filesize})</span>\
</span>`;
            // end of placeholder HTML

            return dom.content.firstChild;
        };
        let decorator = Decoration.widget(pos, toDOM, {id, filename, filesize});
        this.placeholder = this.placeholder.add(tr.doc, [decorator]);
        this.view.dispatch(tr);
    }

    removePlaceholder(id) {
        let decorator = this.placeholder.find(null, null, spec => spec.id === id);
        this.placeholder = this.placeholder.remove(decorator);
        this.view.dispatch(this.view.state.tr);
    }

    updatePlaceholder(id, percentCompleted) {
        let decorator = this.placeholder.find(null, null, spec => spec.id === id);

        if(!decorator || !decorator.length) {
            return;
        }

        decorator = decorator[0];
        let pos = decorator.from;
        this.removePlaceholder(id);
        this.createPlaceholder(pos, id, decorator.spec.filename, decorator.spec.filesize, percentCompleted);
    }

    insertImage(id, uniqueName) {
        let pos = this.findPlaceholderPos(id);

        if(pos === null) {
            console.error("placeholder position not found");
            return;
        }

        let node = editorSchema.nodes.image.create({src: this.getFileURL(uniqueName)},
            editorSchema.nodes.paragraph.create());
        let resolvedPos = this.view.state.doc.resolve(pos);
        let tr = this.view.state.tr;
        tr.setSelection(Selection.near(resolvedPos));
        this.view.dispatch(tr.replaceSelectionWith(node));
    }

    insertPDF(id, uniqueName) {
        let pos = this.findPlaceholderPos(id);

        if(pos == null) {
            console.error("placeholder position not found");
            return;
        }

        let node = editorSchema.nodes.pdf.create({src: this.getFileURL(uniqueName)});
        this.view.dispatch(this.view.state.tr.insert(pos, node));
    }

    insertFile(id, uniqueName, name, size) {
        let pos = this.findPlaceholderPos(id);

        if(pos == null) {
            console.error("placeholder position not found");
            return;
        }

        let node = editorSchema.nodes.file.create({
            file: this.getFileURL(uniqueName),
            name,
            size: getSizeFromBytes(size)
        });
        this.view.dispatch(this.view.state.tr.insert(pos, node));
    }

    findPlaceholderPos(id) {
        let found = this.placeholder.find(null, null, spec => spec.id === id);
        return found.length ? found[0].from : null;
    }

    mapPlaceholder(tr) {
        this.placeholder = this.placeholder.map(tr.mapping, tr.doc);
    }

    setRoom(id) {
        this.room = id;
    }

    setPlaceholder(placeholder) {
        this.placeholder = placeholder;
    }

    getFileURL(uniqueName) {
        return `${EMVI_WIKI_BACKEND_HOST}/api/v1/content/${uniqueName}`;
    }

    getPlaceholder() {
        return this.placeholder;
    }

    removePendingUpload(cancelToken) {
        this.pendingUploads.splice(this.pendingUploads.indexOf(cancelToken), 1);
    }

    cancelPendingUploads() {
        for(let i = 0; i < this.pendingUploads.length; i++) {
            console.log("Cancelling upload...");
            this.pendingUploads[i].cancel();
        }
    }
}

export function uploadPlugin(articleId, langId, onError) {
    let fileUploader = new FileUpload(articleId, langId, onError);
    let plugin = new Plugin({
        view(view) {
            fileUploader.init(view);
            this.fileUploader = fileUploader;
            return fileUploader;
        },
        state: {
            init() {
                return fileUploader;
            },
            apply(tr, fileUploader) {
                if(tr.docChanged) {
                    fileUploader.mapPlaceholder(tr);
                }

                return fileUploader;
            }
        },
        props: {
            decorations(state) {
                let pluginState = this.getState(state);
                return pluginState.getPlaceholder();
            },
            handlePaste(view, event) {
                let items = event.clipboardData.items;

                for(let i = 0; i < items.length; i++) {
                    let item = items[i];

                    if(item.type.indexOf("image") !== -1) {
                        let file = item.getAsFile();
                        let pos = view.state.selection.$from.pos;
                        fileUploader.uploadFile(file, pos);
                        return true;
                    }
                }

                return false;
            }
        }
    });

    plugin.destroy = function () {
        fileUploader.destroy();
    };

    return plugin;
}
