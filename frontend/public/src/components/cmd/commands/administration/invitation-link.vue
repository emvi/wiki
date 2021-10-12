<template>
    <div>
        <p>{{$t("text")}}</p>
        <emvi-cmd-link :index="0"
            :link="link"
            v-on:click="copy"></emvi-cmd-link>
        <emvi-cmd-checkbox :label="$t('label_read_only')"
            :index="1"
            v-model="readOnly"
            name="read_only"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-button icon="link"
            :label="$t('label_generate')"
            :index="2"
            v-on:next="nextRow"
            v-on:previous="previousRow"
            v-on:enter="run"
            v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import {copyToClipboard} from "../../../../util";
    import {OrganizationService} from "../../../../service";
    import emviCmdLink from "../../content/link.vue";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdLink, emviCmdCheckbox, emviCmdButton},
        props: ["enter", "esc", "up", "down"],
        data() {
            return {
                link: "",
                readOnly: false
            };
        },
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 3, this.$store);
            },
            enter(enter) {
                if(enter) {
                    this.run();
                }
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            },
            up(up) {
                if(up) {
                    this.$store.dispatch("selectPreviousRow");
                }
            },
            down(down) {
                if(down) {
                    this.$store.dispatch("selectNextRow");
                }
            }
        },
        mounted() {
            this.loadCode();
        },
        methods: {
            loadCode() {
                OrganizationService.getInvitationCode()
                    .then(({code, read_only}) => {
                        this.link = code ? this.getLink(code) : "";
                        this.readOnly = read_only;
                    })
                    .catch(e => {
                        this.showTechnicalError(e);
                    });
            },
            run() {
                if(this.row === 0) {
                    this.copy();
                }
                else {
                    this.generate();
                }
            },
            copy() {
                if(this.link) {
                    copyToClipboard(this.link);
                    this.$store.dispatch("success", this.$t("toast_copied"));
                }
            },
            generate() {
                this.focusCmdInput();
                OrganizationService.generateInvitationCode(this.readOnly)
                    .then(code => {
                        this.link = this.getLink(code);
                    })
                    .catch(e => {
                        this.showTechnicalError(e);
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
            },
            getLink(code) {
                return `${EMVI_WIKI_WEBSITE_HOST}/join/${code}`;
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "text": "Anyone with an Emvi account and this link can join your organization. In case you're unsure who has it, please generate a new link.",
            "label_read_only": "Read only access",
            "label_generate": "Generate New Invitation Link",
            "toast_copied": "Link copied to clipboard."
        },
        "de": {
            "text": "Jeder mit einem Emvi Account und diesem Link kann der Organisation beitreten. Wenn du unsicher bist wer den Link hat, generiere bitte einen neuen.",
            "label_read_only": "Nur Lesezugriff",
            "label_generate": "Neuen Einladungslink generieren",
            "toast_copied": "Link in die Zwischenablage kopiert."
        }
    }
</i18n>
