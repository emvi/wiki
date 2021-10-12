<template>
    <div :class="{'widget': true, 'active': active && !isTouch, 'touch': isTouch}" ref="filter">
        <div class="headline no-select">
            <h2>
                {{title}}
                <emvi-shortcut :shortcut="$t('shift')" v-show="active" v-if="columns === 1">
                    <span>&#43;</span>
                    <span class="key">&uarr;</span>
                    <span class="key">&darr;</span>
                    {{$t("shortcut_sort")}}
                </emvi-shortcut>
                <emvi-shortcut :shortcut="$t('shift')" v-show="active" v-if="columns > 1">
                    <span>&#43;</span>
                    <span class="key">&uarr;</span>
                    <span class="key">&darr;</span>
                    <span class="key">&larr;</span>
                    <span class="key">&rarr;</span>
                    {{$t("shortcut_sort")}}
                </emvi-shortcut>
            </h2>
            <small v-if="!disableResultCount">{{count}} results</small>
        </div>
        <div class="widget-content" ref="results">
            <component :is="card"
                v-for="(result, index) in results"
                :key="result.id"
                :entity="result"
                :index="index"
                :active="index === selected && active"
                :scroll-area="$refs.results"
                :enter="enterEvent"
                :tab="tabEvent"
                :up="upEvent"
                :down="downEvent"
                v-on:preview="togglePreview"></component>
            <emvi-placeholder-card v-for="card in 3" :key="card" v-show="cancelToken"></emvi-placeholder-card>
            <emvi-empty :icon="icon" v-show="!cancelToken && results.length === 0"></emvi-empty>
            <infinite-loading v-on:infinite="search"></infinite-loading>
        </div>
    </div>
</template>

<script>
    import InfiniteLoading from "vue-infinite-loading";
    import {CancelToken} from "axios";
    import {mapGetters} from "vuex";
    import {scrollToTop} from "../../util";
    import emviArticleCard from "../cards/article-card.vue";
    import emviListCard from "../cards/list-card.vue";
    import emviMemberCard from "../cards/member-card.vue";
    import emviGroupCard from "../cards/group-card.vue";
    import emviTagCard from "../cards/tag-card.vue";
    import emviActivityCard from "../cards/activity-card.vue";
    import emviPlaceholderCard from "../cards/placeholder-card.vue";
    import emviEmpty from "./empty.vue";
    import emviShortcut from "../cmd/content/shortcut.vue";

    export default {
        components: {
            InfiniteLoading,
            emviArticleCard,
            emviListCard,
            emviMemberCard,
            emviGroupCard,
            emviTagCard,
            emviActivityCard,
            emviPlaceholderCard,
            emviEmpty,
            emviShortcut
        },
        props: {
            title: "",
            icon: "",
            columns: {default: 1},
            card: "",
            index: {default: 0},
            active: {default: false},
            performSearch: null,
            disableResultCount: {default: false},
            enter: {default: false},
            tab: {default: false},
            esc: {default: false},
            up: {default: false},
            down: {default: false}
        },
        data() {
            return {
                cancelToken: null,
                selected: 0,
                results: [],
                count: 0,
                previewActive: false,
                enterEvent: false,
                tabEvent: false,
                upEvent: false,
                downEvent: false
            };
        },
        computed: {
            ...mapGetters(["selection"])
        },
        watch: {
            active(active) {
                if(active) {
                    scrollToTop(this.$refs.filter);
                    this.$store.dispatch("select", this.selected);
                }
            },
            selection(selection) {
                if(this.active) {
                    this.selected = selection;
                }
            },
            enter(enter) {
                if(enter && this.active && this.results.length) {
                    this.enterEvent = true;
                    this.$nextTick(() => {
                        this.enterEvent = false;
                    });
                }
            },
            tab(tab) {
                if(tab && this.active && this.results.length) {
                    this.tabEvent = true;
                    this.$nextTick(() => {
                        this.tabEvent = false;
                    });
                }
            },
            esc(esc) {
                if(esc && this.active) {
                    this.$store.dispatch("select", 0);
                }
            },
            up(up) {
                if(up && this.active) {
                    if (!this.previewActive) {
                        this.$store.dispatch("selectPrevious", this.results.length - 1);
                    } else {
                        this.upEvent = true;
                        this.$nextTick(() => {
                            this.upEvent = false;
                        });
                    }
                }
            },
            down(down) {
                if(down && this.active) {
                    if(!this.previewActive) {
                        this.$store.dispatch("selectNext", this.results.length-1);
                    }
                    else {
                        this.downEvent = true;
                        this.$nextTick(() => {
                            this.downEvent = false;
                        });
                    }
                }
            }
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

                this.performSearch(filter, this.cancelToken)
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
                        this.showTechnicalError(e);
                    });
            },
            togglePreview({show, index}) {
                if(!this.active) {
                    this.$emit("toggleactive", this.index);
                }

                this.$nextTick(() => {
                    this.$store.dispatch("select", index);
                    this.previewActive = show;
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "shift": "Shift",
            "shortcut_sort": "Jump between widgets"
        },
        "de": {
            "shift": "Shift",
            "shortcut_sort": "Zwischen den Widgets springen"
        }
    }
</i18n>
