<template>
    <div :class="{'form': true, 'active': showActive}" ref="entry">
        <fieldset :class="{'toggle': true, 'active': showActive}">
            <input ref="element"
                   type="checkbox"
                   :id="name"
                   :disabled="disabled"
                   :checked="value"
                   v-on:input="$emit('input', $event.target.checked)"
                   v-on:focus="setSelected"
                   v-on:keydown.tab.stop.prevent="tab"
                   v-on:keydown.enter.stop.prevent="$emit('enter')"
                   v-on:keydown.esc.stop.prevent="$emit('esc')"
                   v-on:keydown.up.stop="up"
                   v-on:keydown.down.stop="down" />
            <label :for="name">
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
            name: {default: "", required: true} // unique name or else it won't work
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
