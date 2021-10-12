<template>
    <emvi-preview-list :perform-search="search"
        icon="group"
        v-on:loaded="loaded"></emvi-preview-list>
</template>

<script>
    import {addAttrToListElements} from "../../util";
    import {UsergroupService} from "../../service";
    import emviPreviewList from "./preview-list.vue";

    export default {
        components: {emviPreviewList},
        props: ["id"],
        methods: {
            search() {
                return new Promise((resolve, reject) => {
                    UsergroupService.getMember(this.id)
                        .then(({results, count}) => {
                            let users = [];

                            for(let i = 0; i < results.length; i++) {
                                users.push(results[i].user);
                            }

                            resolve({results: addAttrToListElements(users, "type", "user"), count});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            loaded() {
                this.$emit("loaded");
            }
        }
    }
</script>
