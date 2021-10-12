<template>
    <div>
        <emvi-cmd-selection-search-field ref="search"
            v-model="query"
            :placeholder="placeholder"
            :index="0"
            v-on:next="selectNext"
            v-on:previous="selectPrevious"
            v-on:reset="resetQuery"></emvi-cmd-selection-search-field>
        <component :is="element"
            v-for="(result, index) in results"
            disable-remove="true"
            :key="result.type+result.id"
            :entity="result"
            :index="index+1"
            :active="index+1 === row"
            :enter="enter"
            :up="up"
            :down="down"
            v-on:add="add"></component>
        <emvi-cmd-loading v-show="cancelToken"></emvi-cmd-loading>
        <emvi-cmd-emtpy v-show="!cancelToken && !results.length"></emvi-cmd-emtpy>
        <infinite-loading :identifier="resetLoading" v-on:infinite="search"></infinite-loading>
        <h5 v-show="selected.length">{{$t("selected")}}</h5>
        <component :is="element"
            v-for="(result, index) in selected"
            disable-add="true"
            :key="result.type+result.id"
            :entity="result"
            :index="index+results.length+1"
            :active="index+results.length+1 === row"
            :del="del"
            :up="up"
            :down="down"
            v-on:remove="remove"></component>
        <emvi-cmd-button :index="results.length+selected.length+1"
            :label="button"
            :icon="icon"
            v-on:next="selectNext"
            v-on:previous="selectPrevious"
            v-on:enter="save"
            v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import InfiniteLoading from "vue-infinite-loading";
    import {mapGetters} from "vuex";
    import {CancelToken} from "axios";
    import {debounce, findIndexById, removeFromListByIds} from "../../../util";
    import {updateSelectedRow} from "../util";
    import emviCmdEmtpy from "../content/empty.vue";
    import emviCmdLoading from "../content/loading.vue";
    import emviCmdButton from "../form/button.vue";
    import emviCmdSelectionSearchField from "./form/search-field.vue";
    import emviCmdSelectionArticle from "./results/article.vue";
    import emviCmdSelectionUser from "./results/user.vue";
    import emviCmdSelectionGroup from "./results/group.vue";
    import emviCmdSelectionUserGroup from "./results/user-group.vue";

    const searchDebounce = 250;

    export default {
        components: {
            InfiniteLoading,
            emviCmdEmtpy,
            emviCmdLoading,
            emviCmdButton,
            emviCmdSelectionSearchField,
            emviCmdSelectionArticle,
            emviCmdSelectionUser,
            emviCmdSelectionGroup,
            emviCmdSelectionUserGroup
        },
        props: {
            element: {default: ""},
            preselected: null,
            performSearch: null,
            placeholder: {default: ""},
            button: {default: ""},
            icon: {default: ""},
            enter: {default: false},
            tab: {default: false},
            del: {default: false},
            esc: {default: false},
            up: {default: false},
            down: {default: false}
        },
        data() {
            return {
                lock: false, // prevent adding multiple results at once
                searchDebounced: null,
                query: "",
                cancelToken: null,
                resetLoading: 0,
                results: [],
                resultCount: 0,
                count: 0,
                selected: []
            };
        },
        computed: {
            ...mapGetters(["row"])
        },
        watch: {
            tab(tab) {
                if(tab) {
                    if(tab.shiftKey) {
                        this.selectPrevious();
                    }
                    else {
                        this.selectNext();
                    }
                }
            },
            del(del) {
                if(del) {
                    del.preventDefault();
                    del.stopPropagation();
                }
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            },
            up(up) {
                if(up && !this.detailsActive) {
                    this.selectPrevious();
                }
            },
            down(down) {
                if(down && !this.detailsActive) {
                    this.selectNext();
                }
            },
            row(row) {
                this.setRow(row);
            },
            query() {
                this.searchDebounced();
            }
        },
        mounted() {
            this.searchDebounced = debounce(() => {
                this.resetResults();
            }, searchDebounce);

            if(this.preselected) {
                this.selected = this.preselected;
            }

            this.focus();
        },
        methods: {
            focus() {
                if(this.row === 0) {
                    this.$refs.search.$refs.input.focus();
                }
                else {
                    this.focusCmdInput();
                }
            },
            save() {
                this.$emit("save", this.selected);
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
                        let len = results.length;
                        results = removeFromListByIds(results, this.selected);
                        this.results = this.results.concat(results);
                        this.resultCount += len;
                        this.count = count;
                        this.cancelToken = null;

                        if(this.resultCount >= this.count) {
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
            cancel() {
                if(this.row !== 0) {
                    this.$store.dispatch("selectRow", 0);
                }
                else {
                    this.$store.dispatch("popColumn");
                }
            },
            add(entity) {
                if(this.lock) {
                    return;
                }

                this.lock = true;

                this.$nextTick(() => {
                    if(findIndexById(this.selected, entity.id) < 0) {
                        this.selected.push(entity);
                        let index = findIndexById(this.results, entity.id);

                        if(index > -1) {
                            this.results.splice(index, 1);
                        }
                    }

                    this.lock = false;
                });
            },
            remove(id) {
                if(this.lock) {
                    return;
                }

                this.lock = true;

                this.$nextTick(() => {
                    let index = findIndexById(this.selected, id);

                    if(index > -1) {
                        let removed = this.selected.splice(index, 1);
                        index = findIndexById(this.results, id);

                        if(index < 0) {
                            this.results.push(removed[0]);
                        }
                    }

                    this.lock = false;
                });
            },
            resetResults() {
                if(this.cancelToken) {
                    this.cancelToken.cancel();
                }

                this.cancelToken = null;
                this.resetLoading++;
                this.results = [];
                this.resultCount = 0;
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
            setRow(row) {
                updateSelectedRow(row, this.results.length+this.selected.length+2, this.$store);
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

<i18n>
    {
        "en": {
            "selected": "Selection"
        },
        "de": {
            "selected": "Auswahl"
        }
    }
</i18n>
