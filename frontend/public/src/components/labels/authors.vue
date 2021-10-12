<template>
    <small class="label" :title="title">
        {{text}}
        <span v-if="count > 0">+ {{count}}</span>
    </small>
</template>

<script>
    export default {
        props: ["authors"],
        data() {
            return {
                text: "",
                count: 0,
                title: ""
            };
        },
        watch: {
            authors() {
                let names = [];

                for (let i = 0; i < this.authors.length; i++) {
                    names.push(`${this.authors[i].firstname} ${this.authors[i].lastname}`);
                }

                this.text = names[0];
                this.count = names.length > 1 ? names.length - 1 : 0;
                this.title = names.join(", ");
            }
        }
    }
</script>
