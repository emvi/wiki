import axios from "axios";

export const MemberService = new class {
	toggleModerator(user_id) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/member/${user_id}/moderator`)
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	toggleAdmin(user_id) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/member/${user_id}/admin`)
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	toggleReadOnly(user_id) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/member/${user_id}/readonly`)
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	removeMember(user_id, remove_permissions) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/member/${user_id}`, {params: {remove_permissions}})
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	inviteMember(emails, read_only) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/member`, {emails, read_only})
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	cancelInvitation(invitation_id) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/member`, {params: {invitation_id}})
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	getInvitations() {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/member`)
			.then((r) => {
				resolve(r.data || []);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}
};
