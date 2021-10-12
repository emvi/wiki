<template>
    <emvi-cmd-selection-result :index="index" icon="history" :details="details" v-on:click="view" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <div class="item">
                <p>{{entity.user.firstname}} {{entity.user.lastname}}</p>
                <small>
                    {{entity.def_time | moment("LLL")}}
                    <span class="dot" v-if="entity.commit">·</span>
                    {{entity.commit}}
                </small>
            </div>
            <emvi-cmd-shortcut shortcut="Enter" v-show="active && !isActive" v-on:click="view">
                {{$t("shortcut_view")}}
            </emvi-cmd-shortcut>
            <emvi-cmd-shortcut shortcut="TAB" icon="restore" v-show="write && !entity.is_latest" v-on:click="showReset">
                {{$t("shortcut_reset")}}
            </emvi-cmd-shortcut>
            <emvi-cmd-shortcut :shortcut="$t('key_delete')" icon="trash" v-show="write" v-on:click="showDelete">
                {{$t("shortcut_remove")}}
            </emvi-cmd-shortcut>
            <span class="label" v-if="isActive">
                {{$t("active")}}
            </span>
        </template>
        <template slot="details">
            <div v-show="resetActive">
                <p>{{$t("confirmation_reset")}}</p>
                <emvi-cmd-input :label="$t('label_message')"
                    :index="index"
                    :selection="selection"
                    :selection-index="0"
                    :error="validation['message']"
                    v-model="message"
                    v-on:select="setSelection"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:enter="reset"
                    v-on:esc="cancel"></emvi-cmd-input>
                <emvi-cmd-selection-button icon="back"
                    :label="$t('label_no')"
                    :index="index"
                    :selection="selection"
                    :selection-index="1"
                    v-on:click="cancel"
                    v-on:select="setSelection"></emvi-cmd-selection-button>
                <emvi-cmd-selection-button icon="save"
                    color="green"
                    :label="$t('label_yes_reset')"
                    :index="index"
                    :selection="selection"
                    :selection-index="2"
                    v-on:click="reset"
                    v-on:select="setSelection"></emvi-cmd-selection-button>
            </div>
            <div v-show="removeActive">
                <p>{{$t("confirmation_remove")}}</p>
                <emvi-cmd-selection-button icon="back"
                    :label="$t('label_no')"
                    :index="index"
                    :selection="selection"
                    :selection-index="0"
                    v-on:click="cancel"
                    v-on:select="setSelection"></emvi-cmd-selection-button>
                <emvi-cmd-selection-button icon="trash"
                    color="red"
                    :label="$t('label_yes_remove')"
                    :index="index"
                    :selection="selection"
                    :selection-index="1"
                    v-on:click="remove"
                    v-on:select="setSelection"></emvi-cmd-selection-button>
            </div>
        </template>
    </emvi-cmd-selection-result>
</template>

