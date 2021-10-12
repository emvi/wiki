import Vue from "vue";
import Vuex from "vuex";
import VueMoment from "vue-moment";
import VueMq from "vue-mq";
import moment from "moment/moment";
import {mapGetters} from "vuex";
import NewStore from "./store/store.js";
import {addAxiosResponseInterceptor, addAxiosTokenInterceptor} from "./axios";
import {toast} from "./components";
import {ErrorService} from "./service";
import {getRouter} from "./router";
import {getVueI18n} from "./i18n";
import "./mixin";
import {getCookie, getLocale} from "./util";

Vue.use(Vuex);
Vue.use(VueMoment, {moment});
Vue.use(VueMq, {
	breakpoints: {
		512: 512,
		768: 768,
		1024: 1024,
		1025: Infinity,
	}
});
Vue.config.productionTip = false;
Vue.config.devtools = false;
addAxiosResponseInterceptor();
addAxiosTokenInterceptor();
let store = NewStore();
Vue.component("toast", toast);

new Vue({
	el: "#app",
	store,
	router: getRouter(store),
	i18n: getVueI18n(),
	computed: {
		...mapGetters(["user", "language", "colormode"])
	},
	watch: {
		user() {
			this.initLanguage();
			this.setColorMode();
		}
	},
	mounted() {
		this.loadUser();
		this.initLanguage();
		this.initRecaptcha();
		this.completeJoin();
	},
	methods: {
		loadUser() {
			if(getCookie("access_token")) {
				this.$store.dispatch("loadUser");
			}
		},
		initLanguage() {
			let userLanguage = null;

			if(this.language) {
				userLanguage = this.language;
			}

			this.$i18n.locale = getLocale(userLanguage);
			ErrorService.setLocale(this.$i18n.locale);
			this.$moment.locale(this.$i18n.locale);
		},
		initRecaptcha() {
			let self = this;
			window.onload = function() {
				self.$store.commit("setRecaptcha", grecaptcha);
			};
		},
		setColorMode() {
			let root = document.documentElement;

			if(this.colormode === 1) {
				root.className = "theme-light theme-pro-deu";
			}
			else if(this.colormode === 2) {
				root.className = "theme-light theme-tri";
			}
			else {
				root.className = "theme-light";
			}
		},
		completeJoin() {
			if(this.$route.path !== "/join") {
				let code = localStorage.getItem("join_code");

				if(code) {
					this.$router.push(`join?code=${code}`);
				}
			}

			localStorage.removeItem("join_code");
		}
	}
});
