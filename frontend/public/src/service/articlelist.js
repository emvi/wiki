import axios from "axios";

export const ArticlelistService = new class {
	getArticlelist(list_id) {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${list_id}`)
			.then(r => {
				resolve({list: r.data.list, moderator: r.data.moderator, observed: r.data.observed, bookmarked: r.data.bookmarked});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	saveArticlelist(id, names, is_public, client_access) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist`, {id, names, public: is_public, client_access})
			.then(r => {
				resolve(r.data.id);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	deleteArticlelist(id) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${id}`)
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	addEntry(list_id, article_ids) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${list_id}/entry`, {article_ids})
			.then(r => {
				let articles = [];

				for(let i = 0; i < r.data.length; i++) {
					articles.push(r.data[i].article);
				}

				resolve(articles);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	removeEntry(list_id, article_id) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${list_id}/entry`, {params: {article_ids: article_id}})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getEntries(list_id, filter, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		if(!filter){
			filter = {};
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${list_id}/entry`, {params: filter, cancelToken})
			.then(r => {
				resolve({results: r.data.entries || [], count: r.data.count, start_pos: r.data.start_pos});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	swapEntries(list_id, article_id_a, article_id_b) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${list_id}/entry`, {article_id_a, article_id_b})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	moveEntry(list_id, article_id, direction) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${list_id}/entry`, {article_id, direction})
				.then(r => {
					resolve(r);
				})
				.catch(e => {
					reject(e);
				});
		});
	}

	getMember(list_id, filter, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		if(!filter){
			filter = {};
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${list_id}/member`, {params: filter, cancelToken})
			.then(r => {
				let results = r.data.member || [];
				let out = [];

				for(let i = 0; i < results.length; i++) {
					if(results[i].user_id) {
						results[i].user.type = "user";
						results[i].user.member_id = results[i].id;
						results[i].user.is_moderator = results[i].is_moderator;
						out.push(results[i].user);
					}
					else {
						results[i].user_group.type = "group";
						results[i].user_group.member_id = results[i].id;
						out.push(results[i].user_group);
					}
				}

				resolve({results: out, count: r.data.count});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	addMember(list_id, user_ids, group_ids) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${list_id}/member`, {user_ids, group_ids})
			.then(r => {
				resolve(r.data || []);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	removeMember(list_id, member_id) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${list_id}/member`, {params: {member_ids: member_id}})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	toggleModerator(list_id, member_id) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/${list_id}/member`, {member_id})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getPrivateLists(offset, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/articlelist/private`, {params: {offset}, cancelToken})
			.then(r => {
				resolve(r.data || []);
			})
			.catch(e => {
				reject(e);
			});
		});
	}
};
