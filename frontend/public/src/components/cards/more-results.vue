<template>
    <div>
        {{results}} {{$t("text")}}
    </div>
</template>

<script>
    export default {
        props: ["results"]
    }
</script>

<i18n>
    {
        "en": {
            "text": "more results"
        },
        "de": {
            "text": "weitere Ergebnisse"
        }
    }
</i18n>
