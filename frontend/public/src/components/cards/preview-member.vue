<template>
    <div>
        <emvi-cmd-loading v-show="!loaded"></emvi-cmd-loading>
        <div class="padding-16" v-show="loaded">
            <emvi-profile-header :entity="user">{{title}}</emvi-profile-header>
            <small>{{$t("info")}}</small>
            <p v-if="user.organization_member.info">{{user.organization_member.info}}</p>
            <p v-if="!user.organization_member.info">-</p>
            <small>{{$t("username")}}</small>
            <p>{{user.organization_member.username}}</p>
            <small>{{$t("email")}}</small>
            <p>{{user.email}}</p>
        </div>
    </div>
</template>

<script>
    import {UserService} from "../../service";
    import emviCmdLoading from "../content/loading.vue";
    import emviProfileHeader from "../content/profile-header.vue";

    export default {
        components: {emviCmdLoading, emviProfileHeader},
        props: ["id"],
        data() {
            return {
                user: {organization_member: {}},
                loaded: false
            };
        },
        computed: {
            title() {
                return `${this.user.firstname} ${this.user.lastname}`;
            },
            picture() {
                return this.user.picture;
            }
        },
        mounted() {
            this.loadMember();
        },
        methods: {
            loadMember() {
                if(!this.loaded) {
                    this.resetError();
                    UserService.getUser(this.id)
                        .then(user => {
                            this.user = user;
                            this.loaded = true;
                            this.$emit("loaded");
                        })
                        .catch(e => {
                            this.setError(e);
                        });
                }
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "info": "Info",
            "username": "Username",
            "email": "Email"
        },
        "de": {
            "info": "Info",
            "username": "Nutzername",
            "email": "E-Mail"
        }
    }
</i18n>
