<template>
    <div v-bind:class="{'form--element--select': true, 'form--has-error': error}">
        <select class="form--element--field has-value" v-bind:value="value" v-on:input="$emit('input', $event.target.value)">
            <slot></slot>
        </select>
        <label class="form--element--label">{{label}}</label>
        <span class="form--element--hint" v-if="error">{{error}}</span>
        <span class="form--element--hint" v-if="hint">{{hint}}</span>
    </div>
</template>

<script>
    export default {
        props: ["value", "label", "hint", "error"]
    }
</script>
