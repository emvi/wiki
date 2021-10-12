<template>
    <emvi-layout>
        <emvi-layout-narrow>
            <emvi-page-title disable-result-count="true">{{$t("title")}}</emvi-page-title>
        </emvi-layout-narrow>
    </emvi-layout>
</template>

<script>
    import {TitleMixin} from "./title";
    import {emviLayout, emviLayoutNarrow, emviPageTitle} from "../components";

    export default {
        mixins: [TitleMixin],
        components: {emviLayout, emviLayoutNarrow, emviPageTitle}
    }
</script>

<i18n>
    {
        "en": {
            "title": "The page does not exist."
        },
        "de": {
            "title": "Die Seite existiert nicht."
        }
    }
</i18n>
