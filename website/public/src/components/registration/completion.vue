<template>
    <div class="reg-form">
        <h3 class="reg-headline">
            {{$t("title")}}
        </h3>
        <form v-on:submit.prevent="action">
            <div class="form--element--full">
                <div class="form--element--checkbox no-select">
                    <checkbox v-model="terms_of_service" :error="validation['terms_of_service']" :autofocus="true">
                        {{$t('label_terms_1')}}
                        <a href="/terms" target="_blank" ref="noreferrer">{{$t('label_terms_2')}}</a>
                        {{$t('label_terms_3')}}
                    </checkbox>
                    <checkbox v-model="privacy" :error="validation['privacy']">
                        {{$t('label_privacy_1')}}
                        <a href="/privacy" target="_blank" ref="noreferrer">{{$t('label_privacy_2')}}</a>
                        {{$t('label_privacy_3')}}
                    </checkbox>
                    <checkbox :label="$t('label_marketing')" v-model="marketing" :error="validation['marketing']"></checkbox>
                </div>
            </div>
            <div class="recaptcha">
                <div id="recaptcha"></div>
            </div>
            <div class="row">
                <input type="submit" :value="$t('button_submit')" class="col-sm-12 reg-button--primary" />
            </div>
        </form>
    </div>
</template>

<script>
    import {AuthService} from "../../service";
    import checkbox from "../html/checkbox.vue";
    import {setCookie} from "../../util";

    export default {
        components: {checkbox},
        props: ["code"],
        data() {
            return {
                terms_of_service: false,
                privacy: false,
                marketing: false,
                recaptchaInit: false
            };
        },
        computed: {
            grecaptcha() {
                return this.$store.state.recaptcha.grecaptcha;
            }
        },
        watch: {
            grecaptcha(value) {
                if(value) {
                    this.renderRecaptcha();
                }
            }
        },
        mounted() {
            if(this.grecaptcha) {
                this.renderRecaptcha();
            }
        },
        methods: {
            action() {
                let recaptchaToken = this.grecaptcha.getResponse();
                AuthService.setRegistrationCompletion(this.code, this.terms_of_service, this.privacy, this.marketing, recaptchaToken)
                    .then(({access_token, expires_in, secure, domain}) => {
                        setCookie("access_token", access_token, expires_in, secure, domain);
                        this.$store.dispatch("loadUser");
                        this.$emit("success");
                    })
                    .catch(e => {
                        this.grecaptcha.reset();
                        this.setError(e);
                    })
            },
            renderRecaptcha() {
                if(!this.recaptchaInit) {
                    this.recaptchaInit = true;

                    this.grecaptcha.render('recaptcha', {
                        'sitekey': EMVI_WIKI_RECAPTCHA_CLIENT_SECRET
                    });
                }
            }
        }
    }
</script>

<style scoped>
    .recaptcha {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 0 0 20px 0;
    }
</style>

<i18n>
    {
        "en": {
            "title": "Complete",
            "label_terms_1": "I confirm that I have read the ",
            "label_terms_2": "Terms and Conditions",
            "label_terms_3": " and agree.",
            "label_privacy_1": "I confirm that I have read the ",
            "label_privacy_2": "Privacy Policy",
            "label_privacy_3": " and agree.",
            "label_marketing": "I would like to receive regular news on product updates and offers. The newsletter can be cancelled at any time. ",
            "button_submit": "Complete registration"
        },
        "de": {
            "title": "Abschließen",
            "label_terms_1": "Ich bestätige, dass ich die ",
            "label_terms_2": "Nutzungsbedingungen",
            "label_terms_3": " gelesen habe und diesen zustimme.",
            "label_privacy_1": "Ich bestätige, dass ich die ",
            "label_privacy_2": "Datenschutzerklärung",
            "label_privacy_3": " gelesen habe.",
            "label_marketing": "Ich möchte regelmäßige Nachrichten zu Produktverbesserungen und Angeboten erhalten. Der Newsletter kann jederzeit abbestellt werden.",
            "button_submit": "Registrierung abschließen"
        }
    }
</i18n>
