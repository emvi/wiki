<template>
    <emvi-list :title="tag.name"
               title-icon="tag"
               icon="article"
               card="emvi-article-card"
               :perform-search="search"
               v-on:open="open">
        <template slot="headline">
            <small :title="tag.def_time | moment('LT')">{{tag.def_time | moment("ll")}}</small>
        </template>
    </emvi-list>
</template>

<script>
    import {SearchService, TagService} from "../service";
    import {slugWithId} from "../util";
    import {setPageTitle} from "./title";
    import {emviList} from "../components";

    export default {
        components: {emviList},
        data() {
            return {
                tag: {def_time: new Date()}
            };
        },
        methods: {
            open(result) {
                this.$router.push(`/read/${slugWithId(result.latest_article_content.title, result.id)}`);
            },
            search(query, filter, cancelToken) {
                if(!this.tag.id) {
                    return this.loadTag()
                        .then(() => {
                            return this.loadArticles(query, filter, cancelToken);
                        })
                        .catch(e => {
                            return Promise.reject(e)
                        });
                }
                else {
                    return this.loadArticles(query, filter, cancelToken);
                }
            },
            loadTag() {
                return new Promise((resolve, reject) => {
                    TagService.getTagByName(this.$route.params.name)
                        .then(tag => {
                            this.$store.dispatch("setMetaVars", [
                                {key: "id", value: tag.id},
                                {key: "name", value: tag.name}
                            ]);
                            this.tag = tag;
                            setPageTitle(this.tag.name);
                            resolve();
                        })
                        .catch(e => {
                            e.status = 404;
                            reject(e);
                        });
                });
            },
            loadArticles(query, filter, cancelToken) {
                filter.tag_ids = this.tag.id;

                return new Promise((resolve, reject) => {
                    SearchService.findArticles(query, filter, cancelToken)
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
