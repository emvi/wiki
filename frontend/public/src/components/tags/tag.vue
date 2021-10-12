<template>
    <span class="tag" v-on:click="open">
        {{tag}}
        <i class="icon icon-close" v-on:click.stop="remove"></i>
    </span>
</template>

<script>
    import {TagService} from "../../service";

    export default {
        props: ["articleId", "tag", "readOnly"],
        methods: {
            open() {
                if(!this.readOnly) {
                    this.$router.push(`/tag/${this.tag}`);
                }
            },
            remove() {
                if(this.articleId) {
                    TagService.removeTag(this.articleId, this.tag)
                        .then(() => {
                            this.$emit("remove", this.tag);
                        })
                        .catch(e => {
                            this.setError(e);
                        });
                }
                else {
                    this.$emit("remove", this.tag);
                }
            }
        }
    }
</script>
