<template>
    <div class="reg-form">
        <h3 class="reg-headline">
            {{$t("title")}}
        </h3>
        <form v-on:submit.prevent="action">
            <field v-model="username" type="text" :label="$t('label_username')" :hint="$t('hint_username')" :error="validation['username']" color="grey" autofocus="true"></field>
            <selection v-model="language" :label="$t('label_language')" :hint="$t('hint_language')">
                <option v-for="l in langs" :value="l.code" :key="l.code">{{l.nativeName}}</option>
            </selection>
            <div class="row">
                <input type="submit" :value="$t('button_submit')" class="col-sm-12 reg-button--primary" />
            </div>
        </form>
    </div>
</template>

<script>
    import ISO6391 from "iso-639-1";
    import {OrganizationService} from "../../service";
    import field from "../html/field.vue";
    import selection from "../html/selection.vue";

    export default {
        components: {field, selection},
        data() {
            return {
                username: "",
                language: "en"
            };
        },
        computed: {
            langs() {
                return ISO6391.getLanguages(ISO6391.getAllCodes());
            }
        },
        mounted() {
            let username = window.localStorage.getItem("new_orga_username");
            let lang = window.localStorage.getItem("new_orga_default_lang");

            if(username) {
                this.username = username;
            }

            if(lang) {
                this.language = lang;
            }
        },
        methods: {
            action() {
                let data = {username: this.username, default_lang: this.language};

                OrganizationService.validateOrganization(data, false, true)
                .then(() => {
                    window.localStorage.setItem("new_orga_username", this.username);
                    window.localStorage.setItem("new_orga_default_lang", this.language);
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
            "title": "Username and language",
            "button_submit": "Next step",
            "label_username": "Your username",
            "label_language": "Language",
            "hint_username": "Cannot be changed",
            "hint_language": "Select in which language articles are primarly created in. Organizations on the Expert plan can add more languages later."
        },
        "de": {
            "title": "Nutzername und Sprache",
            "button_submit": "Nächster Schritt",
            "label_username": "Dein Nutzername",
            "label_language": "Sprache",
            "hint_username": "Kann nicht mehr geändert werden",
            "hint_language": "Lege fest in welcher Sprache Artikel primär erstellt werden. Expert-Organisationen können später weitere Sprachen hinzugefügen."
        }
    }
</i18n>
