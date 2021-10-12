<template>
    <div ref="header">
        <div class="profile-header no-select">
            <emvi-cmd-avatar size="96" icon="organization" color="blue" :entity="organizationEntity"></emvi-cmd-avatar>
            <h1>{{organizationEntity.name}}</h1>
        </div>
        <div class="bullets">
            <span class="blue-100 bg-blue-10">{{organizationStatistics.article_count}} {{$t("articles")}}</span>
            <span class="green-100 bg-green-10">{{organizationStatistics.list_count}} {{$t("lists")}}</span>
            <span class="pink-100 bg-pink-10">{{organizationStatistics.member_count}} {{$t("members")}}</span>
            <span class="purple-100 bg-purple-10">{{organizationStatistics.group_count}} {{$t("groups")}}</span>
            <span class="orange-100 bg-orange-10">{{organizationStatistics.tag_count}} {{$t("tags")}}</span>
            <span class="bg-grey-10-to-grey-80">{{organizationStatistics.storage_usage | size}} {{$t("of")}} {{organizationStatistics.max_storage}} GB</span>
        </div>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {getSizeFromBytes, isEmptyObject, scrollIntoViewArea} from "../../../util";
    import emviCmdAvatar from "../content/avatar.vue";

    export default {
        components: {emviCmdAvatar},
        filters: {
            size(value) {
                return getSizeFromBytes(value);
            }
        },
        computed: {
            ...mapGetters(["row", "organization", "organizationStatistics"]),
            organizationEntity() {
                let entity = this.organization;

                // set type to user in case the picture is set to show it as the avatar
                if(entity.picture) {
                    entity.type = "user";
                }

                return entity;
            }
        },
        watch: {
            row(row) {
                if(row === 0) {
                    scrollIntoViewArea(this.$refs.header, document.getElementById("cmd-results"));
                }
            }
        },
        mounted() {
            this.loadStatistics();
        },
        methods: {
            loadStatistics() {
                if(isEmptyObject(this.organizationStatistics)) {
                    this.$store.dispatch("loadOrganizationStatistics");
                }
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "articles": "Articles",
            "lists": "Lists",
            "members": "Members",
            "groups": "Groups",
            "tags": "Tags",
            "of": "of"
        },
        "de": {
            "articles": "Artikel",
            "lists": "Listen",
            "members": "Mitglieder",
            "groups": "Gruppen",
            "tags": "Tags",
            "of": "von"
        }
    }
</i18n>
