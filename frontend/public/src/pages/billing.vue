<template>
    <emvi-layout disable-events="true">
        <emvi-layout-narrow>
            <div class="headline">
                <h1 v-if="isExpert">{{$t("billing")}}</h1>
                <h1 v-if="!isExpert">{{$t("upgrade")}}</h1>
            </div>
            <emvi-billing-payment-confirmation></emvi-billing-payment-confirmation>
            <emvi-billing-subscribe v-if="!isExpert"></emvi-billing-subscribe>
            <div class="tabs no-select" v-if="isExpert">
                <a href="#" :class="{'router-link-active': tab === 0}" v-on:click.stop.prevent="tab = 0">
                    <h2>{{$t("tab_subscription")}}</h2>
                </a>
                <a href="#" :class="{'router-link-active': tab === 1}" v-on:click.stop.prevent="tab = 1">
                    <h2>{{$t("tab_invoices")}}</h2>
                </a>
            </div>
            <div v-if="isExpert && tab === 0">
                <emvi-billing-customer></emvi-billing-customer>
                <emvi-billing-payment></emvi-billing-payment>
                <emvi-billing-subscription></emvi-billing-subscription>
            </div>
            <div v-if="isExpert && tab === 1">
                <emvi-billing-invoices></emvi-billing-invoices>
            </div>
        </emvi-layout-narrow>
    </emvi-layout>
</template>

<script>
    import {TitleMixin} from "./title";
    import {
        emviLayout,
        emviLayoutNarrow,
        emviBillingPaymentConfirmation,
        emviBillingCustomer,
        emviBillingPayment,
        emviBillingInvoices,
        emviBillingSubscribe,
        emviBillingSubscription
    } from "../components";

    export default {
        mixins: [TitleMixin],
        components: {
            emviLayout,
            emviLayoutNarrow,
            emviBillingPaymentConfirmation,
            emviBillingCustomer,
            emviBillingPayment,
            emviBillingInvoices,
            emviBillingSubscribe,
            emviBillingSubscription
        },
        data() {
            return {
                tab: 0
            };
        }
    }
</script>

<i18n>
    {
        "en": {
            "title": "Billing",
            "billing": "Billing",
            "upgrade": "Upgrade",
            "tab_subscription": "Subscription",
            "tab_invoices": "Invoices"
        },
        "de": {
            "title": "Abrechnung",
            "billing": "Abrechnung",
            "upgrade": "Upgrade",
            "tab_subscription": "Abonnement",
            "tab_invoices": "Rechnungen"
        }
    }
</i18n>
