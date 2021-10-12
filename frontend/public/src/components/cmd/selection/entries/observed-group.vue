<template>
    <emvi-cmd-selection-entity icon="group"
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
    import {ObserveService} from "../../../../service";
    import emviCmdSelectionEntity from "./entity.vue";

    export default {
        mixins: [PropsMixin],
        components: {emviCmdSelectionEntity},
        methods: {
            remove() {
                this.resetError();
                ObserveService.observe({user_group_id: this.entity.id})
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
            "confirmation": "Are you sure you would like to remove this group from your observed groups?",
            "label_yes": "Yes, remove list",
            "toast_removed": "The list has been removed."
        },
        "de": {
            "confirmation": "Möchtest du diese Gruppe wirklich aus deinen beobachteten Gruppen entfernen?",
            "label_yes": "Ja, Liste entfernen",
            "toast_removed": "Die Liste wurde entfernt."
        }
    }
</i18n>
