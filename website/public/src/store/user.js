import {AuthService} from "../service";

export const UserStore = {
	state: {
		user: JSON.parse(window.localStorage.getItem("user"))
	},
	mutations: {
		setUser(state, user) {
			if(user) {
				state.user = user;
				window.localStorage.setItem("user", JSON.stringify(user));
			}
			else {
				state.user = null;
				window.localStorage.removeItem("user");
			}
		}
	},
	actions: {
		login(context, stateParam) {
			if(!stateParam) {
				window.location = `${EMVI_WIKI_AUTH_HOST}/auth/authorize?response_type=token&client_id=${EMVI_WIKI_AUTH_CLIENT_ID}&redirect_uri=${EMVI_WIKI_WEBSITE_HOST}/organizations`;
			}
			else{
				window.location = `${EMVI_WIKI_AUTH_HOST}/auth/authorize?response_type=token&client_id=${EMVI_WIKI_AUTH_CLIENT_ID}&redirect_uri=${EMVI_WIKI_WEBSITE_HOST}/organizations&state=${stateParam}`;
			}
		},
		logout(context) {
			context.commit("setUser", null);
			window.location = `${EMVI_WIKI_AUTH_HOST}/auth/logout?redirect_uri=${EMVI_WIKI_WEBSITE_HOST}`;
		},
		loadUser(context) {
			AuthService.getUser()
				.then(user => {
					context.commit("setUser", user);
				})
				.catch(e => {
					console.error(e);
					context.dispatch("logout");
				});
		}
	},
	getters: {
		user(state) {
			return state.user || {};
		},
		loggedIn(state) {
			return state.user !== undefined && state.user !== null;
		},
		language(state) {
			return state.user ? state.user.language : null;
		},
		colormode(state) {
			return state.user ? state.user.color_mode : null;
		}
	}
};
