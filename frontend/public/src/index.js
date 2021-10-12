import Vue from "vue";
import VueRouter from "vue-router";
import Vuex, {mapGetters} from "vuex";
import VueI18n from "vue-i18n";
import VueMoment from "vue-moment";
import VueMq from "vue-mq";
import moment from "moment/moment";
import InfiniteLoading from "vue-infinite-loading";
import {addAxiosResponseInterceptor, addAxiosTokenInterceptor} from "./axios";
import {getRouter} from "./router";
import {getVueI18n, setMomentAbbreviations} from "./i18n";
import {getCookie, getLocale, getSupportedLocale, isEmptyObject, isObject} from "./util";
import {initialsFilter} from "./filter";
import {ErrorService} from "./service";
import {NotificationsStore, UserStore, WSStore} from "./store";
import {CmdStore} from "./components/cmd/state";
import {PageStore} from "./components/layout/state";
import {ToastStore} from "./components/toast/state";

Vue.use(VueRouter);
Vue.use(Vuex);
Vue.use(VueI18n);
Vue.use(VueMoment, {moment});
Vue.use(VueMq, {
	breakpoints: {
		512: 512,
		768: 768,
		1024: 1024,
		1200: 1200,
		1440: 1440,
		1600: Infinity
	}
});
Vue.filter("initials", initialsFilter);
Vue.use(InfiniteLoading, {
	// disable all slots (we use custom components to show states)
	slots: {
		spinner: {
			render: h => h('div'),
		},
		noResults: {
			render: h => h('div'),
		},
		noMore: {
			render: h => h('div'),
		},
		error: {
			render: h => h('div'),
		}
	}
});
Vue.config.productionTip = false;
Vue.config.devtools = false;

let store = new Vuex.Store({
	modules: {
		user: UserStore,
		ws: WSStore,
		notifications: NotificationsStore,
		cmd: CmdStore,
		page: PageStore,
		toast: ToastStore
	}
});

addAxiosResponseInterceptor(store);
addAxiosTokenInterceptor();
setMomentAbbreviations();

Vue.mixin({
	data() {
		return {
			validation: {},
			err: null
		};
	},
	computed: {
		locale() {
			return getSupportedLocale(this.$i18n.locale);
		},
		isExpert() {
			return this.$store.state.user.organization.expert;
		},
		isAdmin() {
			return this.$store.state.user.member.is_admin;
		},
		isMod() {
			return this.$store.state.user.member.is_moderator;
		},
		isReadOnly() {
			return this.$store.state.user.member.read_only;
		},
		isMobile() {
			return this.$mq <= 768;
		},
		isTouch() {
			return "ontouchstart" in document.documentElement;
		}
	},
	methods: {
		focusCmdInput(ignoreTouch = false) {
			if(!this.isTouch || ignoreTouch) {
				document.getElementById("cmd-input").focus();
			}
		},
		resetError() {
			this.validation = {};
			this.err = null;
			this.$store.dispatch("setError");
		},
		setError(e) {
			if(!e) {
				return;
			}

			if(isObject(e.validation) && !isEmptyObject(e.validation)) {
				this.validation = ErrorService.map(e.validation);
			}
			else if(e.errors && e.errors.length) {
				// display the first error only
				this.err = ErrorService.mapError(e.errors[0].message);
				this.$store.dispatch("setError", this.err);
			}
			else if(e.exceptions && e.exceptions.length) {
				// display the first exception only
				this.err = ErrorService.mapError(e.exceptions[0]);
			}
			else {
				this.err = e;
			}
		},
		showTechnicalError(e) {
			console.error(e);

			if(e.status === 400) {
				this.$store.dispatch("error", `${e.status}: ${JSON.stringify(e.validation)} ${e.errors.join(", ")} ${e.exceptions.join(", ")}`);
			}
			else {
				this.$store.dispatch("error", `An unexpected error occured (status code ${e.status}).`);
			}
		}
	}
});

new Vue({
	el: "#app",
	store,
	router: getRouter(store),
	i18n: getVueI18n(),
	computed: {
		...mapGetters(["user", "darkmode", "colormode"])
	},
	watch: {
		user(user) {
			if(user.introduction) {
				this.$router.push("/intro");
			}
		},
		darkmode() {
			this.toggleColorMode();
		},
		colormode() {
			this.toggleColorMode();
		}
	},
	mounted() {
		this.$store.dispatch("reload");
		this.initLanguage();
		this.connectToCollab();
		this.toggleColorMode();
	},
	methods: {
		initLanguage() {
			let userLanguage = null;

			if(this.user && this.user.language) {
				userLanguage = this.user.language;
			}

			this.$i18n.locale = getLocale(userLanguage);
			ErrorService.setLocale(this.$i18n.locale);
			this.$moment.locale(this.$i18n.locale);
		},
		connectToCollab() {
			this.$store.dispatch("connect", getCookie("access_token"));
		},
		toggleColorMode() {
			let className = this.darkmode ? "theme-dark" : "theme-light";

			if(this.colormode === 1) {
				className += " theme-pro-deu";
			}
			else if(this.colormode === 2) {
				className += " theme-tri";
			}

			document.documentElement.className = className;
		}
	}
});
