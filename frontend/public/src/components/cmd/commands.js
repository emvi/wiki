import {getQueryFromCmd, nameStartsWith} from "./util";
import {toggleBookmark, toggleObserve, togglePinned} from "./command_funcs";

/*
{
    command: {
        en: ["name", "alias"],
        de: ["name", "alias"],
        ...
    },
    hidden: true, // optional, hides the command from search
    pageNames: ["name", ...], // optional
    excludePageNames: ["/path", ...], // optional
    buttonPageNames: ["name", ...], // optional, show button on page
    requires_admin: true, // optional (disable command for normal users and moderators)
    requires_mod: true, // optional (disable command for normal users)
    expert: true, // optional (show expert hint if entry organization)
    description: {
        en: "Description",
        de: "Beschreibung",
        ...
    },
    button_label: { // description is used if not set
        en: "Button Label",
        de: "Button Beschriftung"
        ...
    },
    icon: "icon",
    run(vue) {
        // ...
    },
    disabled(vue) { // programmatically disable/enable entry (optional)
        return true/false
    }
}
*/
export const commands = [
    {
        command: {
            en: ["edit"],
            de: ["bearbeiten"]
        },
        pageNames: ["read"],
        description: {
            en: "Edit this article",
            de: "Bearbeite diesen Artikel"
        },
        icon: "edit",
        run(vue) {
            let langId = vue.$store.state.page.meta.get("lang_id");
            vue.$router.push(`/edit/${vue.$route.params.slug}?lang=${langId}`);
            vue.$store.dispatch("resetCmd");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                !!vue.$store.state.page.meta.get("archived") ||
                !vue.$store.state.page.meta.get("write") ||
                (vue.$store.state.page.meta.get("translationMissing") && !vue.isExpert);
        }
    },
    {
        command: {
            en: ["publish", "save"],
            de: ["veröffentlichen", "speichern"]
        },
        pageNames: ["edit", "edit_new"],
        description: {
            en: "Publish this article",
            de: "Veröffentliche diesen Artikel"
        },
        icon: "publish",
        run(vue) {
            vue.$store.dispatch("pushColumn", "publish");
        }
    },
    {
        command: {
            en: ["discard", "leave", "exit"],
            de: ["verwerfen", "verlassen"]
        },
        pageNames: ["edit", "edit_new"],
        description: {
            en: "Discard changes and leave this article",
            de: "Verwerfe die Änderungen und verlasse diesen Artikel"
        },
        icon: "trash",
        run(vue) {
            vue.$store.dispatch("pushColumn", "discard");
        }
    },
    {
        command: {
            en: ["recommend"],
            de: ["empfehlen"]
        },
        pageNames: ["read"],
        description: {
            en: "Recommend this article",
            de: "Empfehle diesen Artikel"
        },
        icon: "send",
        run(vue) {
            vue.$store.dispatch("pushColumn", "recommend");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                vue.$store.state.page.meta.get("private");
        }
    },
    {
        command: {
            en: ["invite"],
            de: ["einladen"]
        },
        pageNames: ["edit", "edit_new"],
        description: {
            en: "Invite members to edit this article",
            de: "Lade Mitglieder ein diesen Artikel zu bearbeiten"
        },
        icon: "add-user",
        run(vue) {
            vue.$store.dispatch("pushColumn", "invite");
        }
    },
    {
        command: {
            en: ["history", "versions"],
            de: ["historie", "versionen"]
        },
        pageNames: ["read"],
        description: {
            en: "View and manage the versions of this article",
            de: "Siehe und bearbeite die Versionen dieses Artikels"
        },
        icon: "history",
        run(vue) {
            vue.$store.dispatch("pushColumn", "history");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound");
        }
    },
    {
        command: {
            en: ["language", "translation"],
            de: ["sprache", "übersetzung"]
        },
        pageNames: ["read"],
        description: {
            en: "Switch the language you're reading this article in",
            de: "Wechsel die Sprache in der du den Artikel ließt"
        },
        icon: "globe",
        run(vue) {
            vue.$store.dispatch("pushColumn", "translation");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound");
        }
    },
    {
        command: {
            en: ["rtl", "right-to-left"],
            de: ["rnl", "rechts-nach-links"]
        },
        pageNames: ["edit", "edit_new"],
        description: {
            en: "Toggle the writing direction",
            de: "Schalte die Schreibrichtung um"
        },
        button_label: {
            en: "Change Text Direction",
            de: "Schreibrichtung ändern"
        },
        icon: "edit",
        run(vue) {
            let rtl = vue.$store.state.page.meta.get("rtl");
            vue.$store.dispatch("setMeta", {key: "rtl", value: !rtl});
            vue.$store.dispatch("resetCmd");
        }
    },
    {
        command: {
            en: ["permissions"],
            de: ["berechtigungen"]
        },
        pageNames: ["edit", "edit_new"],
        description: {
            en: "Manage permissions",
            de: "Berechtigungen bearbeiten"
        },
        icon: "key",
        run(vue) {
            vue.$store.dispatch("pushColumn", "article-permissions");
        }
    },
    {
        command: {
            en: ["edit"],
            de: ["bearbeiten"]
        },
        pageNames: ["bookmarks_articles"],
        description: {
            en: "Manage bookmarked articles",
            de: "Verwalte deine Lesezeichen"
        },
        icon: "edit",
        run(vue) {
            vue.$store.dispatch("pushColumn", "bookmarks-articles");
        }
    },
    {
        command: {
            en: ["edit"],
            de: ["bearbeiten"]
        },
        pageNames: ["bookmarks_lists"],
        description: {
            en: "Manage bookmarked lists",
            de: "Verwalte deine Lesezeichen"
        },
        icon: "edit",
        run(vue) {
            vue.$store.dispatch("pushColumn", "bookmarks-lists");
        }
    },
    {
        command: {
            en: ["edit"],
            de: ["bearbeiten"]
        },
        pageNames: ["watch_articles"],
        description: {
            en: "Manage observed articles",
            de: "Verwalte deine beobachteten Artikel"
        },
        icon: "edit",
        run(vue) {
            vue.$store.dispatch("pushColumn", "observed-articles");
        }
    },
    {
        command: {
            en: ["edit"],
            de: ["bearbeiten"]
        },
        pageNames: ["watch_lists"],
        description: {
            en: "Manage observed lists",
            de: "Verwalte deine beobachteten Listen"
        },
        icon: "edit",
        run(vue) {
            vue.$store.dispatch("pushColumn", "observed-lists");
        }
    },
    {
        command: {
            en: ["edit"],
            de: ["bearbeiten"]
        },
        pageNames: ["watch_groups"],
        description: {
            en: "Manage observed groups",
            de: "Verwalte deine beobachteten Gruppen"
        },
        icon: "edit",
        run(vue) {
            vue.$store.dispatch("pushColumn", "observed-groups");
        }
    },
    {
        command: {
            en: ["details", "edit", "rename", "tag"],
            de: ["details", "bearbeiten", "umbenennen", "tag"]
        },
        pageNames: ["tag"],
        description: {
            en: "Edit this tag",
            de: "Bearbeite diesen Tag"
        },
        icon: "edit",
        run(vue) {
            vue.$store.dispatch("pushColumn", "tag");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                !vue.isAdmin ||
                !vue.isMod ||
                vue.isReadOnly;
        }
    },
    {
        command: {
            en: ["details"],
            de: ["details"]
        },
        pageNames: ["list"],
        description: {
            en: "Edit this list",
            de: "Bearbeite diese Liste"
        },
        icon: "edit",
        run(vue) {
            vue.$store.dispatch("pushColumn", "list");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                !vue.$store.state.page.meta.get("moderator");
        }
    },
    {
        command: {
            en: ["details"],
            de: ["details"]
        },
        pageNames: ["group"],
        expert: true,
        description: {
            en: "Edit this group",
            de: "Bearbeite diese Gruppe"
        },
        icon: "edit",
        run(vue) {
            vue.$store.dispatch("pushColumn", "group");
        },
        disabled(vue) {
            return vue.$store.state.page.meta.get("immutable") ||
                !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                !vue.isExpert ||
                !vue.$store.state.page.meta.get("moderator");
        }
    },
    {
        command: {
            en: ["edit", "entries", "entry"],
            de: ["bearbeiten", "einträge", "eintrag"]
        },
        pageNames: ["list"],
        description: {
            en: "Manage list entries",
            de: "Bearbeite die Einträge in dieser Liste"
        },
        icon: "list",
        run(vue) {
            vue.$store.dispatch("pushColumn", "list-entries");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                !vue.$store.state.page.meta.get("moderator");
        }
    },
    {
        command: {
            en: ["permissions"],
            de: ["berechtigungen"]
        },
        pageNames: ["list"],
        description: {
            en: "Manage list permissions",
            de: "Berechtigungen bearbeiten"
        },
        icon: "key",
        run(vue) {
            vue.$store.dispatch("pushColumn", "list-permissions");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                !vue.$store.state.page.meta.get("moderator");
        }
    },
    {
        command: {
            en: ["edit", "entries", "entry"],
            de: ["bearbeiten", "einträge", "eintrag"]
        },
        pageNames: ["group"],
        expert: true,
        description: {
            en: "Manage group members",
            de: "Bearbeite die Mitglieder dieser Gruppe"
        },
        icon: "user",
        run(vue) {
            vue.$store.dispatch("pushColumn", "group-member");
        },
        disabled(vue) {
            return vue.$store.state.page.meta.get("immutable") ||
                !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                !vue.isExpert ||
                !vue.$store.state.page.meta.get("moderator");
        }
    },
    {
        command: {
            en: ["watch", "observe"],
            de: ["beobachten"]
        },
        pageNames: ["read"],
        description: {
            en: "Watch this article and receive notifications",
            de: "Beobachte diesen Artikel und erhalte Benachrichtungen"
        },
        icon: "notification",
        run(vue) {
            toggleObserve(vue, "article");
            vue.$store.dispatch("resetCmd");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound");
        }
    },
    {
        command: {
            en: ["watch", "observe"],
            de: ["beobachten"]
        },
        pageNames: ["list"],
        description: {
            en: "Watch this list and receive notifications",
            de: "Beobachte diese Liste und erhalte Benachrichtungen"
        },
        icon: "notification",
        run(vue) {
            toggleObserve(vue, "list");
            vue.$store.dispatch("resetCmd");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound");
        }
    },
    {
        command: {
            en: ["watch", "observe"],
            de: ["beobachten"]
        },
        pageNames: ["group"],
        description: {
            en: "Watch this group and receive notifications",
            de: "Beobachte diese Gruppe und erhalte Benachrichtungen"
        },
        icon: "notification",
        run(vue) {
            toggleObserve(vue, "group");
            vue.$store.dispatch("resetCmd");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound");
        }
    },
    {
        command: {
            en: ["bookmark"],
            de: ["merken", "lesezeichen"]
        },
        pageNames: ["read"],
        description: {
            en: "Bookmark this article",
            de: "Merke dir diesen Artikel"
        },
        icon: "bookmark",
        run(vue) {
            toggleBookmark(vue, "article");
            vue.$store.dispatch("resetCmd");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound");
        }
    },
    {
        command: {
            en: ["bookmark"],
            de: ["merken", "lesezeichen"]
        },
        pageNames: ["list"],
        description: {
            en: "Bookmark this list",
            de: "Merke dir diese Liste"
        },
        icon: "bookmark",
        run(vue) {
            toggleBookmark(vue, "list");
            vue.$store.dispatch("resetCmd");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound");
        }
    },
    {
        command: {
            en: ["pin"],
            de: ["anheften"]
        },
        pageNames: ["read"],
        requires_mod: true,
        description: {
            en: "Pin this article to the dashboard for everyone",
            de: "Hefte diesen Artikel an das Dashboard für alle"
        },
        icon: "pin",
        run(vue) {
            togglePinned(vue, "article");
            vue.$store.dispatch("resetCmd");
        },
        disabled(vue) {
            return !vue.isMod ||
                !!vue.$store.state.page.meta.get("notFound") ||
                vue.$store.state.page.meta.get("private");
        }
    },
    {
        command: {
            en: ["pin"],
            de: ["anpinnen"]
        },
        pageNames: ["list"],
        requires_mod: true,
        description: {
            en: "Pin this list to the dashboard for everyone",
            de: "Hefte diese Liste an das Dashboard für alle"
        },
        icon: "pin",
        run(vue) {
            togglePinned(vue, "list");
            vue.$store.dispatch("resetCmd");
        },
        disabled(vue) {
            return !vue.isMod ||
                !!vue.$store.state.page.meta.get("notFound") ||
                vue.$store.state.page.meta.get("private");
        }
    },
    {
        command: {
            en: ["duplicate", "copy"],
            de: ["duplizieren", "kopieren"]
        },
        pageNames: ["read"],
        description: {
            en: "Duplicate this article",
            de: "Dupliziere diesen Artikel"
        },
        icon: "copy",
        run(vue) {
            vue.$store.dispatch("pushColumn", "duplicate-article");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                vue.$store.state.page.meta.get("translationMissing");
        }
    },
    {
        command: {
            en: ["export"],
            de: ["exportieren"]
        },
        pageNames: ["read"],
        description: {
            en: "Export this article",
            de: "Exportiere diesen Artikel"
        },
        icon: "download",
        run(vue) {
            vue.$store.dispatch("pushColumn", "export-article");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound");
        }
    },
    {
        command: {
            en: ["archive", "restore"],
            de: ["archivieren", "wiederherstellen"]
        },
        pageNames: ["read"],
        description: {
            en: "Archive or restore this article",
            de: "Archiviere diesen Artikel oder stelle ihn wieder her"
        },
        icon: "archive",
        run(vue) {
            vue.$store.dispatch("pushColumn", "archive-article");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                !vue.$store.state.page.meta.get("write");
        }
    },
    {
        command: {
            en: ["delete"],
            de: ["löschen"]
        },
        pageNames: ["read"],
        description: {
            en: "Delete this article",
            de: "Lösche diesen Artikel"
        },
        icon: "trash",
        run(vue) {
            vue.$store.dispatch("pushColumn", "delete-article");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                !vue.$store.state.page.meta.get("write");
        }
    },
    {
        command: {
            en: ["delete"],
            de: ["löschen"]
        },
        pageNames: ["list"],
        description: {
            en: "Delete this list",
            de: "Lösche diese Liste"
        },
        icon: "trash",
        run(vue) {
            vue.$store.dispatch("pushColumn", "delete-list");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                !vue.$store.state.page.meta.get("moderator");
        }
    },
    {
        command: {
            en: ["delete"],
            de: ["löschen"]
        },
        pageNames: ["group"],
        expert: true,
        description: {
            en: "Delete this group",
            de: "Lösche diese Gruppe"
        },
        icon: "trash",
        run(vue) {
            vue.$store.dispatch("pushColumn", "delete-group");
        },
        disabled(vue) {
            return vue.$store.state.page.meta.get("immutable") ||
                !!vue.$store.state.page.meta.get("notFound") ||
                vue.isReadOnly ||
                !vue.isExpert ||
                !vue.$store.state.page.meta.get("moderator");
        }
    },
    {
        command: {
            en: ["delete"],
            de: ["löschen"]
        },
        pageNames: ["tag"],
        description: {
            en: "Delete this tag",
            de: "Lösche diesen Tag"
        },
        icon: "trash",
        run(vue) {
            vue.$store.dispatch("pushColumn", "delete-tag");
        },
        disabled(vue) {
            return !!vue.$store.state.page.meta.get("notFound") ||
                !vue.isAdmin ||
                !vue.isMod ||
                vue.isReadOnly;
        }
    },
    {
        command: {
            en: ["article", "write"],
            de: ["artikel", "schreiben"]
        },
        excludePageNames: ["edit"],
        buttonPageNames: ["articles"],
        description: {
            en: "Create a new article",
            de: "Lege einen neuen Artikel an"
        },
        button_label: {
            en: "New Article",
            de: "Artikel erstellen"
        },
        icon: "article",
        run(vue) {
            vue.$router.push("/edit");
            vue.$store.dispatch("resetCmd");
        },
        disabled(vue) {
            return vue.isReadOnly;
        }
    },
    {
        command: {
            en: ["list"],
            de: ["liste"]
        },
        excludePageNames: ["list"],
        buttonPageNames: ["lists"],
        description: {
            en: "Create a new list",
            de: "Lege einen neue Liste an"
        },
        button_label: {
            en: "New List",
            de: "Liste anlegen"
        },
        icon: "list",
        run(vue) {
            vue.$store.dispatch("pushColumn", "list");
        },
        disabled(vue) {
            return vue.isReadOnly;
        }
    },
    {
        command: {
            en: ["group"],
            de: ["gruppe"]
        },
        excludePageNames: ["group"],
        buttonPageNames: ["groups"],
        expert: true,
        description: {
            en: "Create a new group",
            de: "Lege einen neue Gruppe an"
        },
        button_label: {
            en: "New Group",
            de: "Gruppe anlegen"
        },
        icon: "group",
        run(vue) {
            vue.$store.dispatch("pushColumn", "group");
        },
        disabled(vue) {
            return !vue.isExpert || vue.isReadOnly;
        }
    },
    {
        command: {
            en: ["notifications"],
            de: ["benachrichtigungen"]
        },
        description: {
            en: "Read your notifications",
            de: "Lese deine Benachrichtigungen"
        },
        icon: "notification",
        run(vue) {
            vue.$store.dispatch("pushColumn", "notifications");
        }
    },
    {
        command: {
            en: ["new"],
            de: ["neu"]
        },
        description: {
            en: "Create a new article, list, or group",
            de: "Lege einen neuen Artikel, Liste oder Gruppe an"
        },
        icon: "add",
        run(vue) {
            vue.$store.dispatch("pushColumn", "create");
        },
        disabled(vue) {
            return vue.isReadOnly;
        }
    },
    {
        command: {
            en: ["back", "last"],
            de: ["zurück", "letzte"]
        },
        description: {
            en: "Go back one page",
            de: "Gehe eine Seite zurück"
        },
        icon: "compass",
        run(vue) {
            vue.$router.back();
            vue.$store.dispatch("resetCmd");
        }
    },
    {
        command: {
            en: ["forward", "next"],
            de: ["vorwärts", "nächste"]
        },
        description: {
            en: "Go forward one page",
            de: "Gehe eine Seite weiter"
        },
        icon: "compass",
        run(vue) {
            vue.$router.forward();
            vue.$store.dispatch("resetCmd");
        }
    },
    {
        command: {
            en: ["darkmode", "lightmode"],
            de: ["darkmode", "lightmode", "nachtmodus", "tagmodus"]
        },
        description: {
            en: "Toggle between light and dark mode",
            de: "Wechsel zwischen Tag- und Nachmodus"
        },
        icon: "night",
        run(vue) {
            vue.$store.dispatch("toggleDarkmode");
            vue.$store.dispatch("resetCmd");
        }
    },
    {
        command: {
            en: ["invite"],
            de: ["einladen"]
        },
        requires_admin: true,
        excludePageNames: ["edit"],
        buttonPageNames: ["members"],
        description: {
            en: "Invite new members to this organization",
            de: "Lade neue Mitglieder in diese Organisation ein"
        },
        button_label: {
            en: "Invite Member",
            de: "Mitglieder einladen"
        },
        icon: "add-user",
        run(vue) {
            vue.$store.dispatch("pushColumn", "administration-members-add");
        }
    },
    {
        command: {
            en: ["administration"],
            de: ["administration"]
        },
        requires_admin: true,
        description: {
            en: "Manage organization settings, members, languages and clients",
            de: "Verwalte Organisationseinstellungen, Mitglieder, Sprachen und Clients"
        },
        icon: "setup",
        run(vue) {
            vue.$store.dispatch("pushColumn", "administration");
        }
    },
    {
        command: {
            en: ["settings"],
            de: ["einstellungen"]
        },
        description: {
            en: "Manage profile and notification settings",
            de: "Bearbeite deine Profil- und Benachrichtigungseinstellungen"
        },
        icon: "settings",
        run(vue) {
            vue.$store.dispatch("pushColumn", "settings");
        }
    },
    {
        command: {
            en: ["billing", "expert"],
            de: ["abrechnung", "expert"]
        },
        description: {
            en: "Show billing details and settings",
            de: "Zeige Zahlungsinformationen und Einstellungen"
        },
        icon: "billing",
        run(vue) {
            vue.$router.push("/billing");
            vue.$store.dispatch("resetCmd");
        }
    },
    {
        command: {
            en: ["help", "usage"],
            de: ["hilfe"]
        },
        description: {
            en: "Show help",
            de: "Zeige die Hilfe an"
        },
        icon: "help",
        run(vue) {
            vue.$store.dispatch("pushColumn", "help");
        }
    },
    {
        command: {
            en: ["support"],
            de: ["support"]
        },
        description: {
            en: "Get help",
            de: "Erhalte Hilfe"
        },
        icon: "support",
        run(vue) {
            vue.$store.dispatch("pushColumn", "support");
        }
    },
    {
        command: {
            en: ["organizations"],
            de: ["organisationen"]
        },
        description: {
            en: "Navigate to organization overview",
            de: "Navigiere zur Organisationsübersicht"
        },
        icon: "organization",
        run(vue) {
            window.location = `${EMVI_WIKI_WEBSITE_HOST}/organizations`;
        }
    },
    {
        command: {
            en: ["logout"],
            de: ["abmelden", "logout"]
        },
        description: {
            en: "Log out of Emvi",
            de: "Von Emvi abmelden"
        },
        icon: "logout",
        run(vue) {
            vue.$store.dispatch("logout");
        }
    }
];

// filterCommands returns the commands for the active page name and global commands filtered using given query.
// The query starts with a slash (for command) followed by a command name or alias.
export function filterCommands(activePath, query, locale) {
    if(!query) {
        query = "";
    }

    query = getQueryFromCmd(query);
    let out = [];

    for(let i = 0; i < commands.length; i++) {
        if((!commands[i].pageNames || commands[i].pageNames.length === 0 || commands[i].pageNames.indexOf(activePath) > -1) &&
            (!commands[i].excludePageNames || commands[i].excludePageNames.length === 0 || commands[i].excludePageNames.indexOf(activePath) < 0) &&
            nameStartsWith(commands[i].command[locale], query)) {
            out.push(commands[i]);
        }
    }

    return out;
}

// findCommandAlias returns the first alias for given query. The alias is read from left to right.
export function findCommandAlias(command, query, locale) {
    query = getQueryFromCmd(query);

    for(let i = 0; i < command.command[locale].length; i++) {
        if(command.command[locale][i].startsWith(query)) {
            return command.command[locale][i];
        }
    }

    return "";
}
