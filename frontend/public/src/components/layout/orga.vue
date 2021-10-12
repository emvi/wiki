<template>
    <div class="orga cursor-pointer" v-on:click="toStart">
        <div class="avatar size-48">
            <i class="icon icon-organization blue-100 bg-blue-10" v-if="!orgaPicture"></i>
            <img :src="orgaPicture" alt="" v-if="orgaPicture" />
        </div>
        <h4>{{orga}}</h4>
        <span class="label" v-on:click.stop="toBillingPage" v-if="!isExpert">{{$t("entry")}}</span>
        <span class="label expert" v-on:click.stop="toBillingPage" v-else>{{$t("expert")}}</span>
        <i class="icon icon-chevron grey-50" :title="$t('title')" v-on:click="toOrganizationOverview"></i>
    </div>
</template>

<script>
    export default {
        computed: {
            orgaPicture() {
                return this.$store.state.user.organization.picture;
            },
            orga() {
                return this.$store.state.user.organization.name;
            }
        },
        methods: {
            toStart() {
                if(this.$route.path !== "/") {
                    this.$store.dispatch("resetCmd");
                    this.$router.push("/");
                }
            },
            toOrganizationOverview() {
                window.location = `${EMVI_WIKI_WEBSITE_HOST}/organizations`;
            },
            toBillingPage() {
                if(this.$route.path !== "/billing") {
                    this.$router.push("/billing");
                }
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "title": "Organization overview",
            "entry": "Entry",
            "expert": "Expert"
        },
        "de": {
            "title": "Organisationsübersicht",
            "entry": "Entry",
            "expert": "Expert"
        }
    }
</i18n>
