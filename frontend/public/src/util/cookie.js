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
