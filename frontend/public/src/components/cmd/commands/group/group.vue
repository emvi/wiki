<template>
    <div>
        <emvi-cmd-input :label="$t('label_name')"
                        :index="0"
                        v-model="name"
                        :error="validation['name']"
                        required="true"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="save"
                        v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-input :label="$t('label_info')"
                        :index="1"
                        v-model="info"
                        :error="validation['info']"
                        v-on:next="nextRow"
                        v-on:previous="previousRow"
                        v-on:enter="save"
                        v-on:esc="cancel"></emvi-cmd-input>
        <emvi-cmd-button icon="group"
                         color="purple"
                         :label="isNew ? $t('label_save_new') : $t('label_save_existing')"
                         :index="2"
                         v-on:next="nextRow"
                         v-on:previous="previousRow"
                         v-on:enter="save"
                         v-on:esc="cancel"></emvi-cmd-button>
    </div>
</template>

<script>
    import {updateSelectedRow} from "../../util";
    import {mapGetters} from "vuex";
    import {UsergroupService} from "../../../../service";
    import emviCmdInput from "../../form/input.vue";
    import emviCmdButton from "../../form/button.vue";
    import {slugWithId} from "../../../../util";

    export default {
        components: {emviCmdInput, emviCmdButton},
        props: ["esc"],
        data() {
            return {
                id: "",
                name: "",
                info: ""
            };
        },
        computed: {
            ...mapGetters(["row"]),
            isNew() {
                return !this.id;
            }
        },
        watch: {
            row(row) {
                updateSelectedRow(row, 3, this.$store);
            },
            esc(esc) {
                if(esc) {
                    this.cancel();
                }
            }
        },
        mounted() {
            if(this.$route.name === "group") {
                this.loadGroup();
            }
        },
        methods: {
            loadGroup() {
                this.id = this.$store.state.page.meta.get("id");

                UsergroupService.getUsergroup(this.id)
                    .then(group => {
                        this.name = group.group.name;
                        this.info = group.group.info;
                    });
            },
            save() {
                this.resetError();

                UsergroupService.saveUsergroup(this.id, this.name, this.info)
                    .then(id => {
                        this.$store.dispatch("success", this.$t(this.id ? "toast_saved_existing" : "toast_saved_new"));
                        this.$store.dispatch("resetCmd");

                        if(!this.id) {
                            this.$router.push(`/group/${slugWithId(this.name, id)}`);
                        }
                        else {
                            this.$store.dispatch("setMetaVars", [
                                {key: "updated", value: true},
                                {key: "name", value: this.name},
                                {key: "info", value: this.info}
                            ]);
                        }
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            cancel() {
                this.$store.dispatch("popColumn");
            },
            nextRow() {
                this.$store.dispatch("selectNextRow");
            },
            previousRow() {
                this.$store.dispatch("selectPreviousRow");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "label_name": "Name",
            "label_info": "Description",
            "label_save_new": "Create Group",
            "label_save_existing": "Update Group",
            "toast_saved_new": "The group has been created.",
            "toast_saved_existing": "The group has been updated."
        },
        "de": {
            "label_name": "Name",
            "label_info": "Beschreibung",
            "label_save_new": "Gruppe erstellen",
            "label_save_existing": "Gruppe speichern",
            "toast_saved_new": "Die Gruppe wurde erstellt.",
            "toast_saved_existing": "Die Gruppe wurde aktualisiert."
        }
    }
</i18n>
