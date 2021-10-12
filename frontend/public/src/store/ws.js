import io from "socket.io-client";
import {getSubdomain} from "../util";

export const WSStore = {
    state: {
        socket: null
    },
    mutations: {
        connect(state, token) {
            state.socket = io(EMVI_WIKI_COLLAB_HOST, {
                query: {
                    access_token: token,
                    organization: getSubdomain()
                },
                path: "/api/v1/collab/ws",
                timeout: 10000
            });
        }
    },
    actions: {
        connect(context, token) {
            context.commit("connect", token);
        }
    }
};
