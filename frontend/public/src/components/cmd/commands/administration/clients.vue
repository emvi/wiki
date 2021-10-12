<template>
    <emvi-cmd-selection element="emvi-cmd-selection-client"
                        :perform-search="search"
                        :placeholder="$t('placeholder')"
                        add-view="administration-clients-add"
                        :add-text="$t('add_text')"
                        :enter="enter"
                        :tab="tab"
                        :del="del"
                        :esc="esc"
                        :up="up"
                        :down="down"></emvi-cmd-selection>
</template>

<script>
    import {ClientService} from "../../../../service";
    import emviCmdSelection from "../../selection/selection.vue";

    export default {
        components: {emviCmdSelection},
        props: ["enter", "tab", "del", "esc", "up", "down"],
        data() {
            return {
                clients: null
            };
        },
        methods: {
            search(query) {
                query = query.trim().toLowerCase();

                return new Promise((resolve, reject) => {
                    this.getClients()
                        .then(clients => {
                            if(query === "") {
                                resolve({results: clients, count: 0});
                            }
                            else {
                                let filtered = [];

                                for(let i = 0; i < clients.length; i++) {
                                    if(clients[i].name.toLowerCase().includes(query)) {
                                        filtered.push(clients[i]);
                                    }
                                }

                                resolve({results: filtered, count: filtered.length});
                            }
                        })
                        .catch(e => {
                            reject(e);
                        });
                });
            },
            getClients() {
                return new Promise((resolve, reject) => {
                    if(this.clients !== null) {
                        resolve(this.clients);
                        return;
                    }

                    ClientService.getClients()
                        .then(clients => {
                            this.clients = clients;
                            resolve(this.clients);
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
            "title": "API clients",
            "placeholder": "Filter clients...",
            "add_text": "Add client"
        },
        "de": {
            "title": "API Clients",
            "placeholder": "Clients filtern...",
            "add_text": "Client hinzufügen"
        }
    }
</i18n>
