const {schema} = require("./schema.js");
const {logger} = require("./logger.js");

const MAX_STEP_HISTORY = 100;
const MAX_NON_EXPERT_CLIENTS = 2;

class Article {
	constructor(apiClient) {
		this.apiClient = apiClient;

		// article data
		this.id = "";
		this.title = "";
		this.doc = schema.node("doc", null,
			[schema.node("paragraph", null)]);
		this.version = 0;
		this.tags = [];
		this.lang = ""; // language ID
		this.accessMode = 0;
		this.access = []; // user and groups with write access true/false
		this.clientAccess = false;
		this.authors = []; // list of author IDs that contribute to this change
		this.rtl = false; // right to left text direction

		// technical data
		this.steps = [];
		this.clientIDs = [];
		this.saved = true; // don't save if not modified
		this.initialized = false;
		this.maxClients = 0; // 0 = unlimited, depends on entry/expert version of organization
	}

	load(client, id, lang_id) {
		return new Promise((resolve, reject) => {
			this.apiClient.getArticle(client, id, lang_id)
			.then(r => {
				this.id = r.data.article.id;
				this.title = r.data.content.title;
				this.tags = Article.tagsToStringArray(r.data.article.tags);
				this.lang = r.data.content.language_id;

				if(r.data.content.content) {
					this.doc = schema.nodeFromJSON(JSON.parse(r.data.content.content));
					r.doc = r.data.content.content;
				}
				else {
					r.doc = JSON.stringify(this.doc.toJSON());
				}

				if(r.data.article.write_everyone) {
					this.accessMode = 0;
				}
				else if(r.data.article.read_everyone) {
					this.accessMode = 1;
				}
				else if(!r.data.article.private) {
					this.accessMode = 2;
				}
				else {
					this.accessMode = 3;
				}

				this.access = r.data.article.access;
				this.clientAccess = r.data.article.client_access;
				this.rtl = r.data.content.rtl;
				logger.debug("Article loaded", {article_id: r.data.article.id});
				resolve(r);
			})
			.catch(e => {
				logger.error("Error loading article", {error: e.message});
				reject(e);
			});
		});
	}

	init(client) {
		return new Promise((resolve, reject) => {
			if(this.initialized) {
				resolve();
			}

			this.apiClient.getOrganization(client)
			.then(r => {
				if(!r.data.expert) {
					this.maxClients = MAX_NON_EXPERT_CLIENTS;
				}

				this.initialized = true;
				resolve();
			})
			.catch(e => {
				logger.error("Error initializing article", {error: e.message});
				reject(e);
			});
		})
	}

	// Saves the article. If explicit is set, it will be saved even if it has not been touched.
	save(client, data, explicit) {
		// only save if something changed
		if(this.saved && (!explicit || data.wip)) {
			return new Promise((resolve) => {resolve();});
		}

		// if data is undefined this is an automatic save
		if(!data) {
			data = this.saveSetAutomaticSave();
		}

		let req = {
			user_id: client.user.id,
			id: this.id,
			room_id: client.roomID,
			authors: this.authors,
			message: data.message,
			wip: data.wip,
			read_everyone: this.accessMode === 0 || this.accessMode === 1,
			write_everyone: this.accessMode === 0,
			private: this.accessMode === 3,
			client_access: this.clientAccess,
			access: this.access,
			title: this.title.substr(0, 100),
			content: JSON.stringify(this.doc.toJSON()),
			rtl: this.rtl,
			tags: this.tags,
			language_id: this.lang
		};
		this.saved = true;
		logger.debug("Saving article");

		return new Promise((resolve, reject) => {
			this.apiClient.saveArticle(client, req)
			.then(r => {
				logger.info("Article saved", {id: r.data.id});
				this.id = r.data.id;
				resolve(r);
			})
			.catch(e => {
				logger.error("Error saving article", {error: e.message});
				reject(e.response ? e.response.data : e);
			});
		});
	}

	saveSetAutomaticSave() {
		// default title if not set
		if(!this.title){
			this.title = "New untitled article"
		}

		// if this is a new article, set to private
		if(!this.id) {
			this.accessMode = 3;
		}

		return {message: "Work in progress", wip: true};
	}

