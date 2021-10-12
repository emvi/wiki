<template>
    <div class="hint red">
        <div class="hint-content">
            <p>{{$t("text_archived")}}: {{reason}}</p>
            <small>{{$t("hint")}}</small>
        </div>
        <i v-on:click.stop.prevent="cmdArchive" class="icon icon-back size-40 cursor-pointer"></i>
        <i v-on:click.stop.prevent="cmdDelete" class="icon icon-trash size-40 cursor-pointer"></i>
    </div>
</template>

<script>
    export default {
        props: ["reason"],
        methods: {
            cmdArchive() {
                this.$store.dispatch("closeCmd");
                this.$nextTick(() => {
                    this.$store.dispatch("resetCmd", this.$t("archive_cmd"));
                });
            },
            cmdDelete() {
                this.$store.dispatch("closeCmd");
                this.$nextTick(() => {
                    this.$store.dispatch("resetCmd", this.$t("delete_cmd"));
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "text_archived": "This article has been archived",
            "hint": "Use /archive to restore it or delete to remove it permanently.",
            "archive_cmd": "/archive",
            "delete_cmd": "/delete"
        },
        "de": {
            "text_archived": "Dieser Artikel wurde archiviert",
            "hint": "Verwende /archivieren um ihn wiederherzustellen oder /löschen um ihn entgültig zu entfernen.",
            "archive_cmd": "/archivieren",
            "delete_cmd": "/löschen"
        }
    }
</i18n>
