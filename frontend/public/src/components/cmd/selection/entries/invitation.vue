<template>
    <emvi-cmd-selection-result :index="index" icon="mail" :details="details" v-on:mouseenter="hover = true" v-on:mouseleave="hover = false">
        <template>
            <div class="item">
                <p>{{entity.email}} {{entity.def_time | moment('LL')}}</p>
            </div>
            <emvi-cmd-shortcut :shortcut="$t('key_delete')" icon="trash" v-show="active" v-on:click="showCancel">
                {{$t("shortcut_cancel")}}
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
                v-on:click="action"
                v-on:select="setSelection"></emvi-cmd-selection-button>
        </template>
    </emvi-cmd-selection-result>
</template>

<script>
    import {SelectionMixin} from "../mixin";
    import {MemberService} from "../../../../service";
    import emviCmdSelectionResult from "../result.vue";
    import emviCmdShortcut from "../../content/shortcut.vue";
    import emviCmdSelectionButton from "../form/button.vue";

    export default {
        mixins: [SelectionMixin],
        components: {emviCmdSelectionResult, emviCmdShortcut, emviCmdSelectionButton},
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
                    }
                    else {
                        this.action();
                    }
                }
            },
            del(del) {
                if(del && this.active) {
                    this.showCancel();
                }
            },
            esc(esc) {
                if(esc && this.details) {
                    this.cancel();
                }
            }
        },
        methods: {
            showCancel() {
                this.toggleDetails(!this.details);
            },
            action() {
                this.resetError();
                MemberService.cancelInvitation(this.entity.id)
                    .then(() => {
                        this.$store.dispatch("success", this.$t("toast_saved"));
                        this.$emit("remove", this.entity.id);
                        this.cancel();
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
            "key_delete": "Del",
            "shortcut_cancel": "Cancel",
            "confirmation": "Are you sure you would like to cancel this invitation?",
            "label_no": "No",
            "label_yes": "Yes, cancel invitation",
            "toast_saved": "The invitation has been cancelled."
        },
        "de": {
            "key_delete": "Entf",
            "shortcut_cancel": "Abbrechen",
            "confirmation": "Möchtest du diese Einladung wirklich abbrechen?",
            "label_no": "Nein",
            "label_yes": "Ja, Einladung abbrechen",
            "toast_saved": "Die Einladung wurde abgebrochen."
        }
    }
</i18n>
