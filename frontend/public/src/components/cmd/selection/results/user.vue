<template>
    <emvi-cmd-selection-entity :index="index"
        :disable-add="disableAdd"
        :disable-remove="disableRemove"
        :entity="entity"
        :enter="enter"
        :del="del"
        :up="up"
        :down="down"
        v-on:add="add"
        v-on:remove="remove"></emvi-cmd-selection-entity>
</template>

<script>
    import {PropsMixin} from "../mixin";
    import emviCmdSelectionEntity from "./entity.vue";

    export default {
        mixins: [PropsMixin],
        components: {
            emviCmdSelectionEntity
        },
        props: {
            disableAdd: {default: false},
            disableRemove: {default: false}
        },
        methods: {
            add(entity) {
                this.$emit("add", entity);
            },
            remove(id) {
                this.$emit("remove", id);
            }
        }
    }
</script>
