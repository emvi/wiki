<template>
    <div>
        <emvi-cmd-select :label="$t('label_language')"
                         :index="0"
                         :disabled="!isExpert"
                         :options="languageOptions"
                         :error="validation['code']"
                         v-model="language"
                         v-on:next="nextRow"
                         v-on:previous="previousRow"
                         v-on:enter="save"
                         v-on:esc="cancel"></emvi-cmd-select>
        <emvi-cmd-button icon="save"
                         :label="isExpert ? $t('label_save') : $t('label_save')+' '+$t('expert')"
                         :index="1"
                         :disabled="!isExpert"
                         v-on:next="nextRow"
                         v-on:previous="previousRow"
                         v-on:enter="save"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import ISO6391 from "iso-639-1";
    import {updateSelectedRow} from "../../util";
    import {LangService} from "../../../../service";
    import emviCmdSelect from "../../form/select.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdSelect, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                language: ""
            };
        },
        computed: {
            ...mapGetters(["row"]),
            languageOptions() {
                let langs = ISO6391.getLanguages(ISO6391.getAllCodes());
                let options = [];

                for(let i = 0; i < langs.length; i++) {
                    options.push({value: langs[i].code, label: langs[i].nativeName});
                }

                return options;
            }
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 2, this.$store);
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

                LangService.addLang(this.language)
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
            "label_language": "Language",
            "label_save": "Save",
            "toast_saved": "Saved.",
            "expert": "(requires and Expert organization)"
        },
        "de": {
            "label_language": "Sprache",
            "label_save": "Speichern",
            "toast_saved": "Gespeichert.",
            "expert": "(benötigt eine Expert Organisation)"
        }
    }
</i18n>
