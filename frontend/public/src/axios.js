import axios from "axios";
import {getCookie, getSubdomain} from "./util";

const subdomain = getSubdomain();

export function addAxiosResponseInterceptor(store) {
    axios.interceptors.response.use(undefined, err => {
        if(axios.isCancel(err)) {
            return Promise.reject(err);
        }

        if(err.response.status === 401) {
            store.dispatch("redirectLogin");
            return;
        }

        return Promise.reject({
            status: err.response.status,
            validation: err.response.data.validation,
            errors: err.response.data.errors,
            exceptions: err.response.data.exceptions
        });
    });
}

export function addAxiosTokenInterceptor() {
    axios.interceptors.request.use(config => {
        const token = getCookie("access_token");

        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
            config.headers.Organization = subdomain;
        }

        return config;
    }, err => {
        return Promise.reject(err);
    });
}
