import axios from "axios";

export const ClientService = new class {
    getClients() {
        return new Promise((resolve, reject) => {
            axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/client`)
            .then(r => {
                resolve(r.data || []);
            })
            .catch(e => {
                reject(e);
            });
        });
    }

    getClient(id) {
        return new Promise((resolve, reject) => {
            axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/client/${id}`)
            .then(r => {
                resolve(r.data);
            })
            .catch(e => {
                reject(e);
            });
        });
    }

    saveClient(id, name, scopes) {
        if(!id) {
            id = "";
        }

        return new Promise((resolve, reject) => {
            axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/client`, {id, name, scopes})
            .then(r => {
                resolve(r);
            })
            .catch(e => {
                reject(e);
            });
        });
    }

    deleteClient(id) {
        return new Promise((resolve, reject) => {
            axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/client/${id}`)
            .then(r => {
                resolve(r);
            })
            .catch(e => {
                reject(e);
            });
        });
    }
};
