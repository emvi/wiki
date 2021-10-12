<template>
    <div>
        <emvi-billing-invoice v-for="invoice in invoices"
            :key="invoice.id"
            :invoice="invoice"></emvi-billing-invoice>
        <button v-on:click="loadInvoices(true)">{{$t("load_more")}}</button>
    </div>
</template>

<script>
    import {BillingService} from "../../service";
    import emviBillingInvoice from "./invoice.vue";

    export default {
        components: {emviBillingInvoice},
        data() {
            return {
                invoices: []
            };
        },
        mounted() {
            this.loadInvoices();
        },
        methods: {
            loadInvoices(more = false) {
                let startInvoiceId = "";

                if(more && this.invoices.length) {
                    startInvoiceId = this.invoices[this.invoices.length-1].id;
                }

                BillingService.getInvoices(startInvoiceId)
                    .then(invoices => {
                        this.invoices = this.invoices.concat(invoices);
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
            "headline": "Invoices",
            "load_more": "Load more"
        },
        "de": {
            "headline": "Rechnungen",
            "load_more": "Mehr laden"
        }
    }
</i18n>
