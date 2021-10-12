<template>
    <div style="margin: 24px 0;">
        <h3>{{$t("headline")}}</h3>
        <div v-if="!showUpdateCard">
            <small style="display: block; margin-top: 12px;">{{$t("card")}}</small>
            <p style="text-transform: capitalize;">{{paymentMethod.brand}} ...{{paymentMethod.last_4}}</p>
            <small style="display: block; margin-top: 12px;">{{$t("expiration")}}</small>
            <p>{{paymentMethod.exp_month}}/{{paymentMethod.exp_year}}</p>
            <button v-on:click="showUpdateCard = true">{{$t("update")}}</button>
        </div>
        <div v-show="showUpdateCard">
            <form v-on:submit.stop.prevent="update" style="margin: 16px 0;">
                <div id="card-element"></div>
                <button type="submit" class="blue-100 bg-blue-10">{{$t("save")}}</button>
                <button v-on:click.stop.prevent="showUpdateCard = false">{{$t("cancel")}}</button>
                <p class="red-100" v-if="error">{{error}}</p>
            </form>
        </div>
    </div>
</template>

<script>
    import {BillingService} from "../../service";
    import {getCodeList} from "country-list";

    export default {
        data() {
            return {
                stripe: null,
                cardElement: null,
                showUpdateCard: false,
                loading: false,
                error: "",
                paymentMethod: {}
            };
        },
        mounted() {
            this.stripe = Stripe(EMVI_WIKI_STRIPE_PUBLIC_KEY);
            this.createCardElement();
            this.loadCard();
        },
        methods: {
            createCardElement() {
                let elements = this.stripe.elements();
                this.cardElement = elements.create("card");
                this.cardElement.mount("#card-element");
                this.cardElement.on("change", e => {
                    if(e.error) {
                        this.error = e.error.message;
                    }
                    else {
                        this.error = "";
                    }
                });
            },
            loadCard() {
                BillingService.getCustomer()
                    .then(({payment_method}) => {
                        this.paymentMethod = payment_method;
                    })
                    .catch(e => {
                        console.error(e);
                        this.setError(e);
                    });
            },
            update() {
                if(this.loading) {
                    return;
                }

                this.loading = true;
                this.stripe.createPaymentMethod({type: "card", card: this.cardElement})
                    .then((result) => {
                        if(result.error) {
                            this.loading = false;
                            this.error = result.error.message;
                        }
                        else {
                            BillingService.updatePaymentMethod(result.paymentMethod.id)
                                .then(() => {
                                    this.showUpdateCard = false;
                                    this.$store.dispatch("reload");
                                    this.$store.dispatch("success", this.$t("saved"));
                                    this.loadCard();
                                })
                                .catch(e => {
                                    this.loading = false;
                                    console.error(e);
                                    this.setError(e);
                                });
                        }
                    });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "headline": "Card",
            "card": "Card",
            "expiration": "Expiration",
            "update": "Update",
            "save": "Save",
            "cancel": "Cancel",
            "saved": "Saved."
        },
        "de": {
            "headline": "Kreditkarte",
            "card": "Kreditkarte",
            "expiration": "Ablaufdatum",
            "update": "Ändern",
            "save": "Speichern",
            "cancel": "Abbrechen",
            "saved": "Gespeichert."
        }
    }
</i18n>
