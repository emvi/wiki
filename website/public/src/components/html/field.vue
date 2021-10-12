<template>
    <div v-bind:class="formClasses">
        <input v-bind:class="{'form--element--field': true, 'has-value': value}"
            placeholder=" "
            :type="type"
            v-bind:value="value"
            v-on:input="$emit('input', $event.target.value)"
            v-bind:disabled="disabled"
            ref="inputField" />
        <label class="form--element--label">{{label}}</label>
        <small class="form--element--hint" v-if="error">{{error}}</small>
        <small class="form--element--hint" v-if="hint">{{hint}}</small>
    </div>
</template>

<script>
    export default {
        props: {
            value: null,
            type: null,
            label: null,
            disabled: null,
            error: null,
            hint: null,
            autofocus: null,
            color: {default: "blue"}
        },
        watch: {
            error(value) {
                this.formClasses["form--has-error"] = value;
            }
        },
        data() {
            return {
                formClasses: {},
            };
        },
        created() {
            this.formClasses["form--element form--element--"+this.color] = true;
        },
        mounted() {
            if(this.autofocus !== undefined) {
                this.$refs.inputField.focus();
            }
        }
    }
</script>