<script>
    import {slugWithId} from "../../../../util";
    import {SelectionMixin} from "../mixin";
    import {ArticleService} from "../../../../service";
    import emviCmdSelectionResult from "../result.vue";
    import emviCmdEntity from "../../content/entity.vue";
    import emviCmdShortcut from "../../content/shortcut.vue";
    import emviCmdSelectionToggle from "../form/toggle.vue";
    import emviCmdSelectionButton from "../form/button.vue";
    import emviCmdInput from "../form/input.vue";

    export default {
        mixins: [SelectionMixin],
        components: {
            emviCmdSelectionResult,
            emviCmdEntity,
            emviCmdShortcut,
            emviCmdSelectionToggle,
            emviCmdSelectionButton,
            emviCmdInput
        },
        data() {
            return {
                resetActive: false,
                removeActive: false,
                maxSelectionIndex: 1,
                message: ""
            };
        },
        computed: {
            isActive() {
                let version = this.$store.state.page.meta.get("version");
                return this.entity.version === version || version === 0 && this.entity.is_latest;
            },
            write() {
                return this.active && this.$store.state.page.meta.get("write");
            }
        },
        watch: {
            selection(selection) {
                if(this.resetActive && selection !== 0) {
                    this.focusCmdInput();
                }
            },
            enter(enter) {
                if(enter && this.active) {
                    this.view();
                }
            },
            tab(tab) {
                if(tab && this.active && this.write && !this.removeActive && !this.entity.is_latest) {
                    this.showReset();
                }
            },
            del(del) {
                if(del && this.active && this.write && !this.resetActive) {
                    this.showDelete();
                }
            },
            esc(esc) {
                if(esc && this.details) {
                    this.cancel();
                }
            },
            details(details) {
                if(!details) {
                    this.resetActive = false;
                    this.removeActive = false;
                }
            }
        },
        methods: {
            view() {
                if(this.resetActive) {
                    if(this.selection === 1) {
                        this.cancel();
                    } else {
                        this.reset();
                    }
                }
                else if(this.removeActive) {
                    if(this.selection === 0) {
                        this.cancel();
                    } else {
                        this.remove();
                    }
                }
                else {
                    this.open();
                }
            },
            showReset() {
                this.maxSelectionIndex = 2;
                this.resetActive = !this.resetActive;
                this.removeActive = false;
                this.message = "";
                this.toggleDetails(this.resetActive);
            },
            showDelete() {
                this.maxSelectionIndex = 1;
                this.resetActive = false;
                this.removeActive = !this.removeActive;
                this.toggleDetails(this.removeActive);
            },
            open() {
                let activeVersion = this.$store.state.page.meta.get("version");
                let version = this.entity.version;

                if(version === activeVersion || activeVersion === 0 && this.entity.is_latest) {
                    return;
                }

                let id = this.$store.state.page.meta.get("id");
                let langId = this.$store.state.page.meta.get("langId");
                let title = this.$store.state.page.meta.get("title");
                this.$store.dispatch("resetCmd");

                if(this.entity.is_latest) {
                    this.$router.push(`/read/${slugWithId(title, id)}?lang=${langId}`);
                }
                else {
                    this.$router.push(`/read/${slugWithId(title, id)}?version=${version}&lang=${langId}`);
                }
            },
            reset() {
                this.resetError();

                if(this.entity.is_latest) {
                    return;
                }

                let id = this.$store.state.page.meta.get("id");
                let langId = this.$store.state.page.meta.get("langId");
                let title = this.$store.state.page.meta.get("title");

                ArticleService.resetArticle(id, langId, this.entity.version, this.message)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_reset"));
                        this.$store.dispatch("resetCmd");

                        if(this.$route.query.version) {
                            this.$router.push(`/read/${slugWithId(title, id)}?lang=${langId}`);
                        }
                        else {
                            this.$router.go();
                        }
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            remove() {
                this.resetError();
                ArticleService.deleteHistoryEntry(this.entity.id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_removed"));
                        this.$emit("remove", this.entity.id);
                        this.cancel();
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            cancel() {
                this.resetActive = false;
                this.removeActive = false;
                this.toggleDetails(false);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "active": "Active",
            "key_delete": "Del",
            "shortcut_view": "Open version",
            "shortcut_reset": "Reset version",
            "shortcut_remove": "Remove version",
            "confirmation_reset": "Are you sure you would like to reset to this version? This won't modify the history but add a new version on top.",
            "confirmation_remove": "Are you sure you would like to remove the version from the articles history?",
            "label_message": "Reason",
            "label_no": "No",
            "label_yes_reset": "Yes, reset to this version",
            "label_yes_remove": "Yes, remove version",
            "toast_reset": "The article version has been reset.",
            "toast_removed": "The articles history has been modified."
        },
        "de": {
            "active": "Aktiv",
            "key_delete": "Entf",
            "shortcut_view": "Version öffnen",
            "shortcut_reset": "Version zurücksetzen",
            "shortcut_remove": "Version entfernen",
            "confirmation_reset": "Möchtest du wirklich zu dieser Version zurücksetzen? Die Historie wird dadurch nicht verändert. Stattdessen wird eine neue Version hinzugefügt.",
            "confirmation_remove": "Möchtest du die Version wirklich aus der Artikelhistorie entfernen?",
            "label_message": "Grund",
            "label_no": "Nein",
            "label_yes_reset": "Ja, auf diese Version zurücksetzen",
            "label_yes_remove": "Ja, Version entfernen",
            "toast_reset": "Die Artikelversion wurde zurückgesetzt.",
            "toast_removed": "Die Artikelhistorie wurde angepasst."
        }
    }
</i18n>
