<template>
    <div>
        <emvi-mobile-head v-if="isMobile"></emvi-mobile-head>
        <emvi-layout v-on:enter="enter" v-on:tab="tab" v-on:esc="esc" v-on:up="up" v-on:down="down" v-on:left="left" v-on:right="right" hide-navigation="true" hide-actions="true">
            <emvi-layout-tiles>
                <template slot="title">
                    <h1>{{$t("title")}}</h1>
                </template>
                <emvi-widget :title="$t('title_pinned_articles')"
                    card="emvi-article-card"
                    icon="article"
                    :columns="filterColumns"
                    :perform-search="loadPinnedArticles"
                    :index="0"
                    :active="activeFilter === 0"
                    :enter="enterEvent"
                    :tab="tabEvent"
                    :esc="escEvent"
                    :up="upEvent"
                    :down="downEvent"
                    v-on:toggleactive="toggleActiveFilter"></emvi-widget>
                <emvi-widget :title="$t('title_pinned_lists')"
                    card="emvi-list-card"
                    icon="list"
                    :columns="filterColumns"
                    :perform-search="loadPinnedLists"
                    :index="1"
                    :active="activeFilter === 1"
                    :enter="enterEvent"
                    :tab="tabEvent"
                    :esc="escEvent"
                    :up="upEvent"
                    :down="downEvent"
                    v-on:toggleactive="toggleActiveFilter"></emvi-widget>
                <emvi-widget :title="$t('title_last_updated')"
                    card="emvi-article-card"
                    icon="article"
                    :columns="filterColumns"
                    :perform-search="loadLastUpdatedArticles"
                    :index="2"
                    :active="activeFilter === 2"
                    :enter="enterEvent"
                    :tab="tabEvent"
                    :esc="escEvent"
                    :up="upEvent"
                    :down="downEvent"
                    v-on:toggleactive="toggleActiveFilter"></emvi-widget>
                <emvi-widget :title="$t('title_in_progress')"
                    card="emvi-article-card"
                    icon="article"
                    :columns="filterColumns"
                    :perform-search="loadInProgressArticles"
                    :index="3"
                    :active="activeFilter === 3"
                    :enter="enterEvent"
                    :tab="tabEvent"
                    :esc="escEvent"
                    :up="upEvent"
                    :down="downEvent"
                    v-on:toggleactive="toggleActiveFilter"></emvi-widget>
            </emvi-layout-tiles>
        </emvi-layout>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {PinnedService, SearchService} from "../service";
    import {FilterMixin} from "./filter";
    import {TitleMixin} from "./title";
    import {emviLayout, emviLayoutTiles, emviWidget, emviMobileHead} from "../components";

    export default {
        mixins: [TitleMixin, FilterMixin],
        components: {emviLayout, emviLayoutTiles, emviWidget, emviMobileHead},
        data() {
            return {
                maxFilterIndex: 3,
                filterColumns: 2
            };
        },
        computed: {
            ...mapGetters(["user"])
        },
        watch: {
            $mq(mq) {
                if(mq > 1200) {
                    this.filterColumns = 2;
                }
                else {
                    this.filterColumns = 1;
                }
            }
        },
        methods: {
            loadPinnedArticles(filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    PinnedService.getPinned(true, false, filter.offset, 0, cancelToken)
                        .then(r => {
                            resolve({results: r.articles, count: r.articlesCount});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            loadPinnedLists(filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    PinnedService.getPinned(false, true, 0, filter.offset, cancelToken)
                        .then(r => {
                            resolve({results: r.lists, count: r.listsCount});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            loadLastUpdatedArticles(filter, cancelToken) {
                filter.updated_start = this.$moment().subtract(30, "d").format("YYYY-MM-DD");
                filter.sort_updated = "desc";

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
            loadInProgressArticles(filter, cancelToken) {
                filter.wip = true;
                filter.authors = this.user.id;

                return new Promise((resolve, reject) => {
                    SearchService.findArticles("", filter, cancelToken)
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
            "title": "Start",
            "title_pinned_articles": "Pinned Articles",
            "title_pinned_lists": "Pinned Lists",
            "title_last_updated": "Last Updated",
            "title_in_progress": "In Progress"
        },
        "de": {
            "title": "Start",
            "title_pinned_articles": "Angepinnte Artikel",
            "title_pinned_lists": "Angepinnte Listen",
            "title_last_updated": "Zuletzt aktualisiert",
            "title_in_progress": "In Bearbeitung"
        }
    }
</i18n>
