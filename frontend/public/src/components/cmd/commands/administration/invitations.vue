<template>
    <emvi-cmd-selection element="emvi-cmd-selection-invitation"
                        :perform-search="search"
                        :placeholder="$t('placeholder')"
                        add-view="administration-members-add"
                        :add-text="$t('add_text')"
                        :enter="enter"
                        :tab="tab"
                        :del="del"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import {MemberService} from "../../../../service";
    import emviCmdSelection from "../../selection/selection.vue";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "tab", "del", "esc", "up", "down"],
        data() {
            return {
                results: null
            };
        },
        methods: {
            search(query) {
                return new Promise((resolve, reject) => {
                    if(this.results === null) {
                        MemberService.getInvitations()
                            .then(invitations => {
                                this.results = invitations;
                                resolve({results: this.filterResults(query), count: 0});
                            })
                            .catch(e => {
                                reject(e);
                            });
                    }
                    else {
                        resolve({results: this.filterResults(query), count: 0});
                    }
                });
            },
            filterResults(query) {
                query = query.trim();

                if(!query) {
                    return this.results;
                }

                let results = [];

                for(let i = 0; i < this.results.length; i++) {
                    if(this.results[i].email.includes(query)) {
                        results.push(this.results[i]);
                    }
                }

                return results;
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "placeholder": "Filter invitations...",
            "add_text": "Invite Members"
        },
        "de": {
            "placeholder": "Einladungen filtern...",
            "add_text": "Mitglieder einladen"
        }
    }
</i18n>
