<template>
    <div>
        <emvi-cmd-input :label="$t('label_name')"
            :index="0"
            :disabled="!isExpert"
            :error="validation['name']"
            v-model="name"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-checkbox :label="$t('label_scope_organization')"
            :index="1"
            :disabled="!isExpert"
            name="scopeOrganization"
            v-model="scopeOrganization"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_language')"
            :index="2"
            :disabled="!isExpert"
            name="scopeLanguage"
            v-model="scopeLanguage"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_articles')"
            :index="3"
            :disabled="!isExpert"
            name="scopeArticles"
            v-model="scopeArticles"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_article_authors')"
            :index="4"
            :disabled="!isExpert"
            name="scopeArticleAuthors"
            v-model="scopeArticleAuthors"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_article_authors_mails')"
            :index="5"
            :disabled="!isExpert"
            name="scopeArticleAuthorsMails"
            v-model="scopeArticleAuthorsMails"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_article_history')"
            :index="6"
            :disabled="!isExpert"
            name="scopeArticleHistory"
            v-model="scopeArticleHistory"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_lists')"
            :index="7"
            :disabled="!isExpert"
            name="scopeLists"
            v-model="scopeLists"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_tags')"
            :index="8"
            :disabled="!isExpert"
            name="scopeTags"
            v-model="scopeTags"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_pinned')"
            :index="9"
            :disabled="!isExpert"
            name="scopePinned"
            v-model="scopePinned"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_search_all')"
            :index="10"
            :disabled="!isExpert"
            name="scopeSearchAll"
            v-model="scopeSearchAll"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_search_articles')"
            :index="11"
            :disabled="!isExpert"
            name="scopeSearchArticles"
            v-model="scopeSearchArticles"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_search_lists')"
            :index="12"
            :disabled="!isExpert"
            name="scopeSearchLists"
            v-model="scopeSearchLists"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_scope_search_tags')"
            :index="13"
            :disabled="!isExpert"
            name="scopeSearchTags"
            v-model="scopeSearchTags"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-button icon="save"
            :label="isExpert ? $t('label_save') : $t('label_save')+' '+$t('expert')"
            :index="14"
            :disabled="!isExpert"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import {ClientService} from "../../../../service";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdInput, emviCmdCheckbox, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                name: "",
                scopeOrganization: false,
                scopeLanguage: false,
                scopeArticles: false,
                scopeArticleAuthors: false,
                scopeArticleAuthorsMails: false,
                scopeArticleHistory: false,
                scopeLists: false,
                scopeTags: false,
                scopePinned: false,
                scopeSearchArticles: false,
                scopeSearchLists: false,
                scopeSearchTags: false,
                scopeSearchAll: false
            };
        },
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 15, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            }
        },
        methods: {
            save() {
                this.resetError();
                let scopes = [
                    {name: "organization", read: this.scopeOrganization, write: false},
                    {name: "language", read: this.scopeLanguage, write: false},
                    {name: "articles", read: this.scopeArticles, write: false},
                    {name: "article_authors", read: this.scopeArticleAuthors, write: false},
                    {name: "article_authors_mails", read: this.scopeArticleAuthorsMails, write: false},
                    {name: "article_history", read: this.scopeArticleHistory, write: false},
                    {name: "lists", read: this.scopeLists, write: false},
                    {name: "tags", read: this.scopeTags, write: false},
                    {name: "pinned", read: this.scopePinned, write: false},
                    {name: "search_articles", read: this.scopeSearchArticles, write: false},
                    {name: "search_lists", read: this.scopeSearchLists, write: false},
                    {name: "search_tags", read: this.scopeSearchTags, write: false},
                    {name: "search_all", read: this.scopeSearchAll, write: false}
                ];

                ClientService.saveClient(null, this.name, scopes)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_saved"));
                        this.$store.dispatch("popColumn");
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            cancel() {
                this.$store.dispatch("popColumn");
            },
            nextRow() {
                this.$store.dispatch("selectNextRow");
            },
            previousRow() {
                this.$store.dispatch("selectPreviousRow");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "label_name": "Name",
            "label_scope_organization": "Grant read access to Organization details",
            "label_scope_language": "Grant read access to languages",
            "label_scope_articles": "Grant read access to Articles",
            "label_scope_article_authors": "Show authors of Articles",
            "label_scope_article_authors_mails": "Show author email addresses",
            "label_scope_article_history": "Grant access to Article history",
            "label_scope_lists": "Grant read access to Lists",
            "label_scope_tags": "Grant read access to Tags",
            "label_scope_pinned": "Grant read access to pinned Articles and Lists",
            "label_scope_search_all": "Grant access to search and filter all elements",
            "label_scope_search_articles": "Grant access to search and filter Articles",
            "label_scope_search_lists": "Grant access to search and filter Lists",
            "label_scope_search_tags": "Grant access to search and filter Tags",
            "label_save": "Save",
            "toast_saved": "Saved.",
            "expert": "(requires an Expert organization)"
        },
        "de": {
            "label_name": "Name",
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
            "label_save": "Speichern",
            "toast_saved": "Gespeichert.",
            "expert": "(benötigt eine Expert Organisation)"
        }
    }
</i18n>
