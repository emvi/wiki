import {AuthService, OrganizationService, UserService} from "../service";
import {getSubdomain} from "../util";

export const UserStore = {
	state: {
		user: JSON.parse(window.localStorage.getItem("user")) || {},
		member: JSON.parse(window.localStorage.getItem("member")) || {},
		organization: JSON.parse(window.localStorage.getItem("organization")) || {},
		organizationStatistics: {},
		darkmode: JSON.parse(window.localStorage.getItem("darkmode")) || false
	},
	mutations: {
		invalidateAuth(state) {
			window.localStorage.removeItem("user");
			window.localStorage.removeItem("member");
			window.localStorage.removeItem("organization");
			state.user = {};
			state.member = {};
			state.organization = {};
		},
		loadUser(state) {
			AuthService.getUser()
				.then(user => {
					window.localStorage.setItem("user", JSON.stringify(user));
					state.user = user;
				})
				.catch(e => {
					console.error(e);
					this.dispatch("redirectLogin");
				});
		},
		loadMember(state) {
			UserService.getMember()
				.then(member => {
					window.localStorage.setItem("member", JSON.stringify(member));
					state.member = member;
				})
				.catch(e => {
					console.error(e);
					this.dispatch("redirectLogin");
				});
		},
		loadOrganization(state) {
			OrganizationService.getOrganization()
				.then(organization => {
					window.localStorage.setItem("organization", JSON.stringify(organization));
					state.organization = organization;
				})
				.catch(e => {
					console.error(e);
					this.dispatch("redirectLogin");
				});
		},
		loadOrganizationStatistics(state) {
			OrganizationService.getStatistics()
				.then(statistics => {
					state.organizationStatistics = statistics;
				})
				.catch(e => {
					this.setError(e);
				});
		},
		setDarkmode(state, darkmode) {
			state.darkmode = darkmode;
			window.localStorage.setItem("darkmode", JSON.stringify(darkmode));
		}
	},
	actions: {
		logout(context) {
			context.commit("invalidateAuth");
			window.location = `${EMVI_WIKI_AUTH_HOST}/auth/logout?redirect_uri=${EMVI_WIKI_WEBSITE_HOST}`;
		},
		redirectOrganizations(context) {
			context.commit("invalidateAuth");

			// redirect to website and logout there or show organization overview in case something went wrong
			window.location = `${EMVI_WIKI_WEBSITE_HOST}/organizations`;
		},
		redirectLogin() {
			window.location = `${EMVI_WIKI_AUTH_HOST}/auth/authorize?response_type=token&client_id=${EMVI_WIKI_AUTH_CLIENT_ID}&redirect_uri=${EMVI_WIKI_WEBSITE_HOST}/organizations&state=${getSubdomain()}`;
		},
		toggleDarkmode(context) {
			context.commit("setDarkmode", !context.state.darkmode);
		},
		loadOrganization(context) {
			context.commit("loadOrganization");
		},
		loadOrganizationStatistics(context) {
			context.commit("loadOrganizationStatistics");
		},
		reload(context) {
			context.commit("loadUser");
			context.commit("loadMember");
			context.commit("loadOrganization");
		}
	},
	getters: {
		user(state) {
			return state.user;
		},
		member(state) {
			return state.member;
		},
		organization(state) {
			return state.organization;
		},
		organizationStatistics(state) {
			return state.organizationStatistics;
		},
		darkmode(state) {
			return state.darkmode;
		},
		colormode(state) {
			return state.user.color_mode;
		}
	}
};
