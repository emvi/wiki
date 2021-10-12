<template>
    <div class="reg-form">
        <h3 class="reg-headline">
            {{$t("title")}}
        </h3>
        <a class="sso-button google" v-bind:href="googleURL">
            <img src="static/img/icon-google.svg" alt="Google">
            {{$t("label_sso_google")}}
        </a>
        <a class="sso-button slack" v-bind:href="slackURL">
            <img src="static/img/icon-slack.svg" alt="Slack">
            {{$t("label_sso_slack")}}
        </a>
        <a class="sso-button github" v-bind:href="githubURL">
            <img src="static/img/icon-github.svg" alt="GitHub">
            {{$t("label_sso_github")}}
        </a>
        <a class="sso-button github" v-bind:href="microsoftURL">
            <img src="static/img/icon-microsoft.svg" alt="Microsoft">
            {{$t("label_sso_microsoft")}}
        </a>
        <div class="spacer-16"></div>
        <div class="or no-select">
            <div class="or--line"></div>
            <div class="or--text">{{$t("or")}}</div>
            <div class="or--line"></div>
        </div>
        <form v-on:submit.prevent="action" v-if="!success">
            <field type="email" :label="$t('label_email')" v-model="email" :error="validation['email']" :autofocus="true"></field>
            <div class="row">
                <input type="submit" :value="$t('button_submit')" class="col-sm-12 reg-button--primary" />
            </div>
        </form>
        <p v-if="success">
            {{$t("hint_success")}}
        </p>
    </div>
</template>

<script>
    import {AuthService} from "../../service";
    import field from "../html/field.vue";

    export default {
        components: {field},
        props: ["initialemail"],
        data() {
            return {
                email: "",
                success: false
            };
        },
        computed: {
            githubURL() {
                return `https://github.com/login/oauth/authorize?client_id=${EMVI_WIKI_GITHUB_CLIENT_ID}&redirect_uri=${EMVI_WIKI_AUTH_HOST}/auth/sso/github&scope=read:user user:email`;
            },
            slackURL() {
                return `https://slack.com/oauth/authorize?scope=identity.basic identity.email identity.avatar&client_id=${EMVI_WIKI_SLACK_CLIENT_ID}&redirect_uri=${EMVI_WIKI_AUTH_HOST}/auth/sso/slack`
            },
            googleURL() {
                return `https://accounts.google.com/o/oauth2/v2/auth?client_id=${EMVI_WIKI_GOOGLE_CLIENT_ID}&redirect_uri=${EMVI_WIKI_AUTH_HOST}/auth/sso/google&response_type=code&scope=https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile`
            },
            microsoftURL() {
                return `https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=${EMVI_WIKI_MICROSOFT_CLIENT_ID}&redirect_uri=${EMVI_WIKI_AUTH_HOST}/auth/sso/microsoft&response_type=code&scope=user.read`
            }
        },
        mounted() {
            if(this.initialemail) {
                this.email = this.initialemail;
            }
        },
        methods: {
            action() {
                AuthService.register(this.email)
                .then(() => {
                    this.success = true;
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
            "title": "Create account",
            "label_email": "Email address",
            "button_submit": "Create your Emvi account",
            "label_sso_google": "Sign up with Google",
            "label_sso_slack": "Sign up with Slack",
            "label_sso_github": "Sign up with GitHub",
            "label_sso_microsoft": "Sign up with Microsoft",
            "hint_success": "Success! Please check your mailbox to confirm your email address.",
            "or": "or"
        },
        "de": {
            "title": "Konto anlegen",
            "label_email": "E-Mail-Adresse",
            "button_submit": "Emvi-Konto erstellen",
            "label_sso_google": "Mit Google anmelden",
            "label_sso_slack": "Mit Slack anmelden",
            "label_sso_github": "Mit GitHub anmelden",
            "label_sso_microsoft": "Mit Microsoft anmelden",
            "hint_success": "Erfolg! Bitte überprüfe dein E-Mail-Postfach, um deine E-Mail-Adresse zu bestätigen.",
            "or": "oder"
        }
    }
</i18n>
