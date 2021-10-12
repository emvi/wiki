<template>
    <emvi-layout v-on:enter="enter" v-on:tab="tab" v-on:esc="esc" v-on:up="up" v-on:down="down">
        <emvi-layout-narrow>
            <template v-if="!notFound">
                <emvi-profile-header :entity="member">{{title}}</emvi-profile-header>
                <small>{{$t("info")}}</small>
                <p v-if="member.organization_member.info">{{member.organization_member.info}}</p>
                <p v-if="!member.organization_member.info">-</p>
                <small>{{$t("username")}}</small>
                <p>{{member.organization_member.username}}</p>
                <small>{{$t("email")}}</small>
                <p>{{member.email}}</p>
                <div v-if="member.id">
                    <emvi-widget :title="$t('title_activities')"
                        card="emvi-activity-card"
                        icon="notification"
                        :perform-search="loadActivities"
                        disable-result-count="true"
                        :index="0"
                        :active="activeFilter === 0"
                        :enter="enterEvent"
                        :tab="tabEvent"
                        :esc="escEvent"
                        :up="upEvent"
                        :down="downEvent"
                        v-on:toggleactive="toggleActiveFilter"></emvi-widget>
                    <emvi-widget :title="$t('title_articles')"
                        card="emvi-article-card"
                        icon="article"
                        :perform-search="loadArticles"
                        :index="1"
                        :active="activeFilter === 1"
                        :enter="enterEvent"
                        :tab="tabEvent"
                        :esc="escEvent"
                        :up="upEvent"
                        :down="downEvent"
                        v-on:toggleactive="toggleActiveFilter"></emvi-widget>
                    <emvi-widget :title="$t('title_lists')"
                        card="emvi-list-card"
                        icon="list"
                        :perform-search="loadLists"
                        :index="2"
                        :active="activeFilter === 2"
                        :enter="enterEvent"
                        :tab="tabEvent"
                        :esc="escEvent"
                        :up="upEvent"
                        :down="downEvent"
                        v-on:toggleactive="toggleActiveFilter"></emvi-widget>
                    <emvi-widget :title="$t('title_groups')"
                        card="emvi-group-card"
                        icon="group"
                        :perform-search="loadGroups"
                        :index="3"
                        :active="activeFilter === 3"
                        :enter="enterEvent"
                        :tab="tabEvent"
                        :esc="escEvent"
                        :up="upEvent"
                        :down="downEvent"
                        v-on:toggleactive="toggleActiveFilter"></emvi-widget>
                </div>
            </template>
            <emvi-not-found v-if="notFound"></emvi-not-found>
        </emvi-layout-narrow>
    </emvi-layout>
</template>

<script>
    import {UserService, FeedService, SearchService} from "../service";
    import {FilterMixin} from "./filter";
    import {setPageTitle} from "./title";
    import {emviLayout, emviLayoutNarrow, emviProfileHeader, emviNotFound, emviWidget} from "../components";

    const activityLimit = 10;

    export default {
        mixins: [FilterMixin],
        components: {emviLayout, emviLayoutNarrow, emviProfileHeader, emviNotFound, emviWidget},
        data() {
            return {
                member: {organization_member: {}},
                notFound: false,
                maxFilterIndex: 3
            };
        },
        computed: {
            title() {
                return `${this.member.firstname} ${this.member.lastname}`;
            },
            picture() {
                return this.member.picture;
            }
        },
        mounted() {
            this.loadMember();
        },
        methods: {
            loadMember() {
                UserService.getUser(undefined, this.$route.params.username)
                    .then(member => {
                        this.member = member;
                        setPageTitle(`${this.member.firstname} ${this.member.lastname}`);
                    })
                    .catch(() => {
                        this.notFound = true;
                    });
            },
            loadActivities(filter, cancelToken) {
                filter.limit = activityLimit;
                filter.user = this.member.id;

                return new Promise((resolve, reject) => {
                    FeedService.getFeed(filter, cancelToken)
                        .then(results => {
                            resolve({results: results, count: results.length === 0 ? 0 : 9999999});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            loadArticles(filter, cancelToken) {
                filter.authors = this.member.id;

                return new Promise((resolve, reject) => {
                    SearchService.findArticles("", filter, cancelToken)
                        .then(({results, count}) => {
                            resolve({results, count});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            loadLists(filter, cancelToken) {
                filter.user_ids = this.member.id;

                return new Promise((resolve, reject) => {
                    SearchService.findArticleLists("", filter, cancelToken)
                        .then(({results, count}) => {
                            resolve({results, count});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            loadGroups(filter, cancelToken) {
                filter.user_ids = this.member.id;

                return new Promise((resolve, reject) => {
                    SearchService.findUserGroups("", filter, cancelToken)
                        .then(({results, count}) => {
                            resolve({results, count});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "info": "Info",
            "username": "Username",
            "email": "Email",
            "title_activities": "Activities",
            "title_articles": "Articles",
            "title_lists": "Lists",
            "title_groups": "Groups"
        },
        "de": {
            "info": "Info",
            "username": "Nutzername",
            "email": "E-Mail",
            "title_activities": "Aktivitäten",
            "title_articles": "Artikel",
            "title_lists": "Listen",
            "title_groups": "Gruppen"
        }
    }
</i18n>

