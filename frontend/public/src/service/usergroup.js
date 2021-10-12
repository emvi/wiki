import axios from "axios";

export const UsergroupService = new class {
	getUsergroup(group_id) {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/usergroup/${group_id}`)
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	saveUsergroup(id, name, info) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/usergroup`, {id, name, info})
			.then(r => {
				resolve(r.data.id);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	deleteUsergroup(id) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/usergroup/${id}`)
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	toggleModerator(group_id, member_id) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/usergroup/${group_id}/member`, {member_id})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	addMember(group_id, user_ids, group_ids) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/usergroup/${group_id}/member`, {user_ids, group_ids})
			.then(r => {
				let member = r.data || [];

				// remove group ID because this might cause trouble rendering users and groups in one list
				for(let i = 0; i < member.length; i++) {
					member[i].user_group_id = undefined;
				}

				resolve(member);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	removeMember(group_id, member_id) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/usergroup/${group_id}/member`, {params: {member_ids: member_id}})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getMember(group_id, filter, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		if(!filter){
			filter = {};
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/usergroup/${group_id}/member`, {params: filter, cancelToken})
			.then(r => {
				resolve({results: r.data.member || [], count: r.data.count});
			})
			.catch(e => {
				reject(e);
			});
		});
	}
};
