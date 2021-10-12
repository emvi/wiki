<template>
    <div class="reg-form">
        <h3 class="reg-headline">
            {{$t("title")}}
        </h3>
        <p>{{$t("text")}}</p>
        <p class="error" v-if="err">{{err}}</p>
        <small>{{$t("orga_name")}}</small>
        <p>{{name}}</p>
        <br />
        <small>{{$t("url")}}</small>
        <p>{{domain}}.emvi.com</p>
        <br />
        <small>{{$t("username")}}</small>
        <p>{{username}}</p>
        <br />
        <small>{{$t("language")}}</small>
        <p>{{lang}}</p>
        <form v-on:submit.prevent="action">
            <div class="row">
                <input type="submit" :value="$t('button_submit')" class="col-sm-12 reg-button--primary" />
            </div>
        </form>
    </div>
</template>

<script>
    import {OrganizationService} from "../../service";
    import ISO6391 from "iso-639-1";

    export default {
        data() {
            return {
                name: window.localStorage.getItem("new_orga_name"),
                domain: window.localStorage.getItem("new_orga_domain"),
                username: window.localStorage.getItem("new_orga_username"),
                default_lang: window.localStorage.getItem("new_orga_default_lang")
            };
        },
        computed: {
            lang() {
                return ISO6391.getNativeName(this.default_lang);
            }
        },
        methods: {
            action() {
                let data = {
                    name: this.name,
                    domain: this.domain,
                    username: this.username,
                    default_lang: this.default_lang
                };

                OrganizationService.createOrganization(data)
                    .then(() => {
                        window.localStorage.removeItem("new_orga_name");
                        window.localStorage.removeItem("new_orga_domain");
                        window.localStorage.removeItem("new_orga_username");
                        window.localStorage.removeItem("new_orga_default_lang");
                        this.$emit("next", data.domain);
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
            "title": "Confirmation",
            "button_submit": "Complete",
            "text": "The following free organization will be created:",
            "orga_name": "Organization name",
            "url": "URL",
            "username": "Your username",
            "language": "Standard language"
        },
        "de": {
            "title": "Bestätigung",
            "button_submit": "Abschließen",
            "text": "Die folgende kostenlose Organisation wird erstellt:",
            "orga_name": "Name der Organisation",
            "url": "URL",
            "username": "Dein Nutzername",
            "language": "Standardsprache"
        }
    }
</i18n>
