<template>
    <div class="tabs no-select">
        <router-link v-for="tab in tabs" :key="tab.path" :to="tab.path">
            <h2>{{tab.label}}</h2>
        </router-link>
        <span class="shortcut" v-if="!isTouch">
            <span>
                <span class="key">{{$t("shift")}}</span>
                <span>&#43;</span>
                <span class="key">&lt;-</span>
                <span class="key">-&gt;</span>
            </span>
        </span>
    </div>
</template>

<script>
    export default {
        props: ["tabs"],
        data() {
            return {
                keydownHandler: null,
                active: 0
            };
        },
        mounted() {
            this.bindKeys();
            this.setActive();
        },
        beforeDestroy() {
            this.unbindKeys();
        },
        methods: {
            bindKeys() {
                this.keydownHandler = e => {
                    if(!e.shiftKey) {
                        return;
                    }

                    let prevent = true;

                    switch(e.code) {
                        case "ArrowLeft":
                            this.previousTab();
                            break;
                        case "ArrowRight":
                            this.nextTab();
                            break;
                        default:
                            prevent = false;
                    }

                    if(prevent) {
                        e.preventDefault();
                        e.stopPropagation();
                    }
                };
                window.addEventListener("keydown", this.keydownHandler);
            },
            unbindKeys() {
                window.removeEventListener("keydown", this.keydownHandler);
            },
            setActive() {
                let path = this.$route.path.toLowerCase();

                for(let i = 0; i < this.tabs.length; i++) {
                    if(path === this.tabs[i].path) {
                        this.active = i;
                        break;
                    }
                }
            },
            nextTab() {
                this.active++;

                if(this.active > this.tabs.length-1) {
                    this.active = 0;
                }

                this.openTab();
            },
            previousTab() {
                this.active--;

                if(this.active < 0) {
                    this.active = this.tabs.length-1;
                }

                this.openTab();
            },
            openTab() {
                this.$router.push(this.tabs[this.active].path);
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "shift": "Shift"
        },
        "de": {
            "shift": "Shift"
        }
    }
</i18n>
