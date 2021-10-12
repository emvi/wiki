<template>
    <emvi-cmd-selection-result :index="index" :disabled="disabled" icon="article" :details="details" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <emvi-cmd-entity :entity="entity"></emvi-cmd-entity>
            <emvi-cmd-shortcut :shortcut="$t('key_delete')" icon="trash" v-show="active" v-on:click="showDelete">
                {{$t("shortcut_remove")}}
            </emvi-cmd-shortcut>
            <emvi-cmd-shortcut :shortcut="$t('shift')" v-show="active">
                <span>&#43;</span>
                <span class="key">&uarr;</span>
                <span class="key">&darr;</span>
                {{$t("shortcut_sort")}}
            </emvi-cmd-shortcut>
        </template>
        <template slot="details">
            <p>{{$t("confirmation")}}</p>
            <emvi-cmd-selection-button icon="back"
                :label="$t('label_no')"
                :index="index"
                :selection="selection"
                :selection-index="0"
                v-on:click="cancel"
                v-on:select="setSelection"></emvi-cmd-selection-button>
            <emvi-cmd-selection-button icon="trash"
                color="red"
                :label="$t('label_yes')"
                :index="index"
                :selection="selection"
                :selection-index="1"
                v-on:click="remove"
                v-on:select="setSelection"></emvi-cmd-selection-button>
        </template>
    </emvi-cmd-selection-result>
</template>

<script>
    import {SelectionMixin} from "../mixin";
    import {ArticlelistService} from "../../../../service";
    import emviCmdSelectionResult from "../result.vue";
    import emviCmdEntity from "../../content/entity.vue";
    import emviCmdShortcut from "../../content/shortcut.vue";
    import emviCmdSelectionButton from "../form/button.vue";

    export default {
        mixins: [SelectionMixin],
        components: {
            emviCmdSelectionResult,
            emviCmdEntity,
            emviCmdShortcut,
            emviCmdSelectionButton
        },
        data() {
            return {
                maxSelectionIndex: 1
            };
        },
        watch: {
            enter(enter) {
                if(enter && this.active && this.details) {
                    if(this.selection === 0) {
                        this.cancel();
                    } else {
                        this.remove();
                    }
                }
            },
            del(del) {
                if(del && this.active) {
                    this.showDelete();
                }
            },
            esc(esc) {
                if(esc && this.details) {
                    this.cancel();
                }
            },
            up(up) {
                if(up && up.shiftKey && this.active && !this.details) {
                    this.sort(-1);
                }
            },
            down(down) {
                if(down && down.shiftKey && this.active && !this.details) {
                    this.sort(1);
                }
            }
        },
        methods: {
            showDelete() {
                this.toggleDetails(!this.details);
            },
            sort(dir) {
                this.resetError();
                let listId = this.$store.state.page.meta.get("id");

                ArticlelistService.moveEntry(listId, this.entity.id, dir)
                    .then(() => {
                        this.$emit("swap", {id: this.entity.id, direction: dir});
                        this.$store.dispatch("setMeta", {key: "updateList", value: true});
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            remove() {
                this.resetError();
                let listId = this.$store.state.page.meta.get("id");

                ArticlelistService.removeEntry(listId, this.entity.id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_removed"));
                        this.$emit("remove", this.entity.id);
                        this.cancel();
                        this.$store.dispatch("setMeta", {key: "updateList", value: true});
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            cancel() {
                this.toggleDetails(false);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "confirmation": "Are you sure you would like to remove the article from this list?",
            "label_yes": "Yes, remove article",
            "label_no": "No",
            "toast_removed": "The article has been removed.",
            "shortcut_sort": "Move",
            "shortcut_remove": "Remove",
            "shift": "Shift",
            "key_delete": "Del"
        },
        "de": {
            "confirmation": "Möchtest du diesen Artikel wirklich aus dieser Liste entfernen?",
            "label_yes": "Ja, Artikel entfernen",
            "label_no": "Nein",
            "toast_removed": "Der Artikel wurde entfernt.",
            "shortcut_sort": "Verschieben",
            "shortcut_remove": "Entfernen",
            "shift": "Shift",
            "key_delete": "Entf"
        }
    }
</i18n>
