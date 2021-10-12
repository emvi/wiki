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
                v-on:focus.stop.prevent="focus"
                v-on:keydown.tab.stop.prevent="tab"
                v-on:keydown.enter.stop.prevent="$emit('enter')"
                v-on:keydown.esc.stop.prevent="$emit('esc')"
                v-on:keydown.up.stop="up"
                v-on:keydown.down.stop="down"
                v-on:click.stop.prevent />
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
    import {FormMixin} from "../mixin";

    export default {
        mixins: [FormMixin],
        props: {
            value: {default: ""},
            disabled: {default: false},
            hint: {default: ""},
            error: {default: ""},
            required: {default: false},
            optional: {default: false},
            label: "",
            icon: "",
            type: {default: "text"}
        },
        watch: {
            active(active) {
                if(active) {
                    this.focus();
                }
            }
        },
        mounted() {
            if(this.active) {
                this.focus();
            }
        },
        methods: {
            focus() {
                this.$store.dispatch("selectRow", this.index);
                this.$emit("select", this.selectionIndex);
                this.$refs.element.focus();
            },
            tab(e) {
                if(e.shiftKey) {
                    this.up();
                }
                else {
                    this.down();
                }
            },
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
