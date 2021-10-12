<template>
    <div :class="{'form': true, 'active': showActive}" ref="entry">
        <fieldset :class="{'active': showActive}">
            <textarea ref="element"
                    placeholder=" "
                    :disabled="disabled"
                    :value="value"
                    v-on:input="$emit('input', $event.target.value)"
                    v-on:focus="setSelected"
                    v-on:keydown.tab.stop.prevent="tab"
                    v-on:keydown.esc.stop.prevent="$emit('esc')"
                    v-on:keydown.enter.stop
                    v-on:keydown.up.stop
                    v-on:keydown.down.stop></textarea>
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
        mixins: [FormMixin]
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
