<template>
    <div>
        <emvi-cmd-input :label="$t('label_commit')"
            :index="0"
            :disabled="!publishNow"
            v-model="message"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="publish"
            v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-checkbox :label="$t('label_save')"
            :hint="$t('hint_save')"
            :index="1"
            v-model="publishNow"
            name="publish"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="publish"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-select :label="$t('label_language')"
            :index="2"
            :options="languageOptions"
            :disabled="!!articleId.length"
            v-model="language"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="publish"
            v-on:esc="cancel"></emvi-cmd-select>
        <emvi-cmd-button icon="publish"
            :label="$t('label_publish')"
            :index="3"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="publish"
            v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import {LangService} from "../../../../service";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import emviCmdSelect from "../../form/select.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdInput, emviCmdCheckbox, emviCmdSelect, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                articleId: this.$store.state.page.meta.get("id"),
                languageOptions: [],
                message: "",
                publishNow: true,
                language: this.$store.state.page.meta.get("langId")
            };
        },
        computed: {
            ...mapGetters(["row", "cmdMeta"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 4, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            },
            language(language) {
                this.$store.dispatch("setMeta", {key: "setLanguage", value: language});
            }
        },
        mounted() {
            this.loadLanguages();
        },
        methods: {
            loadLanguages() {
                LangService.getLangs()
                    .then(langs => {
                        let options = [];

                        for(let i = 0; i < langs.length; i++) {
                            options.push({label: langs[i].name, value: langs[i].id});

                            if(langs[i].default && !this.language.length) {
                                this.language = langs[i].id;
                            }
                        }

                        this.languageOptions = options;
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            publish() {
                if(!this.$store.state.page.meta.get("title").trim()) {
                    // see ErrorService for reference
                    this.setError({errors: [{message: "Title empty"}]});
                    return;
                }

                let page = this.$store.state.page.meta.get("page");
                page.save(this.message, this.publishNow);
                this.$store.dispatch("resetCmd");
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
            "label_commit": "Change description (optional)",
            "label_save": "Publish now",
            "label_language": "Language",
            "label_publish": "Publish",
            "hint_save": "If selected, the changes are immediately visible to readers. Otherwise the changes are saved and are not visible until the next publication."
        },
        "de": {
            "label_commit": "Änderungsbeschreibung (optional)",
            "label_save": "Sofort veröffentlichen",
            "label_language": "Sprache",
            "label_publish": "Veröffentlichen",
            "hint_save": "Wenn angewählt, werden die Änderungen sofort für Leser sichtbar. Ansonsten werden die Änderungen gespeichert und sind bis zur nächsten Veröffentlichung nicht sichtbar."
        }
    }
</i18n>
