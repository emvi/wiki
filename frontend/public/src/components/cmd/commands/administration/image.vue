<template>
    <div>
        <emvi-cmd-input :label="$t('label_file')"
                        :index="0"
                        :error="err"
                        type="file"
                        :hint="$t('hint_picture')"
                        v-model="file"
                        v-on:change="setFile"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="save"
                        v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-button :icon="icon"
                         color="green"
                         :label="$t('label_save')"
                         :index="1"
                         v-on:next="nextRow"
                         v-on:previous="previousRow"
                         v-on:enter="save"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import {OrganizationService} from "../../../../service";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdButton from "../../form/button.vue";

    const MAX_SIZE_BYTES = 512000; // 512 KB

    // TODO reimplement deleting image
    export default {
        components: {emviCmdInput, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                icon: "img",
                err: "",
                file: null,
                picture: null
            };
        },
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 2, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.$store.dispatch("popColumn");
                }
            }
        },
        methods: {
            setFile(e) {
                this.picture = e.target.files[0];
            },
            save() {
                if(!this.picture) {
                    this.err = this.$t("err_img");
                    return;
                }

                if(this.picture.size > MAX_SIZE_BYTES) {
                    this.err = this.$t("err_max_size");
                    return;
                }

                this.resetError();
                let form = new FormData();
                form.append("file", this.picture);
                this.icon = "sync";

                OrganizationService.uploadPicture(form)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_saved"));
                        this.$store.dispatch("loadOrganization");
                        this.$store.dispatch("popColumn");
                    })
                    .catch(e => {
                        this.icon = "img";
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
            "label_file": "Image",
            "hint_picture": "The maximum size for the picture is 512 KB.",
            "label_save": "Upload",
            "err_img": "Please select a picture to upload.",
            "err_max_size": "The file exceeds the maximum size.",
            "toast_saved": "The picture has been saved."
        },
        "de": {
            "label_file": "Bild",
            "hint_picture": "Die maximale Größe für das Bild beträgt 512 KB.",
            "label_save": "Hochladen",
            "err_img": "Bitte wähle ein Bild aus.",
            "err_max_size": "Die Datei überschreitet die maximale Größe.",
            "toast_saved": "Das Bild wurde gespeichert."
        }
    }
</i18n>
