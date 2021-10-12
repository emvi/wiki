<template>
    <div>
        <emvi-cmd-loading v-show="loading"></emvi-cmd-loading>
        <emvi-cmd-search-result v-for="(result, index) in results"
                         :key="result.key"
                         :entity="result"
                         :icon="result.type !== 'user' ? result.type : ''"
                         :index="index"
                         :preview="previewActive"
                         :tab="tab"
                         :up="up"
                         :down="down"
                         v-on:preview="togglePreview"
                         v-on:click="navigate(index)">
            <emvi-cmd-entity :entity="result" :query="cmd"></emvi-cmd-entity>
        </emvi-cmd-search-result>
        <emvi-cmd-empty v-show="!loading && !results.length"></emvi-cmd-empty>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {CancelToken} from "axios";
    import {viewSearch} from "./type";
    import {SearchService} from "../../service";
    import {getQueryFromCmd, updateSelectedRow} from "./util";
    import {debounce, slugWithId} from "../../util";
    import emviCmdSearchResult from "./search/search-result.vue";
    import emviCmdEntity from "./content/entity.vue";
    import emviCmdEmpty from "./content/empty.vue";
    import emviCmdLoading from "./content/loading.vue";
    import emviCmdPreview from "./search/preview.vue";

    const searchDebounce = 250;

    export default {
        components: {emviCmdSearchResult, emviCmdEntity, emviCmdEmpty, emviCmdLoading, emviCmdPreview},
        props: ["up", "down", "enter", "tab", "esc"],
        data() {
            return {
                searchDebounced: null,
                searchCancelToken: null,
                loading: true,
                results: [],
                previewActive: false
            };
        },
        computed: {
            ...mapGetters(["cmd", "row"])
        },
        watch: {
            cmd() {
                this.searchDebounced();
            },
            row(row) {
                updateSelectedRow(row, this.results.length, this.$store);
            },
            up(up) {
                if(up && !this.previewActive) {
                    this.$store.dispatch("selectPreviousRow");
                }
            },
            down(down) {
                if(down && !this.previewActive) {
                    this.$store.dispatch("selectNextRow");
                }
            },
            enter(enter) {
                if(enter) {
                    this.navigate();
                }
            },
            tab(tab) {
                if(tab) {
                    this.togglePreview();
                }
            },
            esc(esc) {
                if(esc) {
                    this.$store.dispatch("popColumn");
                }
            }
        },
        mounted() {
            this.searchDebounced = debounce(() => {this.search();}, searchDebounce);
            this.search();
        },
        methods: {
            search() {
                if(this.searchCancelToken) {
                    this.searchCancelToken.cancel();
                }

                this.resetResults();
                this.searchCancelToken = CancelToken.source();

                SearchService.findAll(getQueryFromCmd(this.cmd), null, 0, this.searchCancelToken)
                    .then(results => {
                        this.searchCancelToken = null;
                        this.loading = false;
                        this.$store.dispatch("selectRow", 0);
                        this.setSearchResults(results);
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            setSearchResults(results) {
                let out = [];
                out = out.concat(this.appendSearchResults(results.articles, "article"));
                out = out.concat(this.appendSearchResults(results.lists, "list"));
                out = out.concat(this.appendSearchResults(results.groups, "group"));
                out = out.concat(this.appendSearchResults(results.user, "user"));
                out = out.concat(this.appendSearchResults(results.tags, "tag"));
                this.results = out;
            },
            appendSearchResults(results, type) {
                let out = [];

                if(results) {
                    for(let i = 0; i < results.length; i++) {
                        results[i].key = `${type}_${results[i].id}`;
                        results[i].type = type;
                        out.push(results[i]);
                    }
                }

                return out;
            },
            resetResults() {
                this.results = [];
                this.loading = true;
            },
            togglePreview(index) {
                if(!this.results.length) {
                    return;
                }

                if(index !== undefined) {
                    this.previewActive = !this.previewActive || index !== this.row;
                    this.$store.dispatch("selectRow", index);
                    this.focusCmdInput();
                }
                else {
                    this.previewActive = !this.previewActive;
                }
            },
            navigate(index) {
                if(!this.results.length) {
                    return;
                }

                if(index === undefined) {
                    index = this.row;
                }

                let entity = this.results[index];
                let path = "";

                switch(entity.type) {
                    case "article":
                        path = `/read/${slugWithId(entity.latest_article_content.title, entity.id)}`;
                        break;
                    case "list":
                        path = `/list/${slugWithId(entity.name.name, entity.id)}`;
                        break;
                    case "group":
                        path = `/group/${slugWithId(entity.name, entity.id)}`;
                        break;
                    case "user":
                        path = `/member/${entity.organization_member.username}`;
                        break;
                    default:
                        path = `/tag/${entity.name}`;
                }

                if(path !== this.$route.path) {
                    this.$router.push(path);
                }

                this.$store.dispatch("closeCmd");
            }
        }
    }
</script>
