import axios from "axios";

export const InvitationService = new class {
	getInvitations() {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/invitation`)
			.then(r => {
				resolve(r.data.invitations || []);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	deleteInvitation(id) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/invitation`, {params: {id}})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}
};
