<template>
    <div style="margin: 24px 0;">
        <h3>{{$t("headline")}}</h3>
        <div style="margin-top: 12px;" v-if="!organization.subscription_cancelled">
            <p v-html="$t('text_plan', {plan: $t(organization.subscription_plan)})"></p>
            <div v-if="!showChangePlan">
                <button v-on:click="showChangePlan = true">{{$t("change_plan")}}</button>
            </div>
            <div v-if="showChangePlan">
                <form v-on:submit.stop.prevent="changePlan">
                    <div class="plan">
                        <emvi-billing-price plan="yearly" v-model="plan"></emvi-billing-price>
                        <emvi-billing-price plan="monthly" v-model="plan"></emvi-billing-price>
                    </div>
                    <button type="submit" class="blue-100 bg-blue-10">{{$t("save")}}</button>
                    <button v-on:click.stop.prevent="showChangePlan = false">{{$t("cancel")}}</button>
                </form>
            </div>
        </div>
        <div style="margin-top: 16px;">
            <p>{{$t("text_subscription")}}</p>
            <button class="red-100 bg-red-10" v-on:click="cancelSubscription" v-if="!organization.subscription_cancelled">{{$t("cancel_subscription")}}</button>
            <button class="blue-100 bg-blue-10" v-on:click="resumeSubscription" v-else>{{$t("resume_subscription")}}</button>
        </div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {BillingService} from "../../service";
    import emviBillingPrice from "./price.vue";

    export default {
        components: {emviBillingPrice},
        data() {
            return {
                showChangePlan: false,
                plan: "yearly"
            };
        },
        computed: {
            ...mapGetters(["organization"])
        },
        mounted() {
            this.plan = this.organization.subscription_plan;
        },
        methods: {
            cancelSubscription() {
                BillingService.cancelSubscription()
                    .then(() => {
                        this.$store.dispatch("reload");
                        this.$store.dispatch("success", this.$t("saved"));
                    })
                    .catch(e => {
                        console.error(e);
                        this.setError(e);
                    });
            },
            resumeSubscription() {
                BillingService.resumeSubscription()
                    .then(() => {
                        this.$store.dispatch("reload");
                        this.$store.dispatch("success", this.$t("saved"));
                    })
                    .catch(e => {
                        console.error(e);
                        this.setError(e);
                    });
            },
            changePlan() {
                BillingService.changePlan(this.plan)
                    .then(() => {
                        this.showChangePlan = false;
                        this.$store.dispatch("reload");
                        this.$store.dispatch("success", this.$t("saved"));
                    })
                    .catch(e => {
                        console.error(e);
                        this.setError(e);
                    });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "headline": "Subscription and Plan",
            "text_subscription": "You can cancel your subscription at any time. The subscription will end with your next billing cycle. You can resume your subscription before it runs out or start a new one afterwards.",
            "text_plan": "You can change your plan at any time. If you change from the monthly to the annual plan, a new invoice will be created and vice versa we will add the difference to your credit. At the moment you use the <strong>{plan}</strong> payment.",
            "monthly": "monthly",
            "yearly": "yearly",
            "cancel_subscription": "Cancel Subscription",
            "resume_subscription": "Resume Subscription",
            "change_plan": "Change Plan",
            "save": "Save",
            "cancel": "Cancel",
            "saved": "Saved."
    },
        "de": {
            "headline": "Abonnement und Plan",
            "text_subscription": "Du kannst dein Abonnement jederzeit beenden. Das Abonnement läuft dann mit dem Ende des Abrechnungszeitraums aus. Du kannst dein Abonnement fortsetzen, bevor es beendet wird oder später danach ein Neues beginnen.",
            "text_plan": "Du kannst deinen Plan jederzeit ändern. Wenn du vom Monats- zum Jahresplan wechseln, wird eine neue Rechnung erstellt. Umgekehrt addieren wir die Differenz zu deinem Guthaben. Zur Zeit verwendest du die <strong>{plan}</strong> Zahlung.",
            "monthly": "monatliche",
            "yearly": "jährliche",
            "cancel_subscription": "Abonnement beenden",
            "resume_subscription": "Abonnement fortsetzen",
            "change_plan": "Plan ändern",
            "save": "Speichern",
            "cancel": "Abbrechen",
            "saved": "Gespeichert."
        }
    }
</i18n>
