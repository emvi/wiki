import axios from "axios";

export const AuthService = new class {
	getUser() {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/user`)
				.then(r => {
					resolve(r.data);
				})
				.catch(e => {
					reject(e);
				});
		});
	}
};
