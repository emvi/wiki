export function initialsFilter(user) {
    if(user && user.firstname && user.lastname && user.firstname.length && user.lastname.length) {
        return `${user.firstname.charAt(0)}${user.lastname.charAt(0)}`.toUpperCase();
    }

    return "";
}
