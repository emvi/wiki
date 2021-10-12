import slugify from "slugify";

// this is a hack to change the URL without reloading the page, which sucks because it ignores vue-router
// https://github.com/vuejs/vue-router/issues/703
export function setURL(url) {
    window.history.pushState(null, null, url);
}

// Returns the ID for given slug separated by hyphen.
// The ID must be the last part of the parameter.
export function getIdFromSlug(slug) {
    if(typeof slug !== "string") {
        return slug;
    }

    let parts = slug.split("-");
    return parts[parts.length-1];
}

// Appends the given ID to given string and slugifys everything.
export function slugWithId(str, id) {
    return `${slugify(str)}-${id}`;
}
