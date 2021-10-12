<template>
    <div class="entry sticky">
        <div class="avatar size-32">
            <i class="icon icon-filter"></i>
        </div>
        <input ref="input"
               id="cmd-search-input"
               type="text"
               autocomplete="off"
               name="query"
               :placeholder="placeholder"
               :value="value"
               v-on:input="$emit('input', $event.target.value)"
               v-on:focus="focus"
               v-on:keydown.enter.prevent.stop
               v-on:keydown.tab.prevent.stop="tab"
               v-on:keydown.esc.prevent.stop="$emit('reset')"
               v-on:keydown.up.prevent.stop="$emit('previous')"
               v-on:keydown.down.prevent.stop="$emit('next')" />
    </div>
</template>

<script>
    export default {
        props: ["value", "placeholder", "index"],
        methods: {
            tab(e) {
                if(e.shiftKey) {
                    this.$emit("previous");
                }
                else {
                    this.$emit("next");
                }
            },
            focus() {
                this.$store.dispatch("selectRow", this.index);
                this.$emit("focus");
            }
        }
    }
</script>
