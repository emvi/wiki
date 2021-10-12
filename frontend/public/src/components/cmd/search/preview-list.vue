<template>
    <div>
        <emvi-cmd-loading v-show="!loaded"></emvi-cmd-loading>
        <emvi-cmd-preview-empty :icon="icon" v-show="loaded && !results.length"></emvi-cmd-preview-empty>
        <div v-show="loaded">
            <emvi-cmd-preview-list-entry v-for="result in results"
                                         :key="result.id"
                                         :icon="resultIcon"
                                         :type="result.type"
                                         :entity="result">
                <emvi-cmd-entity :entity="result"></emvi-cmd-entity>
            </emvi-cmd-preview-list-entry>
            <emvi-cmd-more-results v-if="count > 20" :results="count-20"></emvi-cmd-more-results>
        </div>
    </div>
</template>

<script>
    import emviCmdPreviewEmpty from "./preview-empty.vue";
    import emviCmdLoading from "../content/loading.vue";
    import emviCmdPreviewListEntry from "./preview-list-entry.vue";
    import emviCmdEntity from "../content/entity.vue";
    import emviCmdMoreResults from "./more-results.vue";

    export default {
        components: {
            emviCmdPreviewEmpty,
            emviCmdLoading,
            emviCmdPreviewListEntry,
            emviCmdEntity,
            emviCmdMoreResults
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
            }
        }
    }
</script>
