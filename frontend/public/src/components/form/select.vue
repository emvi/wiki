<template>
    <fieldset>
        <select :name="name"
                :id="name"
                :tabindex="tabindex"
                :value="value"
                v-on:input="$emit('input', $event.target.value)">
            <option v-for="option in options"
                    :key="option.value"
                    :value="option.value">{{option.label}}</option>
        </select>
        <label :for="name">{{label}}</label>
        <span v-if="error">{{error}}</span>
    </fieldset>
</template>

<script>
    export default {
        props: {
            value: "",
            name: "",
            tabindex: 0,
            label: "",
            options: {default: []},
            error: ""
        }
    }
</script>
