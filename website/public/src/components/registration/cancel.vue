<template>
    <div class="reg-form">
        <h3 class="h3">
            {{$t("title")}}
        </h3>
        <div v-if="!success">
            <p>
                {{$t("hint_cancel_registration")}}
            </p>
            <div class="row">
                <button class="col-sm-12 reg-button--primary" v-on:click="cancelRegistration">{{$t("button_yes")}}</button>
                <button class="col-sm-12 reg-button--secondary" v-on:click="back">{{$t("button_no")}}</button>
            </div>
        </div>
        <div v-if="success">
            <p>{{$t("hint_cancel_registration_success")}}</p>
            <p><router-link to="/" class="reg-button--secondary">{{$t("link_home")}}</router-link></p>
        </div>
    </div>
</template>

<script>
    import {AuthService} from "../../service";

    export default {
        props: ["code"],
        data() {
            return {
                success: false
            };
        },
        methods: {
            cancelRegistration() {
                if(this.code) {
                    AuthService.cancelRegistration(this.code)
                    .then(() => {
                        this.success = true;
                    })
                    .catch(e => {
                        this.setError(e);
                    });
                }
            },
            back() {
                this.$emit("cancel");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "title": "Cancel Registration?",
            "hint_cancel_registration": "Are you sure you want to cancel your registration? Your personal data will be deleted.",
            "hint_cancel_registration_success": "The registration has been cancelled. Your personal information has been deleted.",
            "button_yes": "Yes, cancel registration",
            "button_no": "No, continue",
            "link_home": "Back to Home"
        },
        "de": {
            "title": "Registrierung abbrechen?",
            "hint_cancel_registration": "Bist du sicher, dass du die Registrierung abbrechen möchtest? Deine personenbezogenen Daten werden damit gelöscht.",
            "hint_cancel_registration_success": "Die Registrierung wurde abgebrochen. Deine personenbezogenen Daten wurden gelöscht.",
            "button_yes": "Ja, Registrierung abbrechen",
            "button_no": "Nein, fortfahren",
            "link_home": "Zurück zur Startseite"
        }
    }
</i18n>
