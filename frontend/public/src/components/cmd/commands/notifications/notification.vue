<template>
    <div :class="notificationClass" ref="notification">
        <emvi-cmd-avatar :entity="notification.triggered_by_user" v-on:click="openMember"></emvi-cmd-avatar>
        <span class="item activity">
            <p>
                <strong class="cursor-pointer" v-on:click="openMember" ref="user">{{notification.triggered_by_user.firstname}} {{notification.triggered_by_user.lastname}}</strong>
                <span v-html="notification.notification" ref="notification"></span>
            </p>
            <small class="info">
                {{notification.def_time | moment("from", "now")}}
            </small>
        </span>
        <fieldset class="checkbox" ref="checkbox">
            <input type="checkbox"
                   :id="notification.id"
                   v-model="isRead"
                   v-on:click.stop.prevent="toggleRead"
                   v-on:focus="click" />
            <label :for="notification.id"></label>
        </fieldset>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {FeedService} from "../../../../service";
    import emviCmdAvatar from "../../content/avatar.vue";
    import {scrollIntoViewArea} from "../../../../util";

    export default {
        components: {emviCmdAvatar},
        props: {
            notification: {},
            read: {default: false},
            index: {default: 0},
            tab: {default: false},
            enter: {default: false}
        },
        data() {
            return {
                // sets the focus inside the notification
                // 0 highlights the read/unread toggle
                // 1 highlights the triggering user
                // everything else highlights a link inside the notification text
                focus: 0,
                isRead: this.read
            };
        },
        computed: {
            ...mapGetters(["row"]),
            active() {
                return this.index === this.row;
            },
            notificationClass() {
                return {
                    "entry notification": true,
                    unread: !this.isRead,
                    active: this.active && !this.isTouch,
                    double: this.isMobile
                };
            },
            picture() {
                return this.notification.triggered_by_user.picture;
            },
            links() {
                return this.$refs.notification.getElementsByTagName("A");
            }
        },
        watch: {
            row() {
                if(this.active) {
                    scrollIntoViewArea(this.$refs.notification, document.getElementById("cmd-results"));
                }
            },
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
                    if(this.focus === 0) {
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
            },
            read(read) {
                this.isRead = read;
            }
        },
        beforeMount() {
            this.notification.triggered_by_user.type = "user";
        },
        mounted() {
            if(this.active) {
                this.setFocus();
            }
        },
        methods: {
            focusNext() {
                this.focus++;

                if(this.focus > this.links.length+1) {
                    this.focus = 0;
                }
            },
            setFocus() {
                if(this.focus === 0) {
                    this.$refs.checkbox.classList.add("focus");
                }
                else if(this.focus === 1) {
                    this.$refs.user.classList.add("focus");
                }
                else {
                    for(let i = 0; i < this.links.length; i++) {
                        if(this.focus-2 === i) {
                            this.links[i].classList.add("focus");
                        }
                    }
                }
            },
            removeFocus() {
                this.$refs.checkbox.classList.remove("focus");
                this.$refs.user.classList.remove("focus");

                for(let i = 0; i < this.links.length; i++) {
                    this.links[i].classList.remove("focus");
                }
            },
            click() {
                this.$store.dispatch("selectRow", this.index);
                this.focusCmdInput();
            },
            toggleRead() {
                this.resetError();
                FeedService.toggleNotificationRead(this.notification.id)
                    .then(() => {
                        this.$store.dispatch("toggleNotificationRead", this.notification.id);
                    })
                    .catch(e => {
                        this.setError(e);
                    });
            },
            followLink() {
                if(this.focus === 1) {
                    this.openMember();
                }
                else {
                    this.openLink();
                }
            },
            openMember() {
                this.$store.dispatch("resetCmd");
                this.$router.push(`/member/${this.notification.triggered_by_user.organization_member.username}`).catch(() => {});
            },
            openLink() {
                this.$store.dispatch("resetCmd");
                this.$router.push(this.links[this.focus-2].getAttribute("href").substr(window.location.origin.length)).catch(() => {});
            }
        }
    }
</script>
