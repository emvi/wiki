<template>
    <div class="reg-form">
        <h3 class="reg-headline">
            {{$t("title")}}
        </h3>
        <div class="reg-icon">
            <div class="icon icon-tick"></div>
        </div>
        <p>
            {{$t("text")}}
        </p>
        <form v-on:submit.prevent="addMember">
            <field v-model="email" type="email" :label="$t('label_email')" :hint="$t('hint_email')" color="grey"></field>
        </form>
        <div class="modal--selection">	
            <div class="button--action button--pink button--active" v-for="(mail, index) in emails" :key="email+index">
                <div class="button--icon icon icon-mail"></div>
                <div class="button--label">{{mail}}</div>
                <div class="button--sub-icon icon icon-close-small" v-on:click="removeMember(index)"></div>
            </div>
        </div>
        <div class="row">
            <button v-on:click="invite" class="col-sm-12 reg-button--primary">{{$t("button_invite")}}</button>
        </div>
        <div class="row">
            <button v-on:click="open" class="col-sm-12 reg-button--secondary">{{$t("link_open")}}</button>
        </div>
    </div>
</template>

<script>
    import {OrganizationService} from "../../service";
    import field from "../html/field.vue";

    export default {
        components: {field},
        props: ["domain"],
        data() {
            return {
                email: "",
                emails: []
            };
        },
        methods: {
            addMember() {
                this.email = this.email.trim();

                if(this.email) {
                    this.emails.push(this.email);
                    this.email = "";
                }
            },
            removeMember(index) {
                this.emails.splice(index, 1);
            },
            invite() {
                this.email = this.email.trim();
                let match = /^\S+@\S+$/;

                if(this.email.length && match.test(this.email)) {
                    this.emails.push(this.email);
                    this.email = "";
                }

                OrganizationService.inviteMember(this.emails, this.domain)
                .then(() => {
                    this.emails = [];
                    this.$store.dispatch("showToast", {type: "green", message: this.$t("toast_invited")});
                });
            },
            open() {
                let url = EMVI_WIKI_FRONTEND_HOST.split("/");
                window.location = `${url[0]}//${this.domain}.${url[2]}`;
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "title": "Your organization was created!",
            "label_email": "Email address",
            "hint_email": "ENTER to confirm",
            "button_invite": "Invite members",
            "link_open": "See your organization",
            "text": "Invite your first members:",
            "toast_invited": "Invited members."
        },
        "de": {
            "title": "Deine Organisation wurde angelegt!",
            "label_email": "E-Mail-Adresse",
            "hint_email": "ENTER zum bestätigen",
            "button_invite": "Mitglieder einladen",
            "link_open": "Organisation öffnen",
            "text": "Lade die ersten Mitglieder ein:",
            "toast_invited": "Mitglieder eingeladen."
        }
    }
</i18n>
