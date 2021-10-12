<template>
    <emvi-list :title="$t('title')"
               title-icon="user"
               icon="user"
               card="emvi-member-card"
               :perform-search="search"></emvi-list>
</template>

<script>
    import {SearchService} from "../service";
    import {TitleMixin} from "./title";
    import {emviList} from "../components";
    import {addAttrToListElements} from "../util";

    export default {
        mixins: [TitleMixin],
        components: {emviList},
        methods: {
            search(query, filter, cancelToken) {
                filter.sort_lastname = "asc";
                filter.sort_firstname = "asc";

                return new Promise((resolve, reject) => {
                    SearchService.findUser(query, filter, cancelToken)
                        .then(data => {
                            resolve({results: addAttrToListElements(data.user, "type", "user"), count: data.count});
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
            "title": "Members"
        },
        "de": {
            "title": "Mitglieder"
        }
    }
</i18n>