	update(version, steps, clientID) {
		// if client is behind, refuse changes
		if(this.version !== version) {
			return false;
		}

		// try to apply the steps to a copy of the document first
		let doc = this.doc;
		let newSteps = [];
		let newClientIDs = [];

		steps.forEach(step => {
			let result = step.apply(doc);

			if(result.failed) {
				logger.debug(result.failed);
				return false;
			}

			doc = result.doc;
			newSteps.push(step);
			newClientIDs.push(clientID);
		});

		// add steps and update document to new version
		this.saved = false; // document got touched
		this.version += steps.length;
		this.doc = doc;
		this.steps = this.steps.concat(newSteps);
		this.clientIDs = this.clientIDs.concat(newClientIDs);

		if(this.steps.length > MAX_STEP_HISTORY) {
			this.steps = this.steps.slice(this.steps.length-MAX_STEP_HISTORY);
			this.clientIDs = this.clientIDs.slice(this.clientIDs.length-MAX_STEP_HISTORY);
		}

		return true;
	}

	static tagsToStringArray(tags) {
		let tagStr = [];

		if(tags) {
			for (let i = 0; i < tags.length; i++) {
				tagStr.push(tags[i].name);
			}
		}

		return tagStr;
	}

	addTag(client, tag) {
		return new Promise((resolve, reject) => {
			if(this.tagExists(tag) !== -1) {
				resolve(false);
				return;
			}

			this.apiClient.addTag(client, tag)
			.then(() => {
				this.tags.push(tag);
				resolve(true);
			})
			.catch(e => {
				logger.debug("Error adding tag to article", {error: e.message});
				reject(e);
			});
		});
	}

	removeTag(tag) {
		let index = this.tagExists(tag);

		if(index === -1) {
			return false;
		}

		this.tags.splice(index, 1);
		return true;
	}

	getTags() {
		return this.tags;
	}

	tagExists(tag) {
		let tagLowerCase = tag.toLowerCase();
		
		for(let i = 0; i < this.tags.length; i++){
			if(this.tags[i].toLowerCase() === tagLowerCase) {
				return i;
			}
		}

		return -1;
	}

	getState(since = 0) {
		return {steps: this.steps.slice(since),
			clientIDs: this.clientIDs.slice(since)};
	}

	getArticleID() {
		return this.id;
	}

	getDoc() {
		return this.doc;
	}

	getVersion() {
		return this.version;
	}

	setTitle(title) {
		this.saved = false;
		this.title = title;
	}

	getTitle() {
		return this.title;
	}

	setSaved(saved) {
		this.saved = saved;
	}

	setLanguage(lang) {
		this.lang = lang;
	}

	getLanguageID() {
		return this.lang;
	}

	setAccessMode(mode) {
		this.saved = false;
		this.accessMode = parseInt(mode);
	}

	setClientAccess(access) {
		this.saved = false;
		this.clientAccess = access;
	}

	getAccessMode() {
		return this.accessMode;
	}

	getClientAccess() {
		return this.clientAccess;
	}

	setAccess(access) {
		this.saved = false;
		access = Article.setAccessReferenceIds(access);
		let index = this.accessExists(access);

		if(index === -1) {
			this.access.push(access);
		}
		else {
			this.access[index].write = access.write;
		}
	}

	getAccess() {
		return this.access;
	}

	removeAccess(access) {
		this.saved = false;
		access = Article.setAccessReferenceIds(access);
		let index = this.accessExists(access);

		if(index === -1) {
			return false;
		}

		this.access.splice(index, 1);
		return true;
	}

	accessExists(access) {
		access.user_id = access.user_id || null;
		access.user_group_id = access.user_group_id || null;

		for(let i = 0; i < this.access.length; i++){
			if(access.user_id !== null && this.access[i].user_id === access.user_id ||
				access.user_group_id !== null && this.access[i].user_group_id === access.user_group_id) {
				return i;
			}
		}

		return -1;
	}

	static setAccessReferenceIds(access) {
		if(access.user) {
			access.user_id = access.user.id;
		}
		else {
			access.user_group_id = access.group.id;
		}

		return access;
	}

	addAuthor(user_id) {
		for(let i = 0; i < this.authors.length; i++){
			if(this.authors[i] === user_id) {
				return false;
			}
		}

		this.authors.push(user_id);
		return true;
	}

	removeAuthor(user_id) {
		for(let i = 0; i < this.authors.length; i++){
			if(this.authors[i] === user_id) {
				this.authors.splice(i, 1);
				break;
			}
		}
	}

	hasAuthors() {
		return this.authors.length > 0;
	}

	getAuthors() {
		return this.authors;
	}

	maxClientsReached() {
		return this.maxClients !== 0 && this.authors.length > this.maxClients;
	}

	setRTL(rtl) {
		this.saved = false;
		this.rtl = rtl;
	}

	getRTL() {
		return this.rtl;
	}
}

module.exports.Article = Article;
