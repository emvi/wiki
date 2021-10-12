<template>
    <div>
        <emvi-loading v-show="!loaded"></emvi-loading>
        <emvi-preview-empty :icon="icon" v-show="loaded && !results.length"></emvi-preview-empty>
        <div v-show="loaded">
            <emvi-preview-list-entry v-for="result in results"
                                     :key="result.id"
                                     :icon="resultIcon"
                                     :entity="result"
                                     v-on:click="openResult(result)">
                <emvi-entity :entity="result"></emvi-entity>
            </emvi-preview-list-entry>
            <emvi-more-results v-if="count > 20" :results="count-20"></emvi-more-results>
        </div>
    </div>
</template>

<script>
    import {slugWithId} from "../../util";
    import emviLoading from "../content/loading.vue";
    import emviPreviewEmpty from "./preview-empty.vue";
    import emviMoreResults from "./more-results.vue";
    import emviPreviewListEntry from "./preview-list-entry.vue";
    import emviEntity from "../content/entity.vue";

    export default {
        components: {
            emviLoading,
            emviPreviewEmpty,
            emviMoreResults,
            emviPreviewListEntry,
            emviEntity
        },
        props: ["performSearch", "icon", "resultIcon"],
        data() {
            return {
                results: [],
                count: 0,
                loaded: false
            };
        },
        mounted() {
            this.loadResults();
        },
        methods: {
            loadResults() {
                if(!this.loaded) {
                    this.resetError();
                    this.performSearch()
                        .then(({results, count}) => {
                            this.results = results;
                            this.count = count;
                            this.loaded = true;
                            this.$emit("loaded");
                        })
                        .catch(e => {
                            this.setError(e);
                        });
                }
            },
            openResult(entity) {
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
            }
        }
    }
</script>
