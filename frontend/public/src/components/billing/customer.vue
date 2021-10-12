<template>
    <div style="margin: 24px 0;">
        <h3>{{$t("headline")}}</h3>
        <div v-if="!showUpdateCustomer">
            <small style="display: block; margin-top: 12px;">{{$t("name")}}</small>
            <p>{{name}}</p>
            <small style="display: block; margin-top: 12px;">{{$t("email")}}</small>
            <p>{{email}}</p>
            <small style="display: block; margin-top: 12px;">{{$t("address")}}</small>
            <p>
                {{addressLine1}}<br />
                <template v-if="addressLine2">{{addressLine2}}<br /></template>
                {{city}}, {{postalCode}}<br />
                {{countryName}}
            </p>
            <template v-if="phone">
                <small style="display: block; margin-top: 12px;">{{$t("phone")}}</small>
                <p>{{phone}}</p>
            </template>
            <template v-if="taxNumber">
                <small style="display: block; margin-top: 12px;">{{$t("tax_number")}}</small>
                <p>{{taxNumber}}</p>
            </template>
            <small style="display: block; margin-top: 12px;">{{$t("balance")}}</small>
            <p>{{balance}}</p>
            <button v-on:click="showUpdateCustomer = true">{{$t("update")}}</button>
        </div>
        <div v-if="showUpdateCustomer">
            <form v-on:submit.stop.prevent="update">
                <emvi-input name="name" :label="$t('name')" v-model="name" :error="validation['name']"></emvi-input>
                <emvi-input name="email" :label="$t('email')" v-model="email" :error="validation['email']"></emvi-input>
                <emvi-select name="country" :options="countryOptions" :label="$t('country')" v-model="country"></emvi-select>
                <emvi-input name="addressLine1" :label="$t('address_line_1')" v-model="addressLine1" :error="validation['address_line_1']"></emvi-input>
                <emvi-input name="addressLine2" :label="$t('address_line_2')" v-model="addressLine2"></emvi-input>
                <emvi-input name="postalCode" :label="$t('postal_code')" v-model="postalCode" :error="validation['postal_code']"></emvi-input>
                <emvi-input name="city" :label="$t('city')" v-model="city" :error="validation['city']"></emvi-input>
                <emvi-input name="phone" :label="$t('phone_form')" v-model="phone"></emvi-input>
                <emvi-input name="taxNumber" :label="$t('tax_number')" v-model="taxNumber"></emvi-input>
                <button type="submit" class="blue-100 bg-blue-10">{{$t("save")}}</button>
                <button v-on:click.stop.prevent="showUpdateCustomer = false">{{$t("cancel")}}</button>
            </form>
        </div>
    </div>
</template>

<script>
    import {getName} from "country-list";
    import {BillingService} from "../../service";
    import emviInput from "../form/input.vue";
    import emviSelect from "../form/select.vue";
    import {getCountries} from "./util";

    export default {
        components: {
            emviInput,
            emviSelect
        },
        data() {
            return {
                showUpdateCustomer: false,
                loading: false,
                customer: {},
                name: "",
                email: "",
                country: "",
                addressLine1: "",
                addressLine2: "",
                postalCode: "",
                city: "",
                phone: "",
                taxNumber: ""
            };
        },
        computed: {
            countryOptions: getCountries,
            balance() {
                let balance = this.customer.balance || 0;
                return `$${(balance/100).toFixed(2)}`;
            },
            countryName() {
                return getName(this.country);
            }
        },
        mounted() {
            this.loadCustomer();
        },
        methods: {
            loadCustomer() {
                BillingService.getCustomer()
                    .then(({customer}) => {
                        this.customer = customer;
                        this.name = customer.name;
                        this.email = customer.email;
                        this.country = customer.country.toLowerCase();
                        this.addressLine1 = customer.address_line_1;
                        this.addressLine2 = customer.address_line_2;
                        this.postalCode = customer.postal_code;
                        this.city = customer.city;
                        this.phone = customer.phone;
                        this.taxNumber = customer.tax_number;
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
                BillingService.updateCustomer(this.email,
                    this.name,
                    this.country,
                    this.addressLine1,
                    this.addressLine2,
                    this.postalCode,
                    this.city,
                    this.phone,
                    this.taxNumber)
                    .then(() => {
                        this.showUpdateCustomer = false;
                        this.$store.dispatch("success", this.$t("saved"));
                        this.$store.dispatch("reload");
                        this.loadCustomer();
                    })
                    .catch(e => {
                        this.loading = false;
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
            "headline": "Customer",
            "address": "Address",
            "balance": "Balance",
            "name": "Name",
            "email": "Email",
            "country": "Country",
            "address_line_1": "Address Line 1",
            "address_line_2": "Address Line 2",
            "postal_code": "Postal Code",
            "city": "City",
            "phone_form": "Phone (optional)",
            "phone": "Phone",
            "tax_number": "Tax Number (for business customers only)",
            "update": "Update",
            "save": "Save",
            "cancel": "Cancel",
            "saved": "Saved."
        },
        "de": {
            "headline": "Kunde",
            "address": "Adresse",
            "balance": "Guthaben",
            "name": "Name",
            "email": "E-Mail-Adresse",
            "country": "Land",
            "address_line_1": "Adresszeile 1",
            "address_line_2": "Adresszeile 2",
            "postal_code": "Postleitzahl",
            "city": "Stadt",
            "phone_form": "Telefon (optional)",
            "phone": "Telefon",
            "tax_number": "Steuernummer (nur für Geschäftskunden)",
            "update": "Ändern",
            "save": "Speichern",
            "cancel": "Abbrechen",
            "saved": "Gespeichert."
        }
    }
</i18n>
