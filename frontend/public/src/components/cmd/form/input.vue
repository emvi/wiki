<template>
    <div :class="{'form': true,'active': showActive}" ref="entry">
        <fieldset>
            <input autocomplete="off"
                ref="element"
                placeholder=" "
                :type="type"
                :disabled="disabled"
                :value="value"
                v-on:input="$emit('input', $event.target.value)"
                v-on:change="$emit('change', $event)"
                v-on:focus="setSelected"
                v-on:keydown.tab.stop.prevent="tab"
                v-on:keydown.enter.stop.prevent="$emit('enter')"
                v-on:keydown.esc.stop.prevent="$emit('esc')"
                v-on:keydown.up.stop="up"
                v-on:keydown.down.stop="down" />
            <label>
                {{label}}
                <span v-if="required">({{$t("required")}})</span>
                <span v-if="optional">({{$t("optional")}})</span>
            </label>
            <small v-if="hint">{{hint}}</small>
            <small class="error" v-show="error">{{error}}</small>
        </fieldset>
    </div>
</template>

<script>
    import {FormMixin} from "./mixin";

    export default {
        mixins: [FormMixin],
        props: {
            type: {default: "text"}
        },
        methods: {
            up() {
                this.$emit("previous");
            },
            down() {
                this.$emit("next");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "required": "required",
            "optional": "optional"
        },
        "de": {
            "required": "Pflichtfeld",
            "optional": "optional"
        }
    }
</i18n>
