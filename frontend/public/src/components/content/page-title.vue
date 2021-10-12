<template>
    <div class="headline no-select">
        <h1>
            <i :class="iconClass" v-if="icon"></i>
            <slot></slot>
        </h1>
        <small v-if="!disableResultCount && results === 1">{{results}} {{$t("result")}}</small>
        <small v-if="!disableResultCount && results > 1">{{results}} {{$t("results")}}</small>
    </div>
</template>

<script>
    export default {
        props: ["icon", "color", "results", "disableResultCount"],
        computed: {
            iconClass() {
                let c = {"icon size-40": true};
                c[`icon-${this.icon}`] = true;
                return c;
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "result": "result",
            "results": "results"
        },
        "de": {
            "result": "Ergebnis",
            "results": "Ergebnisse"
        }
    }
</i18n>
