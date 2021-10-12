import {BookmarkService, ObserveService, PinnedService} from "../../service";

// toggleObserve toggles observed object on/off.
// The type must be either article, list or group.
export function toggleObserve(vue, type) {
    let id = vue.$store.state.page.meta.get("id");
    let data = {};

    if(type === "article") {
        data.article_id = id;
    }
    else if(type === "list") {
        data.article_list_id = id;
    }
    else {
        data.user_group_id = id;
    }

    ObserveService.observe(data)
        .then(() => {
            let observed = vue.$store.state.page.meta.get("observed");
            vue.$store.dispatch("setMeta", {key: "observed", value: !observed});
            vue.$store.dispatch("success", vue.$t(observed ? `observe_${type}_off` : `observe_${type}_on`));
            vue.$store.dispatch("resetCmd");
        })
        .catch(e => {
            vue.setError(e);
        });
}

// toggleBookmark toggles bookmarked object on/off.
// The type must be either article or list.
export function toggleBookmark(vue, type) {
    let id = vue.$store.state.page.meta.get("id");
    let data = {};

    if(type === "article") {
        data.article_id = id;
    }
    else {
        data.article_list_id = id;
    }

    BookmarkService.bookmark(data)
        .then(() => {
            let bookmarked = vue.$store.state.page.meta.get("bookmarked");
            vue.$store.dispatch("setMeta", {key: "bookmarked", value: !bookmarked});
            vue.$store.dispatch("success", vue.$t(bookmarked ? `bookmark_${type}_off` : `bookmark_${type}_on`));
            vue.$store.dispatch("resetCmd");
        })
        .catch(e => {
            vue.setError(e);
        });
}

// togglePinned toggles pinned object on/off.
// The type must be either article or list.
export function togglePinned(vue, type) {
    if(!vue.isMod) {
        return;
    }

    let id = vue.$store.state.page.meta.get("id");
    let data = {};

    if(type === "article") {
        data.article_id = id;
    }
    else {
        data.article_list_id = id;
    }

    PinnedService.pin(data)
        .then(() => {
            let pinned = vue.$store.state.page.meta.get("pinned");
            vue.$store.dispatch("setMeta", {key: "pinned", value: !pinned});
            vue.$store.dispatch("success", vue.$t(pinned ? `pin_${type}_off` : `pin_${type}_on`));
            vue.$store.dispatch("resetCmd");
        })
        .catch(e => {
            vue.setError(e);
        });
}
