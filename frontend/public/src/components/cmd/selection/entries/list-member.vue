<template>
    <div>
        <emvi-cmd-selection-list-user v-if="entity.type === 'user'"
            :entity="entity"
            :index="index"
            :details="details"
            :enter="enter"
            :tab="tab"
            :del="del"
            :esc="esc"
            :up="up"
            :down="down"
            v-on:details="emitDetails"
            v-on:remove="remove"></emvi-cmd-selection-list-user>
        <emvi-cmd-selection-entity v-if="entity.type === 'group'"
            icon="group"
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
    </div>
</template>

<script>
    import {PropsMixin} from "../mixin";
    import {ArticlelistService} from "../../../../service";
    import emviCmdSelectionListUser from "./list-user.vue";
    import emviCmdSelectionEntity from "./entity.vue";

    export default {
        mixins: [PropsMixin],
        components: {emviCmdSelectionListUser, emviCmdSelectionEntity},
        methods: {
            remove() {
                this.resetError();
                let listId = this.$store.state.page.meta.get("id");

                ArticlelistService.removeMember(listId, this.entity.member_id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_removed"));
                        this.$emit("remove", this.entity.id);
                        this.toggleDetails(false);
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
            "confirmation": "Are you sure you would like to remove the member from this list?",
            "label_yes": "Yes, remove member",
            "toast_removed": "The member has been removed."
        },
        "de": {
            "confirmation": "Möchtest du dieses Mitglied wirklich aus der Liste entfernen?",
            "label_yes": "Ja, Mitglied entfernen",
            "toast_removed": "Das Mitglied wurde entfernt."
        }
    }
</i18n>
