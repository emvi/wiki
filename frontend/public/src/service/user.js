import axios from "axios";

export const UserService = new class {
	getMember() {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/user/member`)
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	saveMember(member) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/user/member`, member)
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getUser(id, name) {
		let endpoint = id ? `/api/v1/profile/${id}` : "/api/v1/profile";

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}${endpoint}`, {params: {name}})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	setIntroduction(introduction) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/user/introduction`, {introduction})
				.then(r => {
					resolve(r.data);
				})
				.catch(e => {
					reject(e);
				});
		});
	}
};
