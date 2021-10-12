<template>
    <emvi-list :title="group.name"
               title-icon="group"
               icon="user"
               card="emvi-member-card"
               :perform-search="search"
               disable-search="true"
               v-on:open="open">
        <template slot="headline">
            <span class="labels">
                <small>{{group.info}}</small>
                <emvi-published-label :published="group.def_time"></emvi-published-label>
                <emvi-modified-label :mod-time="group.mod_time"></emvi-modified-label>
            </span>
            <emvi-watched-label type="group" :observed="observed"></emvi-watched-label>
        </template>
    </emvi-list>
</template>

<script>
    import {mapGetters} from "vuex";
    import {UsergroupService} from "../service";
    import {getIdFromSlug, slugWithId} from "../util";
    import {setPageTitle} from "./title";
    import {emviList, emviPublishedLabel, emviModifiedLabel, emviWatchedLabel} from "../components";

    export default {
        components: {emviList, emviPublishedLabel, emviModifiedLabel, emviWatchedLabel},
        data() {
            return {
                group: {
                    def_time: new Date(),
                    mod_time: new Date()
                },
                observed: false
            };
        },
        computed: {
            ...mapGetters(["metaUpdate"])
        },
        watch: {
            metaUpdate(update) {
                if(update > 1 && this.$store.state.page.meta.get("updated")) {
                    this.group.name = this.$store.state.page.meta.get("name");
                    this.group.info = this.$store.state.page.meta.get("info");
                    this.$store.dispatch("setMeta", {key: "updated", value: undefined});
                }

                let observed = this.$store.state.page.meta.get("observed");

                if(observed !== undefined) {
                    this.observed = observed;
                }
            }
        },
        methods: {
            open(result) {
                this.$router.push(`/read/${slugWithId(result.latest_article_content.title, result.id)}`);
            },
            search(query, filter, cancelToken) {
                if(!this.group.id) {
                    return this.loadGroup()
                        .then(() => {
                            return this.loadMember(query, filter, cancelToken);
                        })
                        .catch(e => {
                            return Promise.reject(e)
                        });
                }
                else {
                    return this.loadMember(query, filter, cancelToken);
                }
            },
            loadGroup() {
                let id = getIdFromSlug(this.$route.params.slug);

                return new Promise((resolve, reject) => {
                    UsergroupService.getUsergroup(id)
                        .then(group => {
                            this.$store.dispatch("setMetaVars", [
                                {key: "id", value: group.group.id},
                                {key: "immutable", value: group.group.immutable},
                                {key: "observed", value: group.observed},
                                {key: "moderator",  value: group.moderator}
                            ]);
                            this.observed = group.observed;
                            this.group = group.group;
                            setPageTitle(this.group.name);
                            resolve();
                        })
                        .catch(e => {
                            e.status = 404;
                            reject(e);
                        });
                });
            },
            loadMember(query, filter, cancelToken) {
                return new Promise((resolve, reject) => {
                    UsergroupService.getMember(this.group.id, filter, cancelToken)
                        .then(({results, count}) => {
                            let members = [];

                            for(let i = 0; i < results.length; i++) {
                                results[i].user.type = "user";
                                members.push(results[i].user);
                            }

                            resolve({results: members, count});
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "edited_before": "edited",
            "edited_after": " "
        },
        "de": {
            "edited_before": " ",
            "edited_after": "bearbeitet"
        }
    }
</i18n>
