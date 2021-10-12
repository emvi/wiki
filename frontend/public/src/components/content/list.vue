<template>
    <emvi-layout v-on:enter="enter" v-on:tab="tab" v-on:esc="esc" v-on:up="up" v-on:down="down">
        <emvi-layout-narrow>
            <template v-if="!notFound">
                <div class="item-headline" v-if="hasHeadline">
                    <slot name="headline"></slot>
                </div>
                <emvi-page-title :icon="titleIcon" :results="count" :disable-result-count="disableResultCount">{{title}}</emvi-page-title>
                <div v-if="hasSlot">
                    <slot></slot>
                </div>
                <emvi-text-search v-on:update="update" v-if="!disableSearch"></emvi-text-search>
                <component :is="card"
                    v-for="(result, index) in results"
                    :key="result.id"
                    :entity="result"
                    :index="index+rowOffset"
                    :active="index+rowOffset === selection"
                    :enter="enterEvent"
                    :tab="tabEvent"
                    :up="upEvent"
                    :down="downEvent"
                    :query="query"
                    v-on:preview="togglePreview"></component>
                <emvi-placeholder-card v-for="card in 3" :key="card" v-show="cancelToken"></emvi-placeholder-card>
                <emvi-empty :icon="icon" v-show="!cancelToken && results.length === 0"></emvi-empty>
                <infinite-loading forceUseInfiniteWrapper="true" :identifier="resetLoading" v-on:infinite="search"></infinite-loading>
            </template>
            <emvi-not-found v-if="notFound"></emvi-not-found>
        </emvi-layout-narrow>
    </emvi-layout>
</template>

<script>
    import InfiniteLoading from "vue-infinite-loading";
    import {CancelToken} from "axios";
    import {mapGetters} from "vuex";
    import {debounce} from "../../util";
    import emviArticleCard from "../cards/article-card.vue";
    import emviListCard from "../cards/list-card.vue";
    import emviMemberCard from "../cards/member-card.vue";
    import emviGroupCard from "../cards/group-card.vue";
    import emviTagCard from "../cards/tag-card.vue";
    import emviActivityCard from "../cards/activity-card.vue";
    import emviEmpty from "./empty.vue";
    import emviLayout from "../layout/layout.vue";
    import emviLayoutNarrow from "../layout/layout-narrow.vue";
    import emviPageTitle from "./page-title.vue";
    import emviPlaceholderCard from "../cards/placeholder-card.vue";
    import emviTextSearch from "./text-search.vue";
    import emviNotFound from "../layout/not-found.vue";

    const searchDebounce = 250;

    export default {
        components: {
            InfiniteLoading,
            emviLayout,
            emviLayoutNarrow,
            emviTextSearch,
            emviArticleCard,
            emviListCard,
            emviMemberCard,
            emviGroupCard,
            emviTagCard,
            emviActivityCard,
            emviPageTitle,
            emviPlaceholderCard,
            emviEmpty,
            emviNotFound
        },
        props: ["title", "icon", "titleIcon", "card", "performSearch", "disableResultCount", "disableSearch"],
        data() {
            return {
                searchDebounced: null,
                query: "",
                cancelToken: null,
                resetLoading: 0,
                rowOffset: 1,
                results: [],
                count: 0,
                previewActive: false,
                notFound: false,
                enterEvent: false,
                tabEvent: false,
                upEvent: false,
                downEvent: false
            };
        },
        computed: {
            ...mapGetters(["selection", "metaUpdate"]),
            hasHeadline() {
                return !!this.$slots.headline;
            },
            hasSlot() {
                return !!this.$slots.default;
            }
        },
        watch: {
            selection(selection) {
                if(selection === 0) {
                    window.scrollTo(0, 0);
                }
            },
            metaUpdate() {
                if(this.$store.state.page.meta.get("updateList")) {
                    this.resetResults();
                    this.$store.dispatch("setMeta", {key: "updateList", value: undefined});
                }
            }
        },
        mounted() {
            if(this.disableSearch) {
                this.rowOffset = 0;
            }

            this.searchDebounced = debounce(() => {
                this.resetResults();
            }, searchDebounce);
        },
        beforeRouteLeave(from, to, next) {
            this.$store.dispatch("select", 0);
            next();
        },
        methods: {
            search($state) {
                if(this.cancelToken) {
                    this.cancelToken.cancel();
                    this.cancelToken = null;
                }

                this.cancelToken = CancelToken.source();
                let filter = {
                    offset: this.results.length
                };

                this.performSearch(this.query || "", filter, this.cancelToken)
                    .then(({results, count}) => {
                        this.results = this.results.concat(results);
                        this.count = count;
                        this.cancelToken = null;

                        if(this.results.length >= this.count) {
                            $state.complete();
                        }
                        else {
                            $state.loaded();
                        }
                    })
                    .catch(e => {
                        if(e.status && e.status === 404) {
                            this.notFound = true;
                            this.$store.dispatch("setMeta", {key: "notFound", value: true});
                        }
                        else {
                            this.showTechnicalError(e);
                        }
                    });
            },
            update(query) {
                this.query = query;
                this.searchDebounced();
            },
            resetResults() {
                if(this.cancelToken) {
                    this.cancelToken.cancel();
                }

                this.cancelToken = null;
                this.resetLoading++;
                this.results = [];
                this.count = 0;
                this.$store.dispatch("select", 0);
            },
            enter() {
                if(this.results.length && this.selection > this.rowOffset-1) {
                    this.enterEvent = true;
                    this.$nextTick(() => {
                        this.enterEvent = false;
                    });
                }
            },
            tab() {
                if(this.results.length && this.selection > this.rowOffset-1) {
                    this.tabEvent = true;
                    this.$nextTick(() => {
                        this.tabEvent = false;
                    });
                }
            },
            esc() {
                this.$store.dispatch("select", 0);
            },
            up() {
                if(!this.previewActive) {
                    this.$store.dispatch("selectPrevious", this.results.length-1+this.rowOffset);
                }
                else {
                    this.upEvent = true;
                    this.$nextTick(() => {
                        this.upEvent = false;
                    });
                }
            },
            down() {
                if(!this.previewActive) {
                    this.$store.dispatch("selectNext", this.results.length-1+this.rowOffset);
                }
                else {
                    this.downEvent = true;
                    this.$nextTick(() => {
                        this.downEvent = false;
                    });
                }
            },
            togglePreview({show, index}) {
                this.$store.dispatch("select", index);
                this.previewActive = show;
            }
        }
    }
</script>
