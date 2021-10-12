export const EnlargeMixin = {
    methods: {
        enlargeImage(node) {
            while(this.$refs.enlargeContent.lastChild) {
                this.$refs.enlargeContent.removeChild(this.$refs.enlargeContent.lastChild);
            }

            this.$refs.enlargeContent.appendChild(node.cloneNode(true));
            this.enlarge = true;
        },
        enlargeTable(node) {
            while(this.$refs.enlargeContent.lastChild) {
                this.$refs.enlargeContent.removeChild(this.$refs.enlargeContent.lastChild);
            }

            node = node.cloneNode(true);
            node.removeChild(node.firstChild);
            this.$refs.enlargeContent.appendChild(node);
            this.enlarge = true;
        }
    }
};
