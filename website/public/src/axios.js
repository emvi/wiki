import axios from "axios";
import {getCookie} from "./util";

export function addAxiosResponseInterceptor(store) {
    axios.interceptors.response.use(undefined, err => {
        if (axios.isCancel(err)) {
            return Promise.reject(err);
        }

        if (err.response.status === 401) {
            store.dispatch("logout");
            return;
        }

        let errors = {
            validation: err.response.data.validation,
            errors: err.response.data.errors,
            exceptions: err.response.data.exceptions
        };

        return Promise.reject(errors);
    });
}

export function addAxiosTokenInterceptor() {
    axios.interceptors.request.use(config => {
        const token = getCookie("access_token");

        if(token) {
            config.headers.Authorization = `Bearer ${token}`;
        }

        return config;
    }, err => {
        return Promise.reject(err);
    });
}
