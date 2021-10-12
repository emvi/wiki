<template>
    <div>
        <p>{{$t("text")}}</p>
        <emvi-cmd-checkbox v-for="(recommendation, index) in recommendations"
                           :key="recommendation.id"
                           :label="getLabel(recommendation)"
                           :index="index"
                           v-model="recommendation.confirm"
                           name="recommendation"
                           v-on:next="nextRow"
                           v-on:previous="previousRow"
                           v-on:enter="confirm"
                           v-on:esc="cancel"></emvi-cmd-checkbox>
        <emvi-cmd-button icon="send"
                         :label="$t('label_confirm')"
                         :index="recommendations.length+1"
                         v-on:next="nextRow"
                         v-on:previous="previousRow"
                         v-on:enter="confirm"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {updateSelectedRow} from "../../util";
    import {ArticleService} from "../../../../service";
    import emviCmdCheckbox from "../../form/checkbox.vue";
    import emviCmdButton from "../../form/button.vue";

    export default {
        components: {emviCmdCheckbox, emviCmdButton},
        props: ["enter", "tab", "esc", "up", "down"],
        data() {
            return {
                recommendations: []
            };
        },
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            row(row) {
                updateSelectedRow(row, this.recommendations.length+1, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            }
        },
        beforeMount() {
            this.recommendations = this.$store.state.page.meta.get("recommendations");
        },
        methods: {
            confirm() {
                this.resetError();
                let articleId = this.$store.state.page.meta.get("id");
                let userIds = [];

                for(let i = 0; i < this.recommendations.length; i++) {
                    if(this.recommendations[i].confirm) {
                        userIds.push(this.recommendations[i].organization_member.user.id);
                    }
                }

                ArticleService.confirmRecommendations(articleId, userIds)
                    .then(() => {
                        this.$store.dispatch("resetCmd");
                        this.$store.dispatch("closeCmd");
                        this.$store.dispatch("success", this.$t("toast_invited"));
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            getLabel(recommendation) {
                return `${recommendation.organization_member.user.firstname} ${recommendation.organization_member.user.lastname}`;
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
            "text": "The following members requested a read confirmation from you. Please select selects the members you would like to send an answer to.",
            "label_confirm": "Confirm",
            "toast_invited": "Confirmations have been send."
        },
        "de": {
            "text": "Die folgenden Mitglieder haben eine Lesebestätigung angefordert. Bitte wähle aus, wem du eine Antwort senden möchtest.",
            "label_confirm": "Bestätigen",
            "toast_invited": "Die Bestätigungen wurden gesendet."
        }
    }
</i18n>
