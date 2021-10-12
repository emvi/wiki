const languages = [
    "en",
    "de"
];
const defaultLocale= "en";

// Returns the users locale as two character lowercase ISO string or the one passed, if valid.
export function getLocale(userLocale) {
    if(userLocale) {
        userLocale = userLocale.toLowerCase();

        if (languages.indexOf(userLocale) !== -1) {
            return userLocale;
        }
    }

    let langs = navigator.languages.slice();

    for(let i = 0; i < langs.length; i++) {
        langs[i] = langs[i].substr(0, 2).toLowerCase(); // e.g.: en-US -> en

        for(let j = 0; j < languages.length; j++) {
            if(langs[i] === languages[j]) {
                return langs[i];
            }
        }
    }

    return defaultLocale;
}

// Returns the system supported locale or default if not available.
export function getSupportedLocale(userLocale) {
    userLocale = userLocale.toLowerCase();

    if(languages.indexOf(userLocale) !== -1) {
        return userLocale
    }

    return defaultLocale;
}
