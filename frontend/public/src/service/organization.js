import axios from "axios";

export const OrganizationService = new class {
	getOrganization() {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization`)
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	leaveOrganization(name) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/exit`, {name})
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	deleteOrganization(name) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization`, {params: {name}})
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	uploadPicture(form) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/picture`, form)
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
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/picture`)
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			})
		});
	}

	updateOrganization(name, domain) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization`, {name, domain})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			})
		});
	}

	updateOrganizationPermissions(create_group_admin, create_group_mod) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/permissions`, {create_group_admin, create_group_mod})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			})
		});
	}

	getStatistics() {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/statistics`)
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			})
		});
	}

	generateInvitationCode(read_only) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/invitation`, {read_only})
				.then(r => {
					resolve(r.data.code);
				})
				.catch(e => {
					reject(e);
				})
		});
	}

	getInvitationCode() {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/invitation`)
				.then(r => {
					resolve({code: r.data.code, read_only: r.data.read_only});
				})
				.catch(e => {
					reject(e);
				})
		});
	}
};
