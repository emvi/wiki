import axios from "axios";

export const ObserveService = new class {
	observe({article_id, article_list_id, user_group_id}) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/observe`, {article_id, article_list_id, user_group_id})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getObserved(articles, lists, groups, offset_articles, offset_lists, offset_groups, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		let params = {
			articles,
			lists,
			groups,
			offset_articles,
			offset_lists,
			offset_groups
		};

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/observe`, {params, cancelToken})
			.then(r => {
				resolve({articles: r.data.articles || [], lists: r.data.lists || [], groups: r.data.groups || []});
			})
			.catch(e => {
				reject(e);
			});
		});
	}
};
