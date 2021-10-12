export function getCookie(name) {
    name += "=";
    let ca = document.cookie.split(';');

    for(let i = 0; i < ca.length; i++) {
        let c = ca[i];

        while(c.charAt(0) === ' ') {
            c = c.substring(1,c.length);
        }

        if(c.indexOf(name) === 0) {
            return c.substring(name.length, c.length);
        }
    }

    return null;
}

export function setCookie(name, value, expires, secure, domain) {
    let time = new Date();
    time.setSeconds(time.getSeconds()+expires);
    let cookie = `${name}=${value}; expires=${time.toUTCString()}; domain=${domain}; path=/; samesite=strict;`;

    if(secure) {
        cookie += " secure;";
    }

    document.cookie = cookie;
}
