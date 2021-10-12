<template>
    <div>
        <mainmenu :active="active"></mainmenu>
        <div class="main-grid" style="justify-content: center;">
            <div class="main-grid--center hub-content">
                <slot></slot>
                <toast></toast>
            </div>
        </div>
        <div class="bottom-bar main-grid no-select">
            <span class="dot">&copy;</span>
            <a href="/" target="_blank" ref="noreferrer">Emvi Software GmbH</a>
            <span class="dot">&middot;</span>
            <a href="/terms" target="_blank" ref="noreferrer">{{$t("link_terms")}}</a>
            <span class="dot">&middot;</span>
            <a href="/privacy" target="_blank" ref="noreferrer">{{$t("link_privacy")}}</a>
            <span class="dot">&middot;</span>
            <a href="/legal" target="_blank" ref="noreferrer">{{$t("link_legal")}}</a>
        </div>
    </div>
</template>

<script>
    import mainmenu from "./menu.vue";
    import toast from "./toast.vue";

    export default {
        components: {
            mainmenu,
            toast
        },
        props: {
            active: {default: null}
        }
    }
</script>

<i18n>
    {
        "en": {
            "link_terms": "Terms",
            "link_privacy": "Privacy",
            "link_legal": "Legal"
        },
        "de": {
            "link_terms": "Nutzungsbedingungen",
            "link_privacy": "Datenschutzerklärung",
            "link_legal": "Impressum"
        }
    }
</i18n>
