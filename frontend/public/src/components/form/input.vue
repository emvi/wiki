<template>
    <fieldset style="margin: 16px 0;">
        <input :type="type"
               :name="name"
               :id="name"
               :placeholder="placeholder"
               :tabindex="tabindex"
               :value="value"
               v-on:input="$emit('input', $event.target.value)" />
        <label :for="name">{{label}}</label>
        <small class="error" v-if="error">{{error}}</small>
    </fieldset>
</template>

<script>
    export default {
        props: {
            value: "",
            type: {default: "text"},
            name: "",
            placeholder: {default: " "},
            tabindex: 0,
            label: "",
            error: ""
        }
    }
</script>
