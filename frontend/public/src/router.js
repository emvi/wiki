import VueRouter from "vue-router";
import axios from "axios";
import * as pages from "./pages";

const routes = [
    {name: "start", path: "/", component: pages.Start},
    {name: "intro", path: "/intro", component: pages.Introduction},
    {name: "articles", path: "/articles", component: pages.Articles},
    {name: "lists", path: "/lists", component: pages.Lists},
    {name: "members", path: "/members", component: pages.Members},
    {name: "groups", path: "/groups", component: pages.Groups},
    {name: "tags", path: "/tags", component: pages.Tags},
    {name: "activities", path: "/activities", component: pages.Activities},
    {name: "notifications", path: "/notifications", component: pages.Notifications},
    {name: "list", path: "/list/:slug", component: pages.List},
    {name: "member", path: "/member/:username", component: pages.Member},
    {name: "group", path: "/group/:slug", component: pages.Group},
    {name: "tag", path: "/tag/:name", component: pages.Tag},
    {name: "read", path: "/read/:slug", component: pages.Read},
    {name: "edit", path: "/edit/:slug", component: pages.Edit},
    {name: "edit_new", path: "/edit", component: pages.Edit},
    {name: "bookmarks_articles", path: "/bookmarks/articles", component: pages.BookmarksArticles},
    {name: "bookmarks_lists", path: "/bookmarks/lists", component: pages.BookmarksLists},
    {name: "drafts", path: "/drafts", component: pages.Drafts},
    {name: "private_articles", path: "/private/articles", component: pages.PrivateArticles},
    {name: "private_lists", path: "/private/lists", component: pages.PrivateLists},
    {name: "watch_articles", path: "/watch/articles", component: pages.WatchArticles},
    {name: "watch_groups", path: "/watch/groups", component: pages.WatchGroups},
    {name: "watch_lists", path: "/watch/lists", component: pages.WatchLists},
    {name: "billing", path: "/billing", component: pages.Billing},
    {name: "notfound", path: "*", component: pages.NotFound}
];

export function getRouter(store) {
    let router = new VueRouter({routes, mode: "history"});

    router.beforeEach((to, from, next) => {
        window.scrollTo(0, 0);
        store.dispatch("resetMeta");
        store.dispatch("select", 0);

        if(to.path !== "/"){
            axios.get(EMVI_WIKI_BACKEND_HOST+"/api/v1/auth")
            .then(() => {
                next();
            })
            .catch(() => {
                store.dispatch("redirectLogin");
            });
        }
        else{
            // root route authenticates user
            next();
        }
    });

    return router;
}
