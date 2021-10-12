<template>
    <label class="form--element--field form--element--grey">
        <span class="label-active" v-if="value">{{value.name}}</span>
        <input v-bind:class="{'form--element--file': true, 'has-value': value}" type="file" v-on:change="handleFileChange"/>
        <label class="form--element--label">{{$t("label")}}</label>
    </label>
</template>

<script>
    export default {
        props: {
            value: File
        },
        methods: {
            handleFileChange(e) {
                this.$emit("input", e.target.files[0])
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "label": "Select File"
        },
        "de": {
            "label": "Datei auswählen"
        }
    }
</i18n>