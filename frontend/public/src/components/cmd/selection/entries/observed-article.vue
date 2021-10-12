<template>
    <emvi-cmd-selection-entity icon="article"
        :entity="entity"
        :index="index"
        :disabled="entity.archived"
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
    import {ObserveService} from "../../../../service";
    import emviCmdSelectionEntity from "./entity.vue";

    export default {
        mixins: [PropsMixin],
        components: {emviCmdSelectionEntity},
        methods: {
            remove() {
                this.resetError();
                ObserveService.observe({article_id: this.entity.id})
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
            "confirmation": "Are you sure you would like to remove this article from your observed articles?",
            "label_yes": "Yes, remove article",
            "toast_removed": "The article has been removed."
        },
        "de": {
            "confirmation": "Möchtest du diesen Artikel wirklich aus deinen beobachteten Artikeln entfernen?",
            "label_yes": "Ja, Artikel entfernen",
            "toast_removed": "Der Artikel wurde entfernt."
        }
    }
</i18n>
