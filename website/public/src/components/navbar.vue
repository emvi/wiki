<template>
    <nav id="navbar" ref="navbar" v-bind:class="{'header no-select': true, 'header--dynamic': (scroll > 1 && !scrollable) || menu, 'header--scroll background--blue-10': scrollable}">
        <div class="container">
            <div class="row center-lg center-md center-sm center-xs">
                <div class="row col-lg-10 col-md-12 col-sm-12 col-xs-12">
                    <div class="col-lg-3 col-md-2 col-sm-6 col-xs-4 start-xs">
                        <a href="/" class="header--wordmark">
                            <img src="static/img/wordmark-black.svg" alt="Emvi" v-if="!theme" />
                            <img src="static/img/wordmark-white.svg" alt="Emvi" v-if="theme" />
                        </a>
                    </div>
                </div>
            </div>
        </div>
    </nav>
</template>

<script>
    import expand from "../components/expand.vue";
    import {scrollTo} from "../util/scroll.js";

    export default {
        components: {
            expand
        },
        props: {
            scrollable: {default: false},
            theme: {default: false}
        },
        computed: {
            isloggedin() {
                return this.$store.state.user.user;
            }
        },
        data() {
            return {
                menu: false,
                scroll: window.pageYOffset
            };
        },
        mounted() {
            window.addEventListener("scroll", this.scrollPosition);
            document.addEventListener("mouseup", this.documentClick);
        },
        beforeDestroy() {
            window.removeEventListener("scroll", this.scrollPosition);
            document.removeEventListener("mouseup", this.documentClick);
        },
        methods: {
            documentClick(e){
                let navbar = this.$refs.navbar;
                let navbarButton = this.$refs.navbarButton;
                let target = e.target;

                if(this.menu === true && navbar && navbar !== target && !navbar.contains(target) && navbarButton !== target) {
                    this.menu = false;
                }
            },
            login() {
                this.$store.dispatch("login");
            },
            scrollPosition() {
                this.scroll = window.pageYOffset;
            },
            features() {
                if (this.$route.path === "/") {
                    let features = document.getElementById("features");
                    if (this.menu) {
                        scrollTo(features, 64);
                        this.menu = false;
                    } else {
                        scrollTo(features, 0);
                    }
                } else {
                    this.$router.push("/#features");
                }
            },
            usecases() {
                if (this.$route.path === "/") {
                    let usecases = document.getElementById("usecases");
                    if (this.menu) {
                        scrollTo(usecases, 64);
                        this.menu = false;
                    } else {
                        scrollTo(usecases, 0);
                    }
                } else {
                    this.$router.push("/#usecases");
                }
            }
        }
    }
</script>
