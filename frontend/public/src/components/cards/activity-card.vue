<template>
    <emvi-card :entity="entity.triggered_by_user" :active="active" :scroll-area="scrollArea" v-on:iconclick="openMember">
        <template>
            <span class="item activity">
                <p>
                    <strong ref="user" v-on:click="openMember">{{entity.triggered_by_user.firstname}} {{entity.triggered_by_user.lastname}}</strong>
                    <span v-html="entity.feed" ref="feed"></span>
                </p>
            </span>
            <small>{{entity.def_time | moment("from", "now")}}</small>
            <fieldset class="checkbox" ref="checkbox" v-if="notifications">
                <input type="checkbox" :id="entity.id" v-model="read" />
                <label :for="entity.id"></label>
            </fieldset>
        </template>
    </emvi-card>
</template>

<script>
    import {FeedService} from "../../service";
    import emviCard from "./card.vue";
    import emviShortcut from "../cmd/content/shortcut.vue";

    export default {
        components: {emviCard, emviShortcut},
        props: {
            entity: {},
            active: {default: false},
            tab: {default: false},
            enter: {default: false},
            scrollArea: ""
        },
        data() {
            return {
                // sets the focus inside the feed
                // 0 highlights the read/unread toggle
                // 1 highlights the triggering user
                // everything else highlights a link inside the feed text
                focus: 0,
                read: this.entity.read,
                notifications: 0
            };
        },
        computed: {
            picture() {
                return this.entity.triggered_by_user.picture;
            },
            links() {
                return this.$refs.feed.getElementsByTagName("A");
            }
        },
        watch: {
            active(active) {
                if(active) {
                    this.focus = 0;
                    this.setFocus();
                }
                else {
                    this.removeFocus();
                }
            },
            enter(enter) {
                if(enter && this.active) {
                    if(this.focus === 0 && this.notifications) {
                        this.toggleRead();
                    }
                    else {
                        this.followLink();
                    }
                }
            },
            tab(tab) {
                if(tab && this.active) {
                    this.focusNext();
                    this.removeFocus();
                    this.setFocus();
                }
            }
        },
        beforeMount() {
            this.entity.triggered_by_user.type = "user";
        },
        mounted() {
            this.notifications = this.$route.query.notifications ? 1 : 0;

            if(this.active) {
                this.setFocus();
            }
        },
        methods: {
            focusNext() {
                this.focus++;

                if(this.focus > this.links.length+this.notifications) {
                    this.focus = 0;
                }
            },
            setFocus() {
                if(this.focus === 0 && this.notifications) {
                    this.$refs.checkbox.classList.add("focus");
                }
                else if(this.focus === this.notifications) {
                    this.$refs.user.classList.add("focus");
                }
                else {
                    for(let i = 0; i < this.links.length; i++) {
                        if(this.focus-1-this.notifications === i) {
                            this.links[i].classList.add("focus");
                        }
                    }
                }
            },
            removeFocus() {
                if(this.notifications) {
                    this.$refs.checkbox.classList.remove("focus");
                }

                this.$refs.user.classList.remove("focus");

                for(let i = 0; i < this.links.length; i++) {
                    this.links[i].classList.remove("focus");
                }
            },
            toggleRead() {
                this.resetError();
                FeedService.toggleNotificationRead(this.entity.id)
                    .then(() => {
                        this.read = !this.read;
                        this.$store.dispatch("toggleNotificationRead", this.entity.id);
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            followLink() {
                if(this.focus === this.notifications) {
                    this.openMember();
                }
                else {
                    this.openLink();
                }
            },
            openMember() {
                this.$store.dispatch("resetCmd");
                this.$router.push(`/member/${this.entity.triggered_by_user.organization_member.username}`);
            },
            openLink() {
                this.$store.dispatch("resetCmd");
                this.$router.push(this.links[this.focus-1-this.notifications].getAttribute("href").substr(window.location.origin.length));
            }
        }
    }
</script>
