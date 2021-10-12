<template>
    <div class="reg-form">
        <h3 class="reg-headline">
            {{$t("title")}}
        </h3>
        <form v-on:submit.prevent="action">
            <field type="text" :label="$t('label_firstname')" v-model="firstname" :error="validation['firstname']" :autofocus="true"></field>
            <field type="text" :label="$t('label_lastname')" v-model="lastname" :error="validation['lastname']"></field>
            <div class="row">
                <input type="submit" :value="$t('button_submit')" class="col-sm-12 reg-button--primary" />
            </div>
        </form>
    </div>
</template>

<script>
    import {AuthService} from "../../service";
    import field from "../html/field.vue";

    export default {
        components: {field},
        props: ["code"],
        data() {
            return {
                firstname: "",
                lastname: ""
            };
        },
        methods: {
            action() {
                AuthService.setRegistrationPersonalData(this.code, this.firstname, this.lastname)
                .then(() => {
                    this.$emit("success");
                })
                .catch(e => {
                    this.setError(e);
                })
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "title": "Your Name",
            "label_firstname": "Firstname",
            "label_lastname": "Lastname",
            "button_submit": "Next Step"
        },
        "de": {
            "title": "Dein Name",
            "label_firstname": "Vorname",
            "label_lastname": "Nachname",
            "button_submit": "Nächster Schritt"
        }
    }
</i18n>
