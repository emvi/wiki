import moment from "moment/moment";
import VueI18n from "vue-i18n";
import {getLocale} from "./util";

// global translations
const messages = {
    "en": {
        "observe_article_on": "The article is now being observed.",
        "observe_article_off": "The article is not being observed anymore.",
        "observe_list_on": "The list is now being observed.",
        "observe_list_off": "The list is not being observed anymore.",
        "observe_group_on": "The group is now being observed.",
        "observe_group_off": "The group is not being observed anymore.",
        "bookmark_article_on": "The article was added to your bookmarks.",
        "bookmark_article_off": "The article was removed from your bookmarks.",
        "bookmark_list_on": "The list was added to your bookmarks.",
        "bookmark_list_off": "The list was removed from your bookmarks.",
        "pin_article_on": "The article was pinned to dashboard.",
        "pin_article_off": "The article was removed from dashboard.",
        "pin_list_on": "The list was pinned to dashboard.",
        "pin_list_off": "The list was removed from dashboard."
    },
    "de": {
        "observe_article_on": "Der Artikel wird jetzt beobachtet.",
        "observe_article_off": "Der Artikel wird nicht länger beobachtet.",
        "observe_list_on": "Die Liste wird jetzt beobachtet.",
        "observe_list_off": "Die Liste wird nicht länger beobachtet.",
        "observe_group_on": "Die Gruppe wird jetzt beobachtet.",
        "observe_group_off": "Die Gruppe wird nicht länger beobachtet.",
        "bookmark_article_on": "Der Artikel wurde in deiner Leseliste gespeichert.",
        "bookmark_article_off": "Der Artikel wurde aus deiner Leseliste entfernt.",
        "bookmark_list_on": "Die Liste wurde zu deiner Leseliste hinzugefügt.",
        "bookmark_list_off": "Die Liste wurde aus deiner Leseliste entfernt.",
        "pin_article_on": "Der Artikel wurde an das Dashboard gepinnt.",
        "pin_article_off": "Der Artikel wurde vom Dashboard entfernt.",
        "pin_list_on": "Die Liste wurde an das Dashboard gepinnt.",
        "pin_list_off": "Die Liste wurde vom Dashboard entfernt."
    }
};

export function getVueI18n() {
    return new VueI18n({
        locale: getLocale(),
        fallbackLocale: "en",
        silentTranslationWarn: true,
        messages
    });
}

export function setMomentAbbreviations() {
    moment.updateLocale("en", {
        relativeTime: {
            future: "in %s",
            past: "%s ago",
            s: "a few seconds",
            m: "1 min",
            mm: "%d min",
            h: "1 hour",
            hh: "%d hours",
            d: "1 day",
            dd: "%d days",
            M: "1 month",
            MM: "%d months",
            y: "1 year",
            yy: "%d years"
        }
    });

    moment.updateLocale("de", {
        relativeTime: {
            future: "in %s",
            past: "vor %s",
            s: "ein paar Sek.",
            m: "1 Min.",
            mm: "%d Min.",
            h: "1 Std.",
            hh: "%d Std.",
            d: "1 Tag",
            dd: "%d Tagen",
            M: "1 Mon.",
            MM: "%d Mon.",
            y: "1 Jahr",
            yy: "%d Jahren"
        }
    });
}
