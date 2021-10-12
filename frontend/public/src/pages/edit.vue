<template>
    <emvi-layout disable-events="true" hide-navigation="true">
        <emvi-layout-three-columns>
            <template slot="top">
                <emvi-connection-state :connection-state="connectionState"></emvi-connection-state>
                <emvi-authors :authors="authors"></emvi-authors>
            </template>
            <template slot="left">
                <emvi-toc :toc="toc"
                          v-on:jump="jumpToc"
                          v-show="connectionState === 1"></emvi-toc>
            </template>
            <template>
                <div :class="{rtl}" v-show="connectionState === 1">
                    <div class="article-headline">
                        <span class="labels">
                            <emvi-public-status v-show="accessMode === 0"></emvi-public-status>
                            <emvi-public-read-status v-show="accessMode === 1"></emvi-public-read-status>
                            <emvi-limited-status v-show="accessMode === 2"></emvi-limited-status>
                            <emvi-private-status v-show="accessMode === 3"></emvi-private-status>
                            <emvi-external-label v-if="clientAccess"></emvi-external-label>
                            <small class="label" v-if="language.name">{{language.name}}</small>
                        </span>
                    </div>
                    <emvi-article-title :title="title" v-on:change="setTitle"></emvi-article-title>
                    <emvi-tags :tags="tags"
                               :article-id="id"
                               read-only="true"
                               v-on:add="addTag"
                               v-on:remove="removeTag"></emvi-tags>
                    <div id="editor-wrapper" style="position: relative;white-space: pre-wrap;" ref="content">
                        <div ref="editor" v-on:click="contentClick"></div>
                    </div>
                    <div class="content-area" ref="editorContent" style="display: none;"></div>
                </div>
            </template>
        </emvi-layout-three-columns>
    </emvi-layout>
</template>

