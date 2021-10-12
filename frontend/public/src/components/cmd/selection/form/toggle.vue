<template>
    <div :class="{'form': true, 'active': showActive}">
        <fieldset :class="{'toggle': true, 'active': showActive}">
            <input type="checkbox"
                   :id="name"
                   :disabled="disabled"
                   :checked="value"
                   v-on:input="$emit('input', $event.target.checked)"
                   v-on:focus="focus" />
            <label :for="name">{{label}}</label>
        </fieldset>
        <emvi-cmd-shortcut shortcut="Enter" v-show="active">
            {{$t("toggle")}}
        </emvi-cmd-shortcut>
    </div>
</template>

<script>
    import {FormMixin} from "../mixin";
    import emviCmdShortcut from "../../content/shortcut.vue";

    export default {
        mixins: [FormMixin],
        components: {emviCmdShortcut},
        props: {
            value: {default: false},
            disabled: {default: false},
            label: "",
            name: {default: "", required: true}, // unique name or else it won't work
            enter: {default: false}
        },
        watch: {
            enter(enter) {
                if(enter && this.active && !this.disabled) {
                    this.$emit("input", !this.value);
                }
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "toggle": "Toggle"
        },
        "de": {
            "toggle": "Ein-/Ausschalten"
        }
    }
</i18n>
