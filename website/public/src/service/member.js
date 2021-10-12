import axios from "axios";

export const MemberService = new class {
	joinOrganization(username, code, invitation_code) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/accession`, {username, code, invitation_code})
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	getInvitation(code) {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/invitation/${code}`)
			.then((r) => {
				resolve(r.data.organization);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	getInvitationOrganization(code) {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/invitation/organization`, {params: {code}})
				.then((r) => {
					resolve(r.data.organization);
				})
				.catch((e) => {
					reject(e);
				});
		});
	}
};
