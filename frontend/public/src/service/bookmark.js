import axios from "axios";

export const BookmarkService = new class {
	bookmark({article_id, article_list_id}) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/bookmark`, {article_id, article_list_id})
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getBookmarks(articles, lists, offset_articles, offset_lists, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		if(!offset_articles) {
			offset_articles = 0;
		}

		if(!offset_lists) {
			offset_lists = 0;
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/bookmark`, {params: {articles, lists, offset_articles, offset_lists}, cancelToken})
			.then(r => {
				resolve({articles: r.data.articles || [], lists: r.data.lists || []});
			})
			.catch(e => {
				reject(e);
			});
		});
	}
};
