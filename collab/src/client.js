const axios = require("axios");
const {logger} = require("./logger.js");

const AUTH_HOST = process.env.EMVI_WIKI_AUTH_HOST || "https://auth.emvi.com";
const BACKEND_HOST = process.env.EMVI_WIKI_BACKEND_HOST || "https://api.emvi.com";
const GRANT_TYPE = "client_credentials";
const CLIENT_ID = process.env.EMVI_WIKI_AUTH_CLIENT_ID;
const CLIENT_SECRET = process.env.EMVI_WIKI_AUTH_CLIENT_SECRET;
const TOKEN_ENDPOINT = "/api/v1/auth/token";

module.exports.apiClient = new class {
    constructor() {
        this.token_type = "";
        this.token = "";
        this.expires_in = 0;
        this.http = axios.create();
        this.addAxiosInterceptor();
        this.refreshTokenPromise = null;
    }

    refreshToken() {
        if(this.refreshTokenPromise) {
            return this.refreshTokenPromise;
        }

        this.refreshTokenPromise = new Promise((resolve, reject) => {
            let req = {
                grant_type: GRANT_TYPE,
                client_id: CLIENT_ID,
                client_secret: CLIENT_SECRET
            };

            this.http.post(AUTH_HOST+TOKEN_ENDPOINT, req)
            .then(r => {
                logger.info("Obtained client token from server");
                this.token_type = r.data.token_type;
                this.token = r.data.access_token;
                this.expires_in = parseInt(r.data.expires_in);
                this.refreshTokenPromise = null;
                resolve(this.token);
            })
            .catch(e => {
                logger.error("Error obtaining token from auth server", {host: AUTH_HOST, error: e.message});
                reject(e);
            });
        });

        return this.refreshTokenPromise;
    }

    addAxiosInterceptor() {
        this.http.interceptors.response.use(null, e => {
            if(e && e.config && e.response && e.response.status === 401) {
                logger.info("Token expired, refreshing...");

                return this.refreshToken()
                .then(token => {
                    e.config.headers["Authorization"] = `Bearer ${token}`;
                    return this.http.request(e.config);
                })
                .catch(e => {
                    return Promise.reject(e);
                });
            }

            return Promise.reject(e);
        });
    }

    getArticle(client, id, lang_id) {
        return new Promise((resolve, reject) => {
            let config = {
                ...this.getConfig(client),
                params: {
                    version: 99999999, // load latest including wip
                    raw_content: true,
                    lang: lang_id,
                    user_id: client.user.id
                }
            };

            this.http.get(`${BACKEND_HOST}/api/v1/article/${id}`, config)
            .then(r => {
                resolve(r);
            })
            .catch(e => {
                reject(e);
            });
        });
    }

    getOrganization(client) {
        return new Promise((resolve, reject) => {
            let config = {
                ...this.getConfig(client)
            };

            this.http.get(`${BACKEND_HOST}/api/v1/organization`, config)
            .then(r => {
                resolve(r);
            })
            .catch(e => {
                reject(e);
            });
        });
    }

    saveArticle(client, req) {
        return new Promise((resolve, reject) => {
            this.http.post(`${BACKEND_HOST}/api/v1/article`, req, this.getConfig(client))
            .then(r => {
                resolve(r);
            })
            .catch(e => {
                reject(e.response.data);
            });
        });
    }

    addTag(client, tag) {
        return new Promise((resolve, reject) => {
            let config = {
                ...this.getConfig(client),
                params: {tag}
            };

            this.http.get(`${BACKEND_HOST}/api/v1/tag`, config)
            .then(() => {
                resolve();
            })
            .catch(e => {
                reject(e);
            });
        });
    }

    getConfig(client) {
        return {
            headers: {
                "Authorization": "Bearer "+this.token,
                "Organization": client.user.organization,
                "Client": CLIENT_ID
            }
        };
    }
};
