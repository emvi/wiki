import axios from "axios";

export const OrganizationService = new class {
	getOrganizations() {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organizations`)
			.then((r) => {
				resolve(r.data || []);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}

	createOrganization(data) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization`, data)
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	validateOrganization(data, validate_name_domain, validate_username_lang) {
		data.validate_name_domain = validate_name_domain;
		data.validate_username_lang = validate_username_lang;

		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/validation`, data)
			.then(r => {
				resolve(r.data);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	inviteMember(emails, domain) {
		let config = {
			headers: {
				Organization: domain
			}
		};

		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/member`, {emails}, config)
			.then((r) => {
				resolve(r.data);
			})
			.catch((e) => {
				reject(e);
			});
		});
	}
};
