<template>
    <emvi-cmd-selection-entity icon="list"
        :entity="entity"
        :index="index"
        :details="details"
        :remove-confirmation="$t('confirmation')"
        :remove-label-yes="$t('label_yes')"
        :enter="enter"
        :tab="tab"
        :del="del"
        :esc="esc"
        :up="up"
        :down="down"
        v-on:details="emitDetails"
        v-on:remove="remove"></emvi-cmd-selection-entity>
</template>

<script>
    import {PropsMixin} from "../mixin";
    import {BookmarkService} from "../../../../service";
    import emviCmdSelectionEntity from "./entity.vue";

    export default {
        mixins: [PropsMixin],
        components: {emviCmdSelectionEntity},
        methods: {
            remove() {
                this.resetError();
                BookmarkService.bookmark({article_list_id: this.entity.id})
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_removed"));
                        this.$emit("remove", this.entity.id);
                        this.toggleDetails(false);
                        this.$store.dispatch("setMeta", {key: "updateList", value: true});
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "confirmation": "Are you sure you would like to remove this list from your bookmarks?",
            "label_yes": "Yes, remove list",
            "toast_removed": "The list has been removed."
        },
        "de": {
            "confirmation": "Möchtest du diese Liste wirklich aus deinen Lesezeichen entfernen?",
            "label_yes": "Ja, Liste entfernen",
            "toast_removed": "Die Liste wurde entfernt."
        }
    }
</i18n>
