<template>
    <emvi-cmd-selection-result :index="index" icon="command" :details="details" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <div class="item">
                <p>{{entity.name}}</p>
            </div>
            <emvi-cmd-shortcut :shortcut="$t('key_delete')" icon="trash" v-show="active" v-on:click="showDelete">
                {{$t("shortcut_remove")}}
            </emvi-cmd-shortcut>
            <emvi-cmd-shortcut shortcut="Tab" icon="chevron" :rotate="detailsActive" v-show="active" v-on:click="showDetails">
                {{$t("shortcut_details")}}
            </emvi-cmd-shortcut>
        </template>
        <template slot="details">
            <div v-show="detailsActive">
                <emvi-cmd-input :label="$t('label_id')"
                    v-model="entity.client_id"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"></emvi-cmd-input>
                <emvi-cmd-input :label="$t('label_secret')"
                    v-model="entity.client_secret"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"></emvi-cmd-input>
                <emvi-cmd-checkbox :label="$t('label_scope_organization')"
                    name="scopeOrganization"
                    v-model="scopes.organization"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_language')"
                    name="scopeLanguage"
                    v-model="scopes.language"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_articles')"
                    name="scopeArticles"
                    v-model="scopes.articles"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_article_authors')"
                    name="scopeArticleAuthors"
                    v-model="scopes.article_authors"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_article_authors_mails')"
                    name="scopeArticleAuthorsMails"
                    v-model="scopes.article_authors_mails"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_article_history')"
                    name="scopeArticleHistory"
                    v-model="scopes.article_history"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_lists')"
                    name="scopeLists"
                    v-model="scopes.lists"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_tags')"
                    name="scopeTags"
                    v-model="scopes.tags"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_pinned')"
                    name="scopePinned"
                    v-model="scopes.pinned"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_search_all')"
                    name="scopeSearchAll"
                    v-model="scopes.search_all"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_search_articles')"
                    name="scopeSearchArticles"
                    v-model="scopes.search_articles"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_search_lists')"
                    name="scopeSearchLists"
                    v-model="scopes.search_lists"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
                <emvi-cmd-checkbox :label="$t('label_scope_search_tags')"
                    name="scopeSearchTags"
                    v-model="scopes.search_tags"
                    disabled="true"
                    container="cmd-selection-result-details"
                    v-on:next="nextRow"
                    v-on:previous="previousRow"
                    v-on:esc="cancel"></emvi-cmd-checkbox>
            </div>
            <div v-show="removeActive">
                <p>{{$t("confirmation")}}</p>
                <emvi-cmd-selection-button icon="back"
                    :label="$t('label_no')"
                    :index="index"
                    :selection="selection"
                    :selection-index="0"
                    v-on:click="cancel"
                    v-on:select="setSelection"></emvi-cmd-selection-button>
                <emvi-cmd-selection-button icon="trash"
                    color="red"
                    :label="$t('label_yes')"
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
    import {SelectionMixin} from "../mixin";
    import {ClientService} from "../../../../service";
    import emviCmdSelectionResult from "../result.vue";
    import emviCmdShortcut from "../../content/shortcut.vue";
    import emviCmdSelectionButton from "../form/button.vue";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import {scrollArea} from "../../../../util";

    export default {
        mixins: [SelectionMixin],
        components: {
            emviCmdSelectionResult,
            emviCmdShortcut,
            emviCmdSelectionButton,
            emviCmdInput,
            emviCmdCheckbox
        },
        data() {
            return {
                detailsActive: false,
                removeActive: false,
                maxSelectionIndex: 1,
                scopes: {}
            };
        },
        watch: {
            enter(enter) {
                if(enter && this.active && this.details && this.removeActive) {
                    if(this.selection === 0) {
                        this.cancel();
                    } else {
                        this.remove();
                    }
                }
            },
            tab(tab) {
                if(tab && this.active && !this.removeActive) {
                    this.showDetails();
                }
            },
            del(del) {
                if(del && this.active && !this.detailsActive) {
                    this.showDelete();
                }
            },
            esc(esc) {
                if(esc && this.details) {
                    this.cancel();
                }
            },
            up(up) {
                if(up && this.details) {
                    scrollArea(document.getElementById("cmd-selection-result-details"), -1);
                }
            },
            down(down) {
                if(down && this.details) {
                    scrollArea(document.getElementById("cmd-selection-result-details"), 1);
                }
            },
            details(details) {
                if(!details) {
                    this.detailsActive = false;
                    this.removeActive = false;
                }
            }
        },
        mounted() {
            this.buildScopes();
        },
        methods: {
            buildScopes() {
                let scopes = this.entity.scopes || [];

                for(let i = 0; i < scopes.length; i++) {
                    this.scopes[scopes[i].name] = scopes[i].read;
                }
            },
            showDetails() {
                this.maxSelectionIndex = 14;
                this.detailsActive = !this.detailsActive;
                this.removeActive = false;
                this.toggleDetails(this.detailsActive);
            },
            showDelete() {
                this.maxSelectionIndex = 1;
                this.detailsActive = false;
                this.removeActive = !this.removeActive;
                this.toggleDetails(this.removeActive);
            },
            remove() {
                this.resetError();
                ClientService.deleteClient(this.entity.id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_deleted"));
                        this.$emit("remove", this.entity.id);
                        this.cancel();
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            cancel() {
                this.detailsActive = false;
                this.removeActive = false;
                this.toggleDetails(false);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "key_delete": "Del",
            "shortcut_details": "Details",
            "shortcut_remove": "Remove",
            "confirmation": "Are you sure you would like to delete this client?",
            "label_id": "Client ID",
            "label_secret": "Client Secret",
            "label_scope_organization": "Grant read access to organization details",
            "label_scope_language": "Grant read access to languages",
            "label_scope_articles": "Grant read access to articles",
            "label_scope_article_authors": "Show authors of articles",
            "label_scope_article_authors_mails": "Show author mail addresses",
            "label_scope_article_history": "Grant access to article history",
            "label_scope_lists": "Grant read access to lists",
            "label_scope_tags": "Grant read access to tags",
            "label_scope_pinned": "Grant read access to pinned articles and lists",
            "label_scope_search_all": "Grant access to search and filter all elements",
            "label_scope_search_articles": "Grant access to search and filter articles",
            "label_scope_search_lists": "Grant access to search and filter lists",
            "label_scope_search_tags": "Grant access to search and filter tags",
            "label_no": "No",
            "label_yes": "Yes, delete client",
            "toast_deleted": "The client has been deleted."
        },
        "de": {
            "key_delete": "Entf",
            "shortcut_roles": "Details",
            "shortcut_remove": "Entfernen",
            "confirmation": "Möchtest du diesen Client wirklich löschen?",
            "label_id": "Client ID",
            "label_secret": "Client Secret",
            "label_scope_organization": "Erlaube Lesezugriff auf Organisationsdetails",
            "label_scope_language": "Erlaube Lesezugriff auf Sprachen",
            "label_scope_articles": "Erlaube Lesezugriff auf Artikel",
            "label_scope_article_authors": "Zeige Autoren in Artikeln",
            "label_scope_article_authors_mails": "Zeige E-Mail-Adressen von Autoren",
            "label_scope_article_history": "Erlaube Zugriff auf den Artikelverlauf",
            "label_scope_lists": "Erlaube Lesezugriff auf Listen",
            "label_scope_tags": "Erlaube Lesezugriff auf Tags",
            "label_scope_pinned": "Erlaube Lesezugriff auf angepinnte Artikel und Listen",
            "label_scope_search_all": "Erlaube die Suche und Filterung aller Elemente",
            "label_scope_search_articles": "Erlaube die Artikelsuche und -filterung",
            "label_scope_search_lists": "Erlaube die Listensuche und -filterung",
            "label_scope_search_tags": "Erlaube die Tagsuche und -filterung",
            "label_no": "Nein",
            "label_yes": "Ja, Client löschen",
            "toast_saved": "Gespeichert.",
            "toast_deleted": "Der Client wurde gelöscht."
        }
    }
</i18n>
