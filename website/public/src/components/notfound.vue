<template>
    <div>
        <div class="title--article">
            <h3>{{$t("title")}}</h3>
        </div>
        <div class="spacer-16"></div>
        <p>
            <slot></slot>
        </p>
    </div>
</template>

<script>
    export default {}
</script>

<i18n>
    {
        "en": {
            "title": "Not found"
        },
        "de": {
            "title": "Nicht gefunden"
        }
    }
</i18n>
