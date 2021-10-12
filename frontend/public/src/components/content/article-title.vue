<template>
    <textarea ref="title"
              name="title"
              class="title"
              :placeholder="$t('placeholder_title')"
              rows="1"
              maxlength="100"
              v-model="text"
              v-on:keydown.enter.stop.prevent></textarea>
</template>

<script>
    export default {
        props: ["title"],
        data() {
            return {
                text: ""
            };
        },
        watch: {
            title(value) {
                this.text = value;
                this.setTitle(value);
            },
            text(value) {
                this.setTitle(value);
                this.$emit("change", value);
            }
        },
        mounted() {
            this.focus();
        },
        methods: {
            focus() {
                setTimeout(() => {
                    this.$refs.title.focus();

                    if(document.activeElement !== this.$refs.title) {
                        this.focus();
                    }
                }, 20);
            },
            setTitle(title) {
                let textarea = this.$refs.title;
                textarea.value = title;
                textarea.style.height = 0;
                textarea.style.height = textarea.scrollHeight+"px";
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "placeholder_title": "Title"
        },
        "de": {
            "placeholder_title": "Titel"
        }
    }
</i18n>
