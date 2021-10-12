<template>
    <div>
        <p v-html="$t('text', {pricing})"></p>
        <form v-on:submit.stop.prevent="subscribe">
            <div class="plan">
                <emvi-billing-price plan="yearly" v-model="plan"></emvi-billing-price>
                <emvi-billing-price plan="monthly" v-model="plan"></emvi-billing-price>
            </div>
            <p>{{$t("first_invoice")}}</p>
            <h1 style="margin: 16px 0;">
                <span>${{price}}</span>
                <span v-show="plan === 'yearly'">{{$t("yearly")}}</span>
                <span v-show="plan === 'monthly'">{{$t("monthly")}}</span>
            </h1>
            <p>{{$t("tax_rate", {tax})}}</p>
            <p>{{$t("active_members", {billableMemberCount})}}</p>
            <emvi-input name="name" :label="$t('name')" v-model="name" :error="validation['name']"></emvi-input>
            <emvi-input name="email" :label="$t('email')" v-model="email" :error="validation['email']"></emvi-input>
            <emvi-select name="country" :options="countryOptions" :label="$t('country')" v-model="country"></emvi-select>
            <emvi-input name="addressLine1" :label="$t('address_line_1')" v-model="addressLine1" :error="validation['address_line_1']"></emvi-input>
            <emvi-input name="addressLine2" :label="$t('address_line_2')" v-model="addressLine2"></emvi-input>
            <emvi-input name="postalCode" :label="$t('postal_code')" v-model="postalCode" :error="validation['postal_code']"></emvi-input>
            <emvi-input name="city" :label="$t('city')" v-model="city" :error="validation['city']"></emvi-input>
            <emvi-input name="phone" :label="$t('phone')" v-model="phone"></emvi-input>
            <emvi-input name="taxNumber" :label="$t('tax_number')" v-model="taxNumber"></emvi-input>
            <div id="card-element"></div>
            <button class="blue-100 bg-blue-10" type="submit">
                <i class="icon icon-send" v-show="!loading"></i>
                <i class="icon icon-sync icon-is-spinning" v-show="loading"></i>
                {{$t("subscribe")}}
            </button>
            <p class="red-100" v-if="error">{{error}}</p>
        </form>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {BillingService, OrganizationService} from "../../service";
    import emviInput from "../form/input.vue";
    import emviSelect from "../form/select.vue";
    import emviBillingPrice from "./price.vue";
    import {calculateSum, getCountries, getTaxRate} from "./util";

    export default {
        components: {
            emviInput,
            emviSelect,
            emviBillingPrice
        },
        data() {
            return {
                stripe: null,
                cardElement: null,
                loading: false,
                error: "",
                plan: "yearly",
                name: "",
                email: "",
                country: "de", // pre-select Germany by default :)
                addressLine1: "",
                addressLine2: "",
                postalCode: "",
                city: "",
                phone: "",
                taxNumber: "",
                billableMemberCount: 0
            };
        },
        computed: {
            ...mapGetters(["organization", "darkmode"]),
            countryOptions: getCountries,
            pricing() {
                return `${EMVI_WIKI_WEBSITE_HOST}/pricing`;
            },
            price() {
                let sum = calculateSum(this.plan, this.billableMemberCount, this.country, this.taxNumber);

                if(this.$i18n.locale === "de") {
                    sum = sum.replace(".", ",");
                }

                return sum;
            },
            tax() {
                return getTaxRate(this.country, this.taxNumber)*100;
            }
        },
        mounted() {
            this.loadStatistics();
            this.stripe = Stripe(EMVI_WIKI_STRIPE_PUBLIC_KEY);
            this.createCardElement();
        },
        methods: {
            loadStatistics() {
                OrganizationService.getStatistics()
                    .then(statistics => {
                        this.billableMemberCount = statistics.billable_member_count;
                    });
            },
            createCardElement() {
                let elements = this.stripe.elements();
                let textColor = !this.darkmode ? "#000" : "#e1e3e5";
                let errorColor = !this.darkmode ? "#ff1a1a" : "#ff7575";
                this.cardElement = elements.create("card", {
                    style: {
                        base: {
                            fontSize: "16px",
                            iconColor: textColor,
                            color: textColor,
                            "::placeholder": {
                                color: "#797c80",
                            }
                        },
                        invalid: {
                            iconColor: errorColor,
                            color: errorColor,
                        }
                    }
                });
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
            subscribe() {
                if(this.loading) {
                    return;
                }

                this.loading = true;

                // 1. create payment method
                // 2. create customer and subscription (backend)
                // 3. confirm card payment if required (3D secure)
                this.stripe.createPaymentMethod({type: "card", card: this.cardElement})
                    .then((result) => {
                        if(result.error) {
                            this.loading = false;
                            this.error = result.error.message;
                        }
                        else {
                            let paymentMethodId = result.paymentMethod.id;

                            BillingService.subscribe(this.email,
                                this.name,
                                this.country,
                                this.addressLine1,
                                this.addressLine2,
                                this.postalCode,
                                this.city,
                                this.phone,
                                this.taxNumber,
                                this.plan,
                                paymentMethodId)
                                .then(client_secret => {
                                    if(!client_secret.length) {
                                        this.subscriptionCompleted();
                                    }
                                    else {
                                        this.stripe.confirmCardPayment(client_secret, {payment_method: paymentMethodId})
                                            .then((result) => {
                                                if(result.error) {
                                                    this.loading = false;
                                                    this.error = result.error.message;
                                                }
                                                else if(result.paymentIntent.status === "succeeded") {
                                                    this.subscriptionCompleted();
                                                }
                                            });
                                    }
                                })
                                .catch(e => {
                                    this.loading = false;
                                    console.error(e);
                                    this.setError(e);
                                });
                        }
                    });
            },
            subscriptionCompleted() {
                this.$store.dispatch("reload");
                this.$store.dispatch("success", this.$t("success"));
            }
        }
    }
</script>

<style>
    fieldset {
        margin: 16px 0;
    }
</style>

<i18n>
    {
        "en": {
            "text": "An Expert Organization offers many additional features and advantages. All Details can be found on our <a href=\"{pricing}\" target=\"_blank\">price page</a>.",
            "first_invoice": "When subscribing you will be charged the following amount:",
            "yearly": "yearly",
            "monthly": "monthly",
            "tax_rate": "Including {tax}% VAT.",
            "active_members": "For {billableMemberCount} active members with write access. Read-only members are free.",
            "name": "Name",
            "email": "Email",
            "country": "Country",
            "address_line_1": "Address Line 1",
            "address_line_2": "Address Line 2",
            "postal_code": "Postal Code",
            "city": "City",
            "phone": "Phone (optional)",
            "tax_number": "Tax ID (for EU business customers only)",
            "subscribe": "Subscribe",
            "success": "Thank you for your subscription. Your organization has been upgraded to Expert!"
        },
        "de": {
            "text": "Eine Expert-Organisation bietet viele Zusatzfunktionen und Vorteile. Alle Details dazu findest du auf unserer <a href=\"{pricing}\" target=\"_blank\">Preisseite</a>.",
            "first_invoice": "Bei Abschluss wird folgender Betrag berechnet:",
            "yearly": "jährlich",
            "monthly": "monatlich",
            "tax_rate": "Einschließlich {tax}% Mehrwertsteuer.",
            "active_members": "Für {billableMemberCount} aktive Mitglieder mit Schreibzugriff. Mitglieder mit ausschließlich Lesezugriff sind kostenlos.",
            "name": "Name",
            "email": "E-Mail-Adresse",
            "country": "Land",
            "address_line_1": "Adresszeile 1",
            "address_line_2": "Adresszeile 2",
            "postal_code": "Postleitzahl",
            "city": "Stadt",
            "phone": "Telefon (optional)",
            "tax_number": "Steueridentifikationsnummer (nur für Geschäftskunden innerhalb der EU)",
            "subscribe": "Abonnieren",
            "success": "Vielen Dank für dein Abonnement! Deine Organisation wurde auf auf Expert geändert!"
        }
    }
</i18n>
