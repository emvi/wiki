import {getQueryFromCmd, nameStartsWith} from "./util";

/*
{
    page: {
        en: ["name", ...],
        de: ["name", ...],
        ...
    },
    pageName: "name",
    path: "/path",
    requires_admin: true, // optional (hide entry for normal users and moderators)
    title: {
        en: "Title",
        de: "Titel",
        ...
    },
    icon: "icon"
}
*/
export const navigation = [
    {
        page: {
            en: ["start", "home", "dashboard"],
            de: ["start", "dashboard"]
        },
        pageName: "start",
        path: "/",
        title: {
            en: "Start",
            de: "Start"
        },
        icon: "home"
    },
    {
        page: {
            en: ["activities"],
            de: ["aktivitäten"]
        },
        pageName: "activities",
        path: "/activities",
        title: {
            en: "View all Activities",
            de: "Alle Aktivitäten ansehen"
        },
        icon: "activity"
    },
    {
        page: {
            en: ["articles"],
            de: ["artikel"]
        },
        pageName: "articles",
        path: "/articles",
        title: {
            en: "View all Articles",
            de: "Alle Artikel ansehen"
        },
        icon: "article"
    },
    {
        page: {
            en: ["lists"],
            de: ["listen"]
        },
        pageName: "lists",
        path: "/lists",
        title: {
            en: "View all Lists",
            de: "Alle Listen ansehen"
        },
        icon: "list"
    },
    {
        page: {
            en: ["members", "user"],
            de: ["mitglieder", "nutzer"]
        },
        pageName: "members",
        path: "/members",
        title: {
            en: "View all Members",
            de: "Alle Mitglieder ansehen"
        },
        icon: "user"
    },
    {
        page: {
            en: ["groups"],
            de: ["gruppen"]
        },
        pageName: "groups",
        path: "/groups",
        title: {
            en: "View all Groups",
            de: "Alle Gruppen ansehen"
        },
        icon: "group"
    },
    {
        page: {
            en: ["tags"],
            de: ["tags"]
        },
        pageName: "tags",
        path: "/tags",
        title: {
            en: "View all Tags",
            de: "Alle Tags ansehen"
        },
        icon: "tag"
    },
    {
        page: {
            en: ["bookmarks"],
            de: ["lesezeichen"]
        },
        pageName: "bookmarks_articles",
        path: "/bookmarks/articles",
        title: {
            en: "View your bookmarked Articles and Lists",
            de: "Deine Artikel und Listen mit Lesezeichen einsehen"
        },
        icon: "bookmark"
    },
    {
        page: {
            en: ["drafts"],
            de: ["entwürfe"]
        },
        pageName: "drafts",
        path: "/drafts",
        title: {
            en: "View your drafts",
            de: "Deine Entwürfe einsehen"
        },
        icon: "article"
    },
    {
        page: {
            en: ["watchlist"],
            de: ["beobachtungsliste"]
        },
        pageName: "watch_articles",
        path: "/watch/articles",
        title: {
            en: "View your watched Articles, Lists and Groups",
            de: "Deine beobachteten Artikel, Listen und Gruppen einsehen"
        },
        icon: "view"
    },
    {
        page: {
            en: ["private"],
            de: ["privat"]
        },
        pageName: "private_articles",
        path: "/private/articles",
        title: {
            en: "View your private Articles and Lists",
            de: "Deine privaten Artikel und Listen einsehen"
        },
        icon: "lock"
    },
    {
        page: {
            en: ["billing"],
            de: ["abrechnung"]
        },
        pageName: "billing",
        path: "/billing",
        requires_admin: true,
        title: {
            en: "Show billing details and settings",
            de: "Zeige Zahlungsinformationen und Einstellungen"
        },
        icon: "billing"
    }
];

// filterNavigation returns the navigation without the active page name and filters for query.
// The query starts with a dot (for navigation) followed by a page name or alias.
export function filterNavigation(activePageName, query, locale, isAdmin) {
    if(!query) {
        query = "";
    }

    query = getQueryFromCmd(query);
    let out = [];

    for(let i = 0; i < navigation.length; i++) {
        if(navigation[i].pageName !== activePageName && nameStartsWith(navigation[i].page[locale], query) &&
            (!navigation[i].requires_admin || isAdmin)) {
            out.push(navigation[i]);
        }
    }

    return out;
}

// findNavigationAlias returns the first alias for given query. The alias is read from left to right.
export function findNavigationAlias(nav, query, locale) {
    query = getQueryFromCmd(query);

    for(let i = 0; i < nav.page[locale].length; i++) {
        if(nav.page[locale][i].startsWith(query)) {
            return nav.page[locale][i];
        }
    }

    return "";
}
