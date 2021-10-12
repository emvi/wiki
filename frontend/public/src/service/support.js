import axios from "axios";

export const SupportService = new class {
    contact(type, subject, message) {
        return new Promise((resolve, reject) => {
            axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/support`, {type, subject, message})
            .then(r => {
                resolve(r.data);
            })
            .catch(e => {
                reject(e);
            });
        });
    }
};
