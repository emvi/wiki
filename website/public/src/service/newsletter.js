import axios from "axios";

export const NewsletterService = new class {
    confirm(code) {
        return new Promise((resolve, reject) => {
            axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/newsletter`, {code})
            .then(r => {
                resolve(r.data);
            })
            .catch((e) => {
                reject(e);
            });
        });
    }

    unsubscribe(code) {
        return new Promise((resolve, reject) => {
            axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/newsletter`, {params: {code}})
            .then(r => {
                resolve(r.data);
            })
            .catch((e) => {
                reject(e);
            });
        });
    }
};
