import axios from "axios";

export const SearchService = new class {
	findUserAndUserGroup(query, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/search/userusergroup`, {params: {query}, cancelToken})
			.then(r => {
				let results = r.data;
				let out = [];

				for(let i = 0; i < results.length; i++) {
					if(results[i].user) {
						results[i].user.type = "user";
						out.push(results[i].user);
					}
					else {
						results[i].usergroup.type = "group";
						out.push(results[i].usergroup);
					}
				}

				resolve(out);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	findTag(query, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/search/tag`, {params: {query}, cancelToken})
			.then(r => {
				resolve({tags: r.data.tags || []});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	findArticles(query, filter, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		if(!filter) {
			filter = {};
		}

		filter.query = query;

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/search/article`, {params: filter, cancelToken})
			.then(r => {
				resolve({results: r.data.articles || [], count: r.data.count});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	findArticleLists(query, filter, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		if(!filter) {
			filter = {};
		}

		filter.query = query;

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/search/list`, {params: filter, cancelToken})
			.then(r => {
				resolve({results: r.data.lists || [], count: r.data.count});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	findUser(query, filter, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		if(!filter) {
			filter = {};
		}

		filter.query = query;

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/search/user`, {params: filter, cancelToken})
			.then(r => {
				resolve({user: r.data.user || [], count: r.data.count});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	findUserGroups(query, filter, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		if(!filter) {
			filter = {};
		}

		filter.query = query;
		filter.find_groups = true;

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/search/usergroup`, {params: filter, cancelToken})
			.then(r => {
				resolve({results: r.data.groups || [], count: r.data.count});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	findTags(query, filter, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		if(!filter) {
			filter = {};
		}

		filter.query = query;

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/search/tag`, {params: filter, cancelToken})
			.then(r => {
				resolve({tags: r.data.tags || [], count: r.data.count});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	findAll(query, filter, limit, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		if(!limit) {
			limit = 0;
		}

		if(!filter) {
			filter = {
				articles: true,
				lists: true,
				user: true,
				groups: true,
				tags: true,
				feed: true,
				articles_limit: limit,
				lists_limit: limit,
				user_limit: limit,
				groups_limit: limit,
				tags_limit: limit,
				feed_limit: limit
			};
		}

		filter.query = query;

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/search`, {params: filter, cancelToken})
			.then(r => {
				resolve(r.data || []);
			})
			.catch(e => {
				reject(e);
			});
		});
	}
};
