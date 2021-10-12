import {slugWithId} from "../util";

export const MentionMixin = {
    methods: {
        followMention(target, mention, button) {
            let type = target.getAttribute("object");

            if(type === "user") {
                this.openLink(button, `/member/${mention}`);
            }
            else if(type === "group") {
                let title = target.getAttribute("title");
                this.openLink(button, `/group/${slugWithId(title, mention)}`);
            }
            else if(type === "list") {
                let title = target.getAttribute("title");
                this.openLink(button, `/list/${slugWithId(title, mention)}`);
            }
            else if(type === "tag") {
                this.openLink(button, `/tag/${mention}`);
            }
            else {
                let title = target.getAttribute("title");
                this.openLink(button, `/read/${slugWithId(title, mention)}`);
            }
        },
        openLink(button, url) {
            if(button === 0) {
                this.$router.push(url);
            }
        }
    }
};
