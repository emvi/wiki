<template>
    <label v-bind:class="{'form--element--checkbox--label': true, 'form--has-error': error && !value, 'form--element--checkbox--focus': focus, 'form--element--checkbox--disabled': disabled}">
        <input v-bind:value="value"
               v-on:input="$emit('input', $event.target.checked)"
               v-bind:checked="value"
               v-bind:disabled="disabled"
               v-on:focus="focus = true"
               v-on:blur="focus = false"
               ref="checkBox"
               class="form--element--checkbox--field"
               type="checkbox" />
        <i class="form--element--checkbox--button"></i>
        <span>
            <template v-if="label">{{label}}</template>
            <slot></slot>
        </span>
    </label>
</template>

<script>
    export default {
        props: ["value", "label", "disabled", "error", "autofocus"],
        data() {
            return {
                focus: false
            }
        },
        mounted() {
            if(this.autofocus !== undefined) {
                this.$refs.checkBox.focus();
            }
        }
    }
</script>
