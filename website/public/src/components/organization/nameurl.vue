<template>
    <div class="reg-form">
        <h3 class="reg-headline">
            {{$t("title")}}
        </h3>
        <form v-on:submit.prevent="action">
            <field v-model="name" type="text" :label="$t('label_name')" :error="validation['name']" color="grey" autofocus="true"></field>
            <div class="form--element--half">
                <field v-model="domain" type="text" :label="$t('label_domain')" :error="validation['domain']" color="grey"></field>
                <label class="form--element--follow-text">.emvi.com</label>
            </div>
            <div class="row">
                <input type="submit" :value="$t('button_submit')" class="col-sm-12 reg-button--primary" />
            </div>
        </form>
    </div>
</template>

<script>
    import {OrganizationService} from "../../service";
    import field from "../html/field.vue";

    export default {
        components: {field},
        data() {
            return {
                name: "",
                domain: ""
            };
        },
        mounted() {
            let name = window.localStorage.getItem("new_orga_name");
            let domain = window.localStorage.getItem("new_orga_domain");

            if(name) {
                this.name = name;
            }

            if(domain) {
                this.domain = domain;
            }
        },
        methods: {
            action() {
                let data = {name: this.name, domain: this.domain};

                OrganizationService.validateOrganization(data, true, false)
                .then(() => {
                    window.localStorage.setItem("new_orga_name", this.name);
                    window.localStorage.setItem("new_orga_domain", this.domain);
                    this.$emit("next");
                })
                .catch(e => {
                    this.setError(e);
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "title": "Name and URL",
            "button_submit": "Next step",
            "label_name": "Organization name",
            "label_domain": "URL"
        },
        "de": {
            "title": "Name und URL",
            "button_submit": "Nächster Schritt",
            "label_name": "Name der Organisation",
            "label_domain": "URL"
        }
    }
</i18n>
