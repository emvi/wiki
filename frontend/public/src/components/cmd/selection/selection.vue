<template>
    <div>
        <emvi-cmd-submenu v-if="addView"
            icon="add"
            :view="addView"
            :index="0"
            :enter="enter">
            {{addText}}
        </emvi-cmd-submenu>
        <emvi-cmd-selection-search-field v-if="!disableSearch"
            ref="search"
            v-model="query"
            :placeholder="placeholder"
            :index="addView ? 1 : 0"
            v-on:next="selectNext"
            v-on:previous="selectPrevious"
            v-on:reset="resetQuery"
            v-on:focus="showDetails(false)"></emvi-cmd-selection-search-field>
        <component :is="element"
            v-for="(result, index) in results"
            :key="result.type+result.id"
            :entity="result"
            :index="index+rowOffset"
            :active="index+rowOffset === row"
            :details="detailsActive"
            :enter="enter"
            :esc="esc"
            :tab="tab"
            :del="del"
            :up="up"
            :down="down"
            v-on:update="update"
            v-on:swap="swap"
            v-on:remove="remove"
            v-on:details="showDetails"></component>
        <emvi-cmd-loading v-show="cancelToken"></emvi-cmd-loading>
        <emvi-cmd-emtpy v-show="!cancelToken && !results.length"></emvi-cmd-emtpy>
        <infinite-loading :identifier="resetLoading" v-on:infinite="search"></infinite-loading>
    </div>
</template>

<script>
    import InfiniteLoading from "vue-infinite-loading";
    import {mapGetters} from "vuex";
    import {CancelToken} from "axios";
    import {debounce, findIndexById} from "../../../util";
    import {updateSelectedRow} from "../util";
    import emviCmdSubmenu from "../content/submenu.vue";
    import emviCmdEmtpy from "../content/empty.vue";
    import emviCmdLoading from "../content/loading.vue";
    import emviCmdSelectionSearchField from "./form/search-field.vue";
    import emviCmdSelectionLanguage from "./entries/language.vue";
    import emviCmdSelectionMember from "./entries/member.vue";
    import emviCmdSelectionInvitation from "./entries/invitation.vue";
    import emviCmdSelectionClient from "./results/client.vue";
    import emviCmdSelectionBookmarkedArticle from "./entries/bookmarked-article.vue";
    import emviCmdSelectionBookmarkedList from "./entries/bookmarked-list.vue";
    import emviCmdSelectionObservedArticle from "./entries/observed-article.vue";
    import emviCmdSelectionObservedList from "./entries/observed-list.vue";
    import emviCmdSelectionObservedGroup from "./entries/observed-group.vue";
    import emviCmdSelectionListArticle from "./entries/list-article.vue";
    import emviCmdSelectionListMember from "./entries/list-member.vue";
    import emviCmdSelectionGroupMember from "./entries/group-member.vue";
    import emviCmdSelectionArticleMember from "./entries/article-member.vue";
    import emviCmdSelectionArticleHistory from "./entries/article-history.vue";
    import emviCmdSelectionArticleTranslation from "./entries/article-translation.vue";

    const searchDebounce = 250;

    export default {
        components: {
            InfiniteLoading,
            emviCmdSubmenu,
            emviCmdEmtpy,
            emviCmdLoading,
            emviCmdSelectionSearchField,
            emviCmdSelectionLanguage,
            emviCmdSelectionMember,
            emviCmdSelectionInvitation,
            emviCmdSelectionClient,
            emviCmdSelectionBookmarkedArticle,
            emviCmdSelectionBookmarkedList,
            emviCmdSelectionObservedArticle,
            emviCmdSelectionObservedList,
            emviCmdSelectionObservedGroup,
            emviCmdSelectionListArticle,
            emviCmdSelectionListMember,
            emviCmdSelectionGroupMember,
            emviCmdSelectionArticleMember,
            emviCmdSelectionArticleHistory,
            emviCmdSelectionArticleTranslation
        },
        props: {
            addView: {default: ""},
            addText: {default: ""},
            element: {default: ""},
            disableSearch: {default: false},
            performSearch: null,
            updateResults: {default: false},
            placeholder: {default: ""},
            enter: {default: false},
            tab: {default: false},
            del: {default: false},
            esc: {default: false},
            up: {default: false},
            down: {default: false}
        },
        data() {
            return {
                rowOffset: 2,
                searchDebounced: null,
                detailsActive: false,
                query: "",
                cancelToken: null,
                resetLoading: 0,
                results: [],
                count: 0
            };
        },
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            esc(esc) {
                if(esc && !this.detailsActive) {
                    if(!this.disableSearch && this.row !== this.rowOffset-1) {
                        this.$store.dispatch("selectRow", 1);
                    }
                    else {
                        this.$store.dispatch("popColumn");
                    }
                }
            },
            up(up) {
                if(up && !up.shiftKey && !this.detailsActive) {
                    this.selectPrevious();
                }
            },
            down(down) {
                if(down && !down.shiftKey && !this.detailsActive) {
                    this.selectNext();
                }
            },
            row(row) {
                this.setRow(row);
            },
            query() {
                this.searchDebounced();
            },
            updateResults() {
                this.resetResults();
            }
        },
        mounted() {
            this.searchDebounced = debounce(() => {
                this.resetResults();
            }, searchDebounce);

            this.setRowOffset();
        },
        methods: {
            setRowOffset() {
                if(!this.addView) {
                    this.rowOffset--;
                }

                if(this.disableSearch) {
                    this.rowOffset--;
                }

                this.$store.dispatch("selectRow", this.rowOffset-1);
            },
            focus() {
                if(this.row === this.rowOffset-1 && !this.disableSearch) {
                    this.$refs.search.$refs.input.focus();
                }
                else {
                    this.focusCmdInput();
                }
            },
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
                        this.showTechnicalError(e);
                    });
            },
            update() {
                this.resetResults();
            },
            swap({id, direction}) {
                let index = findIndexById(this.results, id);

                if(index > -1) {
                    let swap = 0;

                    if(direction < 0) {
                        swap = index-1;

                        if(swap < 0) {
                            swap = this.results.length-1;
                        }
                    }
                    else {
                        swap = index+1;

                        if(swap > this.results.length-1) {
                            swap = 0;
                        }
                    }

                    let result = this.results[index];
                    this.results[index] = this.results[swap];
                    this.results[swap] = result;
                    this.$store.dispatch("selectRow", swap+1);
                }
            },
            remove(id) {
                let index = findIndexById(this.results, id);

                if(index > -1) {
                    this.results.splice(index, 1);
                }

                this.$store.dispatch("selectPreviousRow");
            },
            resetResults() {
                if(this.cancelToken) {
                    this.cancelToken.cancel();
                }

                this.cancelToken = null;
                this.resetLoading++;
                this.results = [];
                this.count = 0;
            },
            resetQuery() {
                if(this.query === "") {
                    this.$store.dispatch("popColumn");
                }
                else {
                    this.query = "";
                }
            },
            showDetails(detailsActive) {
                this.detailsActive = detailsActive;
            },
            setRow(row) {
                updateSelectedRow(row, this.results.length+this.rowOffset, this.$store);
                this.focus();
            },
            selectNext() {
                this.$store.dispatch("selectNextRow");
            },
            selectPrevious() {
                this.$store.dispatch("selectPreviousRow");
            }
        }
    }
</script>
