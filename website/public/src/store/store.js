import Vuex from "vuex";
import {UserStore} from "./user.js";
import {ToastStore} from "./toast.js";
import {RecaptchaStore} from "./recaptcha.js";

export default function NewStore() {
	return new Vuex.Store({
		modules: {
			user: UserStore,
			toast: ToastStore,
			recaptcha: RecaptchaStore
		}
	});
}
