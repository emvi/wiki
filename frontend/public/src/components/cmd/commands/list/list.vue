<template>
    <div>
        <emvi-cmd-input :label="$t('label_name')+' - '+lang.name"
                        :index="0"
                        :error="validation['name.'+lang.id]"
                        required="true"
                        v-model="name"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="save"
                        v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-input :label="$t('label_info')+' - '+lang.name"
                        :index="1"
                        :error="validation['info.'+lang.id]"
                        v-model="info"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="save"
                        v-on:esc="cancel"></emvi-cmd-input>
        <div v-for="(lang, index) in langs" :key="lang.id">
            <emvi-cmd-input :label="$t('label_name')+' - '+lang.name"
                            :index="index*2+2"
                            :error="validation['name.'+lang.id]"
                            optional="true"
                            v-model="names[lang.id]"
                            v-on:next="nextRow"
                            v-on:previous="previousRow"
                            v-on:enter="save"
                            v-on:esc="cancel"></emvi-cmd-input>
            <emvi-cmd-input :label="$t('label_info')+' - '+lang.name"
                            :index="index*2+3"
                            :error="validation['info.'+lang.id]"
                            optional="true"
                            v-model="infos[lang.id]"
                            v-on:next="nextRow"
                            v-on:previous="previousRow"
                            v-on:enter="save"
                            v-on:esc="cancel"></emvi-cmd-input>
        </div>
        <emvi-cmd-checkbox :label="$t('label_public')"
                           :hint="$t('hint_public')"
                           :index="langs.length*2+2"
                           name="public"
                           v-model="public"
                           v-on:next="nextRow"
                           v-on:previous="previousRow"
                           v-on:enter="save"
                           v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-checkbox :label="$t('label_client_access')"
                           :hint="$t('hint_client_access')"
                           :index="langs.length*2+3"
                           name="client_access"
                           v-model="client_access"
                           v-on:next="nextRow"
                           v-on:previous="previousRow"
                           v-on:enter="save"
                           v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-button icon="list"
                         color="green"
                         :label="isNew ? $t('label_save_new') : $t('label_save_existing')"
                         :index="langs.length*2+4"
                         v-on:next="nextRow"
                         v-on:previous="previousRow"
                         v-on:enter="save"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import {slugWithId} from "../../../../util";
    import {LangService, ArticlelistService} from "../../../../service";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdInput, emviCmdCheckbox, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                langs: [],
                id: "",
                lang: {},
                name: "",
                info: "",
                names: {},
                infos: {},
                public: true,
                client_access: false
            };
        },
        computed: {
            ...mapGetters(["row"]),
            isNew() {
                return !this.id;
            }
        },
        watch: {
            row(row) {
                updateSelectedRow(row, this.langs.length*2+5, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            }
        },
        mounted() {
            this.loadLangs();
        },
        methods: {
            loadLangs() {
                LangService.getLangs()
                    .then(langs => {
                        this.langs = langs;

                        if(this.$route.name === "list") {
                            this.loadList();
                        }
                        else {
                            this.buildNamesAndInfosForNewList();
                        }
                    });
            },
            loadList() {
                this.id = this.$store.state.page.meta.get("id");

                ArticlelistService.getArticlelist(this.id)
                    .then(list => {
                        this.public = list.list.public;
                        this.client_access = list.list.client_access;
                        this.buildNamesAndInfosForExistingList(list.list.names);
                    });
            },
            buildNamesAndInfosForNewList() {
                let removeDefaultIndex = 0;

                for(let i = 0; i < this.langs.length; i++) {
                    if(!this.langs[i].default) {
                        this.names[this.langs[i].id] = "";
                        this.infos[this.langs[i].id] = "";
                    }
                    else {
                        this.lang = this.langs[i];
                        removeDefaultIndex = i;
                    }
                }

                this.langs.splice(removeDefaultIndex, 1);
            },
            buildNamesAndInfosForExistingList(namesAndInfos) {
                for(let i = 0; i < this.langs.length; i++) {
                    if(this.langs[i].default) {
                        this.lang = this.langs[i];
                        this.langs.splice(i, 1);
                        break;
                    }
                }

                let names = [];
                let infos = [];

                for(let i = 0; i < namesAndInfos.length; i++) {
                    if(namesAndInfos[i].language_id !== this.lang.id) {
                        names[namesAndInfos[i].language_id] = namesAndInfos[i].name;
                        infos[namesAndInfos[i].language_id] = namesAndInfos[i].info;
                    }
                    else {
                        this.name = namesAndInfos[i].name;
                        this.info = namesAndInfos[i].info;
                    }
                }

                this.names = names;
                this.infos = infos;
            },
            nextRow() {
                this.$store.dispatch("selectNextRow");
            },
            previousRow() {
                this.$store.dispatch("selectPreviousRow");
            },
            save() {
                this.resetError();
                let names = [];
                names.push({
                    language_id: this.lang.id,
                    name: this.name,
                    info: this.info
                });

                for(let lang in this.names) {
                    names.push({
                        language_id: lang,
                        name: this.names[lang],
                        info: this.infos[lang]
                    });
                }

                ArticlelistService.saveArticlelist(this.id, names, this.public, this.client_access)
                    .then(id => {
                        this.$store.dispatch("success", this.$t(this.id ? "toast_saved_existing" : "toast_saved_new"));
                        this.$store.dispatch("resetCmd");

                        if(!this.id) {
                            this.$router.push(`/list/${slugWithId(this.name, id)}`);
                        }
                        else {
                            this.$store.dispatch("setMetaVars", [
                                {key: "updated", value: true},
                                {key: "name", value: this.name},
                                {key: "info", value: this.info},
                                {key: "public", value: this.public},
                                {key: "client_access", value: this.client_access}
                            ]);
                        }
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            cancel() {
                this.$store.dispatch("popColumn");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "label_name": "Name",
            "label_info": "Description",
            "label_public": "Public List",
            "label_client_access": "API access",
            "label_save_new": "Create List",
            "label_save_existing": "Update List",
            "hint_public": "Any member of this Organization can find this List.",
            "hint_client_access": "Enables access through the client API.",
            "toast_saved_new": "The list has been created.",
            "toast_saved_existing": "The list has been updated."
        },
        "de": {
            "label_name": "Name",
            "label_info": "Beschreibung",
            "label_public": "Öffentliche Liste",
            "label_client_access": "API Zugriff",
            "label_save_new": "Liste erstellen",
            "label_save_existing": "Liste speichern",
            "hint_public": "Jedes Mitglieder dieser Organisation kann die Liste finden.",
            "hint_client_access": "Aktiviert den Zugriff über die Client API.",
            "toast_saved_new": "Die Liste wurde erstellt.",
            "toast_saved_existing": "Die Liste wurde aktualisiert."
        }
    }
</i18n>
