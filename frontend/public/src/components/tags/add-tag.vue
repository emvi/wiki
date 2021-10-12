<template>
    <input v-bind:class="{'tag edit': true, 'error': showError}"
           v-model="tag"
           v-on:keyup.enter="add"
           v-on:keydown.tab="tab"
           v-on:blur="add"
           :placeholder="$t('placeholder')"
           ref="input" />
</template>

<script>
    import {getTextWidth} from "../../util/tags";
    import {TagService} from "../../service";

    const tagMaxLength = 40;
    const placeholderFont = "24px Inter";

    export default {
        props: {
            articleId: {default: ""},
            tags: {default: []}
        },
        watch: {
            tag(value) {
                if(value.length > tagMaxLength) {
                    this.tag = value.substr(0, tagMaxLength);
                }
            }
        },
        data() {
            return {
                tag: "",
                showError: false,
                updateFunc: null
            };
        },
        mounted() {
            this.setInputWidth("");
            let self = this;
            this.updateFunc = function(e) {
                self.setInputWidth(e.target.value);
            };
            this.$refs.input.addEventListener("input", this.updateFunc);
        },
        beforeDestroy() {
            this.$refs.input.removeEventListener("input", this.updateFunc);
        },
        methods: {
            tab(e) {
                if(this.tag.length) {
                    e.preventDefault();
                    e.stopPropagation();
                    this.add();
                }
            },
            add() {
                this.showError = false;
                let tag = this.tag.trim();

                if(tag.length === 0) {
                    return;
                }

                if(this.tags.includes(tag)) {
                    this.showError = true;
                    this.resetError();
                    return;
                }

                if(!this.articleId) {
                    TagService.validateTag(tag)
                        .then(() => {
                            this.$emit("add", tag);
                            this.tag = "";
                            this.setInputWidth("");
                        })
                        .catch(() => {
                            this.showError = true;
                            this.resetError();
                        });
                }
                else {
                    TagService.addTag(this.articleId, tag)
                        .then(() => {
                            this.$emit("add", tag);
                            this.tag = "";
                            this.setInputWidth("");
                        })
                        .catch(() => {
                            this.showError = true;
                            this.resetError();
                        });
                }
            },
            resetError() {
                window.setTimeout(() => {
                    this.showError = false;
                }, 400);
            },
            setInputWidth(tag) {
                if(tag === "") {
                    tag = this.$t("placeholder");
                }

                if(tag.length < 5) {
                    tag = "xxxxx";
                }

                this.$refs.input.style.width = `${getTextWidth(tag, placeholderFont)}px`;
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "placeholder": "Add Tag"
        },
        "de": {
            "placeholder": "Tag hinzufügen"
        }
    }
</i18n>
