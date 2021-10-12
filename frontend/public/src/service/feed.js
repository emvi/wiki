import axios from "axios";

export const FeedService = new class {
	getFeed(filter, cancelToken) {
		if(cancelToken) {
			cancelToken = cancelToken.token;
		}

		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/feed`, {params: filter, cancelToken})
			.then(r => {
				resolve(r.data.feed || []);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	getNotifications(unread, limit) {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/feed`, {params: {notifications: true, unread, offset: 0, limit}})
			.then(r => {
				resolve({feed: r.data.feed || [], count: r.data.count});
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	toggleNotificationRead(id) {
		return new Promise((resolve, reject) => {
			if(!id){
				id = "";
			}

			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/feed`, {id})
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}
};
