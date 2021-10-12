<template>
    <div>
        <router-link :to="link">
            <span :class="{'level-1 article': true, 'active': entry.id === articleId}">
                <i>{{pos}}</i>
                <h4 :class="{'grey-50': entry.archived}">{{entry.latest_article_content.title}}</h4>
            </span>
        </router-link>
        <div class="indent">
            <slot v-if="entry.id === articleId"></slot>
        </div>
    </div>
</template>

<script>
    import {slugWithId} from "../../util";

    export default {
        props: ["entry", "listId", "articleId", "pos"],
        computed: {
            link() {
                return `/read/${slugWithId(this.entry.latest_article_content.title, this.entry.id)}?list=${this.listId}`;
            }
        }
    }
</script>
