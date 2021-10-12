<template>
    <div class="item">
        <p>
            <span v-html="title"></span>
            <emvi-private-label v-if="entity.private"></emvi-private-label>
            <emvi-external-label v-if="entity.client_access"></emvi-external-label>
        </p>
        <small>
            <template v-if="entity.type === 'article'">
                <span v-if="entity.views === 1">{{entity.views}} {{$t("view")}}</span>
                <span v-if="entity.views !== 1">{{entity.views}} {{$t("views")}}</span>
                <span class="dot">·</span>
                <template v-if="entity.published"><span>{{entity.published | moment("ll")}}</span></template>
                <template v-if="!entity.published"><span>{{$t("unpublished")}}</span></template>
                <span class="dot">·</span>
                <span>{{$t("edited_before")}} {{entity.latest_article_content.mod_time | moment("from", "now")}} {{$t("edited_after")}}</span>
            </template>
            <template v-if="entity.type === 'list'">
                <span v-if="entity.articles === 1">{{entity.articles}} {{$t("article")}}</span>
                <span v-if="entity.articles !== 1">{{entity.articles}} {{$t("articles")}}</span>
                <span class="dot">·</span>
                <span>{{entity.name.info}}</span>
                <span class="dot" v-if="entity.name.info">·</span>
                <span>{{entity.def_time | moment("ll")}}</span>
                <span class="dot">·</span>
                <span>{{$t("edited_before")}} {{entity.mod_time | moment("from", "now")}} {{$t("edited_after")}}</span>
            </template>
            <template v-if="entity.type === 'user'">
                <span>{{entity.organization_member.info}}</span>
                <span class="dot" v-if="entity.organization_member.info">·</span>
                <span>{{entity.organization_member.username}}</span>
                <span class="dot">·</span>
                <span>{{entity.email}}</span>
            </template>
            <template v-if="entity.type === 'tag'">
                <span v-if="entity.usages === 1">{{entity.usages}} {{$t("usage")}}</span>
                <span v-if="entity.usages !== 1">{{entity.usages}} {{$t("usages")}}</span>
                <span class="dot">·</span>
                <span>{{entity.def_time | moment("ll")}}</span>
            </template>
            <template v-if="entity.type === 'group'">
                <span v-if="entity.member === 1">{{entity.member}} {{$t("member")}}</span>
                <span v-if="entity.member !== 1">{{entity.member}} {{$t("members")}}</span>
                <span class="dot">·</span>
                <span>{{entity.info}}</span>
                <span class="dot" v-if="entity.info">·</span>
                <span>{{entity.def_time | moment("ll")}}</span>
                <span class="dot">·</span>
                <span>{{$t("edited_before")}} {{entity.mod_time | moment("from", "now")}} {{$t("edited_after")}}</span>
            </template>
        </small>
    </div>
</template>

<script>
    import {markInText} from "../util";
    import emviPrivateLabel from "../../labels/private.vue";
    import emviExternalLabel from "../../labels/external.vue";

    export default {
        components: {emviPrivateLabel, emviExternalLabel},
        props: {
            entity: {},
            query: {default: ""}
        },
        computed: {
            title() {
                let name;

                switch(this.entity.type) {
                    case "article":
                        name = this.entity.latest_article_content.title;
                        break;
                    case "list":
                        name = this.entity.name.name;
                        break;
                    case "user":
                        name = `${this.entity.firstname} ${this.entity.lastname}`;
                        break;
                    case "group":
                        name = this.entity.name;
                        break;
                    default:
                        name = this.entity.name;
                }

                return markInText(this.query, name);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "article": "article",
            "articles": "articles",
            "member": "member",
            "members": "members",
            "usage": "usage",
            "usages": "usages",
            "view": "view",
            "views": "views",
            "unpublished": "Unpublished",
            "edited_before": "edited",
            "edited_after": " "
        },
        "de": {
            "article": "Artikel",
            "articles": "Artikel",
            "member": "Mitglied",
            "members": "Mitglieder",
            "usage": "Verwendung",
            "usages": "Verwendungen",
            "view": "Aufruf",
            "views": "Aufrufe",
            "unpublished": "Unveröffentlicht",
            "edited_before": " ",
            "edited_after": "bearbeitet"
        }
    }
</i18n>
