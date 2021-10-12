<template>
    <div v-bind:class="btnClasses" v-on:click="click">
        <div class="toggle" v-if="toggle && !disabled && icon !== 'toggle'">
            <input type="checkbox"
                v-bind:id="toggle"
                v-bind:checked="checked"
                v-on:click.stop="$emit('checked')" />
            <label v-bind:for="toggle"></label>
        </div>
        <div v-bind:class="iconClasses" v-if="icon && !img"><span v-if="label">{{label}}</span></div>
        <img :src="img" alt="" v-if="!number && img" />
        <div class="button--number" v-if="number">{{number}}</div>
        <div class="button--label">
            <slot></slot>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            checked: {default: false},
            icon: null,
            img: null,
            type: {default: "action"},
            color: {default: "blue"},
            active: {default: false},
            disabled: {default: false}, // greyed out and no pointer
            toggle: null,
            spinner: {default: false},
            upsidedown: {default: false},
            number: null,
            label: null,
            cursor: null
        },
        watch: {
            active(value) {
                this.btnClasses["button--active"] = value && !this.disabled;
                this.btnClasses["button--passive"] = !value && !this.disabled;
            }
        },
        computed: {
            btnClasses() {
                let btnClasses = {
                    "button--active": this.active && !this.disabled,
                    "button--passive": !this.active && !this.disabled,
                    "button--disabled": this.disabled,
                    "cursor-default": this.cursor
                };
                btnClasses["button--" + this.type] = true;
                btnClasses["button--" + this.color] = true;
                return btnClasses;
            },
            iconClasses() {
                let iconClasses = {
                    "button--icon icon": true,
                    "icon-is-spinning": this.spinner,
                    "icon-rotate-180": this.upsidedown,
                    "icon--label": this.label
                };
                iconClasses["icon-" + this.icon] = true;
                return iconClasses;
            }
        },
        methods: {
            click(e) {
                if(this.disabled) {
                    return;
                }

                this.$emit("click", e);
            }
        }
    }
</script>
