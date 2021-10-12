<template>
    <div class="empty">
        <i :class="iconClass"></i>
        <p>{{$t("text")}}</p>
    </div>
</template>

<script>
    export default {
        props: ["icon"],
        computed: {
            iconClass() {
                let c = {icon: true};
                c[`icon-${this.icon}`] = true;
                return c;
            }
        }
    }
</script>


<i18n>
    {
        "en": {
            "text": "No results found"
        },
        "de": {
            "text": "Keine Ergebnisse gefunden"
        }
    }
</i18n>
