<template>
    <emvi-label v-if="readingTime > 0">
        {{formatReadingTime(readingTime)}} {{$t("reading_time")}}
    </emvi-label>
</template>

<script>
    import emviLabel from "./label.vue";

    export default {
        components: {emviLabel},
        props: ["readingTime"],
        methods: {
            formatReadingTime(seconds) {
                if(!seconds) {
                    return seconds;
                }

                if(seconds < 60) {
                    return this.$t("less_than_one_minute");
                }

                let add = seconds%60 > 29 ? 1 : 0;
                return `${Math.abs(Math.floor(seconds/60+add))} ${this.$t("minutes")}`;
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "reading_time": "read",
            "less_than_one_minute": "< 1 min",
            "minutes": "min"
        },
        "de": {
            "reading_time": "Lesezeit",
            "less_than_one_minute": "< 1 min",
            "minutes": "Min."
        }
    }
</i18n>
