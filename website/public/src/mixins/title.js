export default {
    mounted() {
        this.setPageTitle(this.$t("pagetitle"));
    },
    methods: {
        setPageTitle(title) {
            if(title && title !== "pagetitle") {
                document.title = `${title} — ${document.title}`;
            }
        }
    }
}
