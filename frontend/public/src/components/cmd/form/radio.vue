<template>
    <div :class="{'form': true, 'active': showActive}" ref="entry">
        <fieldset>
            <span v-for="(option, key) in options" :key="option.label" ref="element">
                <input type="radio"
                    :disabled="disabled"
                    :id="name+key"
                    :name="name"
                    :value="option.value"
                    :checked="value === option.value"
                    v-on:input="$emit('input', $event.target.value)"
                    v-on:focus="setSelected"
                    v-on:keydown.tab.stop.prevent="tab"
                    v-on:keydown.enter.stop.prevent="$emit('enter')"
                    v-on:keydown.esc.stop.prevent="$emit('esc')"
                    v-on:keydown.up.stop
                    v-on:keydown.down.stop />
                <label :for="name+key">
                    {{option.label}}
                    <span v-if="required">({{$t("required")}})</span>
                    <span v-if="optional">({{$t("optional")}})</span>
                </label>
            </span>
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
            name: {default: "", required: true}, // unique name or else it won't work
            options: {default: () => [], required: true}
        },
        methods: {
            focusElement() {
                let index = 0;

                for(let i = 0; i < this.options.length; i++) {
                    if(this.options[i].value === this.value) {
                        index = i;
                        break;
                    }
                }

                this.$refs.element[index].firstChild.focus();
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
