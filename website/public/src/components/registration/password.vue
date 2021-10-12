<template>
    <div class="reg-form">
        <h3 class="reg-headline">
            {{$t("title")}}
        </h3>
        <form v-on:submit.prevent="action">
            <div v-bind:class="{'password no-select': true, 'level-0': score === -1, 'level-1': score === 0, 'level-2': score === 1, 'level-3': score === 2, 'level-4': score === 3, 'level-5': score === 4}">
                <span v-show="score === -1">{{$t("label_no_pwd")}}</span>
                <span v-show="score !== -1">{{scores[score]}}</span>
                <div class="password-bar"></div>
            </div>
            <field type="password" :label="$t('label_pwd')" v-model="password" :error="validation['password']" :autofocus="true"></field>
            <field type="password" :label="$t('label_pwd_repeat')" v-model="password_repeat" :error="validation['password_repeat']"></field>
            <div class="row">
                <input type="submit" :value="$t('button_submit')" class="col-sm-12 reg-button--primary" />
            </div>
        </form>
    </div>
</template>

<script>
    import zxcvbn from "zxcvbn";
    import {AuthService} from "../../service";
    import field from "../html/field.vue";

    export default {
        components: {field},
        props: ["code"],
        data() {
            return {
                password: "",
                password_repeat: "",
                score: -1,
                scores: [
                    this.$t("score_very_weak"),
                    this.$t("score_weak"),
                    this.$t("score_medium"),
                    this.$t("score_strong"),
                    this.$t("score_very_strong")
                ]
            };
        },
        watch: {
            password(value) {
                if(!value) {
                    this.score = -1;
                }
                else {
                    this.score = zxcvbn(value).score;
                }
            }
        },
        methods: {
            action() {
                AuthService.setRegistrationPassword(this.code, this.password, this.password_repeat)
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
            "title": "Select Password",
            "label_pwd": "Password",
            "label_pwd_repeat": "Re-enter password",
            "label_no_pwd": "Choose a strong password",
            "button_submit": "Next Step",
            "score_very_weak": "Very weak",
            "score_weak": "Weak",
            "score_medium": "Medium",
            "score_strong": "Strong",
            "score_very_strong": "Very strong"
        },
        "de": {
            "title": "Passwort festlegen",
            "label_pwd": "Passwort",
            "label_pwd_repeat": "Passwort wiederholen",
            "label_no_pwd": "Wähle ein starkes Passwort",
            "button_submit": "Nächster Schritt",
            "score_very_weak": "Sehr schwach",
            "score_weak": "Schwach",
            "score_medium": "Mittel",
            "score_strong": "Stark",
            "score_very_strong": "Sehr stark"
        }
    }
</i18n>
