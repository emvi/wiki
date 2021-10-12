import VueRouter from "vue-router";
import Vue from "vue";
import axios from "axios";
import * as pages from "./pages";

Vue.use(VueRouter);

const routes = [
    {path: "/registration", component: pages.Registration},
    {path: "/newsletter/confirm", component: pages.NewsletterConfirm},
    {path: "/newsletter/unsubscribe", component: pages.NewsletterUnsubscribe},
    {path: "/join/:code", component: pages.Join},
    {path: "/join", component: pages.Join},
    {path: "/new", component: pages.NewOrganization, meta: {protected: true}},
    {path: "/organizations", component: pages.Organizations, meta: {protected: true}},
    {path: "/account", component: pages.Account, meta: {protected: true}},
    {path: "*", component: pages.Error404}
];

export function getRouter(store) {
    let router = new VueRouter({
        routes, mode: "history",
        scrollBehavior: function (to, from, savedPosition) {
            if (to.hash) {
                return {selector: to.hash};
            } else {
                return {x: 0, y: 0};
            }
        }
    });

    // router interceptor to check token for protected pages
    router.beforeEach((to, from, next) => {
        window.scrollTo(0, 0);

        if (to.meta.protected) {
            axios.get(EMVI_WIKI_AUTH_HOST + "/api/v1/auth/token")
            .then(() => {
                next();
            })
            .catch((e) => {
                store.dispatch("logout");
                next("/");
            });
        } else {
            next();
        }
    });

    // router interceptor to set page title if available and reset toast
    router.beforeEach((to, from, next) => {
        store.dispatch("resetToast");

        // build title
        let title = "Emvi";

        if (EMVI_WIKI_INTEGRATION) {
            title += " — INTEGRATION";
        }

        document.title = title;
        next();
    });

    return router;
}