<script>
    import {mapGetters} from "vuex";
    import {UserService, LangService} from "../service";
    import {TocMixin} from "./toc";
    import {findIndexById, getIdFromSlug, slugWithId} from "../util";
    import {setPageTitle} from "./title";
    import {Editor} from "../editor";
    import {
        emviLayout,
        emviLayoutThreeColumns,
        emviArticleTitle,
        emviTags,
        emviToc,
        emviConnectionState,
        emviAuthors,
        emviExternalLabel,
        emviPublicStatus,
        emviPublicReadStatus,
        emviLimitedStatus,
        emviPrivateStatus
    } from "../components";

    export default {
        mixins: [TocMixin],
        components: {
            emviLayout,
            emviLayoutThreeColumns,
            emviArticleTitle,
            emviTags,
            emviToc,
            emviConnectionState,
            emviAuthors,
            emviExternalLabel,
            emviPublicStatus,
            emviPublicReadStatus,
            emviLimitedStatus,
            emviPrivateStatus
        },
        data() {
            return {
                dragoverHandler: null,
                dropHandler: null,
                focusedElement: null, // used to set focus after cmd is closed
                lastCursorPos: -1, // used to set the cursor after cmd is closed
                scrollY: 0, // used to set the scroll position after cmd is closed
                connectionState: 0, // 0 = not connected, 1 = connected, 2 = error
                id: "",
                langId: "",
                room: "",
                title: "",
                editor: null,
                tags: [],
                authors: [],
                rtl: false,
                accessMode: 0, // 0 = public r/w, 1 = public read, 2 = limited, 3 = private
                clientAccess: false,
                access: [],
                language: {},
                discarded: false
            };
        },
        computed: {
            ...mapGetters(["user", "metaUpdate", "cmdOpen"])
        },
        watch: {
            cmdOpen(open) {
                if(open) {
                    this.scrollY = window.scrollY;
                    this.focusedElement = document.activeElement;
                    let editorWrapper = document.querySelector(".ProseMirror");

                    if(this.focusedElement === editorWrapper) {
                        this.lastCursorPos = this.editor.getCursorPos();
                    }
                }
                else if(this.focusedElement) {
                    this.focusedElement.focus();

                    if(this.lastCursorPos > -1) {
                        this.editor.setCursorPos(this.lastCursorPos);
                        window.scrollTo(0, this.scrollY);
                        this.lastCursorPos = -1;
                    }
                }
            },
            metaUpdate() {
                this.updateFromMeta();
            }
        },
        mounted() {
            this.id = getIdFromSlug(this.$route.params.slug || "");
            this.langId = this.$route.query.lang || "";
            this.room = this.$route.query.room || "";
            this.setMeta();

            // the combination ID without language is not allowed, redirect back to read page in that case
            if(this.id && !this.langId) {
                this.$router.push(`/read/${this.$route.params.slug}`);
                return;
            }

            this.bindDragAndDrop();
            this.connect();
        },
        beforeRouteLeave(from, to, next) {
            this.editor.closeArticle();

            if(!this.discarded) {
                this.$store.dispatch("success", this.$t("toast_saved"));
            }
            else {
                this.$store.dispatch("success", this.$t("toast_discarded"));
            }

            next();
        },
        beforeDestroy() {
            this.unbindDragAndDrop();
        },
        methods: {
            bindDragAndDrop() {
                this.dragoverHandler = window.addEventListener("dragover", e => {
                    e.preventDefault();
                });
                this.dropHandler = window.addEventListener("drop", e => {
                    e.preventDefault();
                });
            },
            unbindDragAndDrop() {
                window.removeEventListener("dragover", this.dragoverHandler);
                window.removeEventListener("drop", this.dropHandler);
            },
            setMeta() {
                this.$store.dispatch("setMetaVars", [
                    {key: "page", value: this},
                    {key: "id", value: this.id},
                    {key: "langId", value: this.langId},
                    {key: "room", value: this.room}
                ]);
                this.$store.dispatch("setCmdMetaVars", [
                    {key: "id", value: this.id},
                    {key: "langId", value: this.langId},
                    {key: "room", value: this.room}
                ]);
            },
            connect() {
                this.editor = new Editor(this.$store.state.ws.socket,
                    this.id,
                    this.langId,
                    this.$store.state.user.user.id,
                    this.$refs.editor,
                    this.$refs.editorContent,
                    this.onUpdate,
                    this.onSave,
                    this.onLeave,
                    this.onDisconnect,
                    this.onError);
                this.editor.connect(this.id, this.langId, this.room);
            },
            onUpdate(update) {
                // called after user joined the room
                if(update.room) {
                    this.room = update.room;
                    this.connectionState = 1;
                    this.$store.dispatch("setMeta", {key: "room", value: update.room});
                    this.buildTableOfContents(this.getContent().innerHTML);
                }

                if(update.initialTitle !== undefined) {
                    this.title = update.initialTitle;
                    this.$store.dispatch("setMeta", {key: "title", value: this.title});
                    setPageTitle(this.title);
                }
                else if(update.title) {
                    this.title = update.title;
                    this.$store.dispatch("setMeta", {key: "title", value: this.title});
                    setPageTitle(this.title);
                }

                if(update.tags) {
                    this.tags = update.tags.slice(0);
                }

                if(update.author_connected) {
                    this.addAuthor(update.author_connected, update.color);
                }
                else if(update.author_disconnected) {
                    this.removeAuthor(update.author_disconnected);
                }

                if(update.rtl) {
                    this.rtl = update.rtl;
                    this.$store.dispatch("setMeta", {key: "rtl", value: update.rtl});
                }

                if(update.accessMode !== undefined && update.accessMode !== null) {
                    this.accessMode = parseInt(update.accessMode);
                    this.$store.dispatch("setMeta", {key: "accessMode", value: update.accessMode});
                }

                if(update.clientAccess !== undefined && update.clientAccess !== null) {
                    this.clientAccess = update.clientAccess;
                    this.$store.dispatch("setMeta", {key: "clientAccess", value: update.clientAccess});
                }

                if(update.access) {
                    this.access = update.access;
                    this.$store.dispatch("setMeta", {key: "access", value: update.access});
                }

                if(update.stateUpdated) {
                    this.buildTableOfContents(this.getContent().innerHTML);
                }

                if(update.language) {
                    LangService.getLang(update.language)
                        .then(lang => {
                            this.language = lang;
                        })
                        .catch(e => {
                            this.setError(e);
                        });
                }
            },
            onSave(id, user_id, message, wip) {
                this.id = id;

                if(user_id === this.user.id) {
                    this.onLeave();
                }
                else {
                    let index = findIndexById(this.authors, user_id);

                    if(index > -1) {
                        let author = this.authors[index];

                        if(wip) {
                            this.$store.dispatch("success", this.$t("toast_saved_by", {firstname: author.firstname, lastname: author.lastname, message}));
                        }
                        else {
                            this.$store.dispatch("success", this.$t("toast_published_by", {firstname: author.firstname, lastname: author.lastname, message}));
                        }
                    }
                }
            },
            onLeave() {
                if(!this.id) {
                    this.$router.push("/").catch(() => {});
                }
                else {
                    this.$router.push(`/read/${slugWithId(this.title, this.id)}?lang=${this.langId}`);
                }
            },
            onDisconnect() {
                this.connectionState = 2;
                this.$store.dispatch("closeCmd");
                this.$nextTick(() => {
                    this.$store.dispatch("resetCmd", this.$t("disconnect_cmd"));
                    this.$store.dispatch("pushColumn", "disconnect");
                });
            },
            onError(error) {
                if(error.connection_error === "article_not_found") {
                    this.$router.push(`/read/${slugWithId(this.title, this.id)}?lang=${this.langId}`);
                }
                else if(error.upload_error) {
                    if(error.upload_error === "file_size") {
                        this.$store.dispatch("error", this.$t("toast_upload_size"));
                    }
                    else {
                        this.$store.dispatch("error", this.$t("toast_upload"));
                    }
                }
                else if(error.connection_error === "connected_already") {
                    this.$store.dispatch("error", this.$t("toast_connected_already"));
                }
                else if(error.connection_error === "max_clients") {
                    this.$store.dispatch("error", this.$t("toast_max_clients"));
                }
                else {
                    console.error(error);
                }
            },
            contentClick(e) {
                // prevent following links
                if(e.button === 0) {
                    e.preventDefault();
                }
            },
            getContent() {
                return document.getElementsByClassName("ProseMirror")[0];
            },
            save(message, publish) {
                if(!publish) {
                    message = "Work in progress";
                }

                this.editor.save(message, !publish);
            },
            leave(discard) {
                this.discarded = !!discard;

                if(this.authors.length === 0) {
                    this.editor.leave();
                }
                else {
                    this.onLeave();
                }
            },
            updateFromMeta() {
                let rtl = this.$store.state.page.meta.get("rtl");
                let accessMode = this.$store.state.page.meta.get("accessMode");
                let clientAccess = this.$store.state.page.meta.get("clientAccess");
                let setAccess = this.$store.state.page.meta.get("setAccess");
                let removeAccess = this.$store.state.page.meta.get("removeAccess");
                let setLanguage = this.$store.state.page.meta.get("setLanguage");

                if(rtl !== undefined && rtl !== this.rtl) {
                    this.rtl = !this.rtl;
                    this.editor.setRTL(this.rtl);
                }

                if(accessMode !== undefined && parseInt(accessMode) !== this.accessMode) {
                    this.accessMode = parseInt(accessMode);
                    this.editor.setAccessMode(parseInt(accessMode));
                }

                if(clientAccess !== undefined && clientAccess !== this.clientAccess) {
                    this.clientAccess = clientAccess;
                    this.editor.setClientAccess(clientAccess);
                }

                if(setAccess !== undefined) {
                    for(let i = 0; i < setAccess.length; i++) {
                        this.editor.setAccess(setAccess[i]);
                    }

                    this.$store.dispatch("setMeta", {key: "setAccess", value: undefined});
                }

                if(removeAccess !== undefined) {
                    this.editor.removeAccess(removeAccess);
                    this.$store.dispatch("setMeta", {key: "removeAccess", value: undefined});
                }

                if(setLanguage !== undefined) {
                    this.editor.setLanguage(setLanguage);
                    this.$store.dispatch("setMeta", {key: "setLanguage", value: undefined});
                }
            },
            addAuthor(id, color) {
                if(id === this.user.id || findIndexById(this.authors, id) > -1) {
                    return;
                }

                UserService.getUser(id)
                    .then(user => {
                        user.color = color;
                        user.type = "user";
                        this.authors.push(user);
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            removeAuthor(id) {
                let index = findIndexById(this.authors, id);

                if(index !== -1) {
                    this.authors.splice(index, 1);
                }
            },
            setTitle(title) {
                this.editor.setTitle(title);
                this.$store.dispatch("setMeta", {key: "title", value: title});
                setPageTitle(title);
            },
            addTag(tag) {
                this.editor.addTag(tag);
            },
            removeTag(tag) {
                this.editor.removeTag(tag);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "toast_saved_by": "The article was saved by {firstname} {lastname}: {message}",
            "toast_published_by": "The article was published by {firstname} {lastname}: {message}",
            "toast_upload_size": "The file is too large. The limit is 50 MB.",
            "toast_upload": "Error uploading file. Please try again.",
            "toast_connected_already": "You're editing this article already. Do you opened it in a different tab or on another device?",
            "toast_max_clients": "The article is editing by the maximum number of authors right know. Please upgrade to an Expert organization to allow more concurrent authors.",
            "toast_saved": "Saved.",
            "toast_discarded": "Changes discarded.",
            "disconnect_cmd": "/disconnect"
        },
        "de": {
            "toast_saved_by": "Der Artikel wurde von {firstname} {lastname} gespeichert: {message}",
            "toast_published_by": "Der Artikel wurde von {firstname} {lastname} veröffentlicht: {message}",
            "toast_upload_size": "Die Datei ist für den Upload zu groß. Das Limit beträgt 50 MB.",
            "toast_upload": "Fehler beim Hochladen der Datei. Bitte versuche es erneut.",
            "toast_connected_already": "Du bearbeitest diesen Artikel bereits. Hast du ihn vielleicht in einem anderen Tab oder auf einem anderen Gerät geöffnet?",
            "toast_max_clients": "Der Artikel wird aktuell von der maximalen Anzahl Autoren bearbeitet. Bitte upgrade auf eine Expert Organisation um mehr gleichzeitige Autoren zu ermöglichen.",
            "toast_saved": "Gespeichert.",
            "toast_discarded": "Änderungen verworfen.",
            "disconnect_cmd": "/disconnect"
        }
    }
</i18n>
