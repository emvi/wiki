<template>
    <emvi-label :title="publishedFullDate">
        <template v-if="published">
            {{published | moment("ll")}}
        </template>
        <template v-if="!published">
            {{$t("unpublished")}}
        </template>
    </emvi-label>
</template>

<script>
    import emviLabel from "./label.vue";

    export default {
        components: {emviLabel},
        props: {
            published: {default: () => new Date()}
        },
        computed: {
            publishedFullDate() {
                if(this.published) {
                    return this.$moment(this.published).format("LLL");
                }

                return "";
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "unpublished": "unpublished"
        },
        "de": {
            "unpublished": "unveröffentlicht"
        }
    }
</i18n>
