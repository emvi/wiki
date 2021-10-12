<template>
    <div>
        <emvi-cmd-select :label="$t('label_format')"
            :index="0"
            :options="formatOptions"
            v-model="format"
            :error="err"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-select>
        <emvi-cmd-checkbox :label="$t('label_attachments')"
            :index="1"
            name="attachments"
            v-model="exportAttachments"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-button icon="download"
            color="green"
            :label="$t('label_action')"
            :index="2"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import axios from "axios";
    import {mapGetters} from "vuex";
    import {slugWithId} from "../../../../util";
    import {updateSelectedRow} from "../../util";
    import emviCmdSelect from "../../form/select.vue";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdSelect, emviCmdCheckbox, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                formatOptions: [
                    {value: "html", label: "HTML"},
                    {value: "markdown", label: "Markdown"},
                ],
                format: "html",
                exportAttachments: true,
                err: ""
            };
        },
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 3, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            }
        },
        methods: {
            run() {
                let id = this.$store.state.page.meta.get("id");
                let langId = this.$store.state.page.meta.get("langId");
                let title = this.$store.state.page.meta.get("title");
                let url = `${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${id}/export?language_id=${langId}&format=${this.format}&export_attachments=${this.exportAttachments}`;

                axios.get(url, {responseType: "blob"})
                    .then(r => {
                        url = window.URL.createObjectURL(new Blob([r.data]));
                        let link = document.createElement('a');
                        link.href = url;
                        link.setAttribute('download', `${slugWithId(title, id)}.zip`);
                        link.click();
                        this.$store.dispatch("resetCmd");
                    })
                    .catch(e => {
                        this.err = this.$t("export_error");
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
            "label_format": "Format",
            "label_attachments": "Include Attachments",
            "label_action": "Export",
            "export_error": "The article could not be exported. Is it published?"
        },
        "de": {
            "label_format": "Format",
            "label_attachments": "Inklusive Anhänge",
            "label_action": "Exportieren",
            "export_error": "Der Artikel konnte nicht veröffentlicht werden. Ist er veröffentlicht?"
        }
    }
</i18n>
