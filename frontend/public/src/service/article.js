import axios from "axios";
import {tagsToStringArray} from "../util";

export const ArticleService = new class {
	getArticle(id, lang, version, update_views) {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${id}`, {params: {lang, version, update_views}})
			.then(r => {
				let data = r.data || {};

				if(!data.article) {
					data.article = {};
				}

				if(!data.article.tags) {
					data.article.tags = [];
				}
				else{
					data.article.tags = tagsToStringArray(data.article.tags);
				}

				if(!data.authors) {
					data.authors = [];
				}

				resolve(data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getArticleHistory(id, lang, offset, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${id}/history`, {params: {lang, offset}, cancelToken})
			.then(r => {
				resolve({history: r.data.history || [], count: r.data.count});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getArticlePreview(id, lang, preview_paragraph) {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${id}/preview`, {params: {lang, preview_paragraph}})
			.then(r => {
				resolve(r.data.content);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	recommendArticle(article_id, user, groups, message, receive_read_confirmation) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${article_id}/recommendation`, {user, groups, message, receive_read_confirmation})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	inviteArticle(article_id, language_id, room_id, user, groups, message) {
		let url = `${EMVI_WIKI_BACKEND_HOST}/api/v1/article/invite`;

		if(article_id) {
			url = `${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${article_id}/invite`;
		}

		return new Promise((resolve, reject) => {
			axios.put(url, {language_id, room_id, user, groups, message})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	archiveArticle(article_id, message, deleteImmediatly) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${article_id}/archive`, {message, delete: deleteImmediatly})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	deleteArticle(id) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${id}`)
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	resetArticle(article_id, language_id, version, commit) {
		let data = {
			language_id,
			version,
			commit
		};

		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${article_id}/reset`, data)
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	copyArticle(article_id, language_id) {
		if(!language_id) {
			language_id = 0;
		}

		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${article_id}/copy`, {language_id})
			.then(r => {
				resolve(r.data.article_id);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getPrivateArticles(offset, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/private`, {params: {offset}, cancelToken})
			.then(r => {
				resolve(r.data || []);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getDrafts(offset, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/draft`, {params: {offset}, cancelToken})
			.then(r => {
				resolve(r.data || []);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	addToList(article_id, article_list_ids) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${article_id}/list`, {article_list_ids})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	clap(article_id, claps) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${article_id}/clap`, {claps})
			.then(r => {
				resolve(r.data.claps_added);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	confirmRecommendations(article_id, notify_user_ids) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/${article_id}/recommendation`, {notify_user_ids})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	deleteHistoryEntry(content_id) {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/article/history`, {params: {content_id}})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}
};
