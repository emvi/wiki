const defaultPageTitle = "Emvi";

export const TitleMixin = {
    beforeMount() {
        let title = this.$t("title");

        if(title !== "title") {
            document.title = `${title} — ${defaultPageTitle}`;
        }
        else {
            document.title = defaultPageTitle;
        }
    }
};

export function setPageTitle(title) {
    document.title = `${title} — ${defaultPageTitle}`;
}
