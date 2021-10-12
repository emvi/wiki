<template>
    <card show-right="true">
        <template slot="icon">
            <i class="icon large-icon--grey icon-mail card--content--icon" v-on:click="join()"></i>
        </template>
        <template>
            <div class="card--content--center--info">{{$t("label_open_invitation")}}</div>
            <div class="card--content--center--item" v-on:click="join()">{{invitation.organization_name}}</div>
        </template>
        <template slot="right">
            <btn icon="trash" color="red" v-on:click="$emit('delete', invitation.id)">{{$t("button_delete")}}</btn>
            <btn icon="arrow" v-on:click="join()">{{$t("button_join")}}</btn>
        </template>
    </card>
</template>

<script>
    import card from "./card.vue";
    import btn from "../html/btn.vue";

    export default {
        components: {
            card,
            btn
        },
        props: ["invitation"],
        methods: {
            join() {
                this.$router.push(`/join?code=${this.invitation.code}&email=${this.invitation.email}`);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "label_open_invitation": "Open invitation",
            "button_join": "Join now",
            "button_delete": "Delete"
        },
        "de": {
            "label_open_invitation": "Offene Einladung",
            "button_join": "Jetzt beitreten",
            "button_delete": "Löschen"
        }
    }
</i18n>
