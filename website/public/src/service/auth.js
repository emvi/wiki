import axios from "axios";

export const AuthService = new class {
	getUser() {
		return new Promise((resolve, reject) => {
			// call the actual backend server to receive user information
			// this ensures the server can update the user information when changed
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/user`)
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	register(email) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_AUTH_HOST}/api/v1/auth/registration`, {email})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	cancelRegistration(code) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_AUTH_HOST}/api/v1/auth/registration`, {params: {code}})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getRegistrationStep(code) {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_AUTH_HOST}/api/v1/auth/registration`, {params: {code}})
			.then(r => {
				resolve(r.data.step);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	setRegistrationPassword(code, password, password_repeat) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_AUTH_HOST}/api/v1/auth/registration/password`, {code, password, password_repeat})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	setRegistrationPersonalData(code, firstname, lastname) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_AUTH_HOST}/api/v1/auth/registration/personal`, {code, firstname, lastname})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	setRegistrationCompletion(code, accept_terms_of_service, accept_privacy, accept_marketing, recaptcha_token) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_AUTH_HOST}/api/v1/auth/registration/completion`, {code, accept_terms_of_service, accept_privacy, accept_marketing, recaptcha_token})
			.then(r => {
				resolve({
					access_token: r.data.access_token,
					expires_in: r.data.expires_in,
					secure: r.data.secure,
					domain: r.data.domain
				});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	uploadPicture(form) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/user/picture`, form)
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			})
		});
	}

	deletePicture() {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/user/picture`)
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	savePersonalData(firstname, lastname, language, accept_marketing) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_AUTH_HOST}/api/v1/auth/user/data`, {
				firstname,
				lastname,
				language,
				accept_marketing
			})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	saveMail(email) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_AUTH_HOST}/api/v1/auth/user/email`, {email})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	savePassword(old_password, new_password, new_password_repeat) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_AUTH_HOST}/api/v1/auth/user/password`, {
				old_password,
				new_password,
				new_password_repeat
			})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	saveColorMode(color_mode) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/user/colormode`, {color_mode: parseInt(color_mode)})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}
};
