import axios from "axios";

export const BillingService = new class {
	subscribe(email, name, country, address_line_1, address_line_2, postal_code, city, phone, tax_number, interval, payment_method_id) {
		let data = {
			email,
			name,
			country,
			address_line_1,
			address_line_2,
			postal_code,
			city,
			phone,
			tax_number,
			interval,
			payment_method_id
		};

		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/subscription`, data)
				.then(r => {
					resolve(r.data.client_secret);
				})
				.catch(e => {
					reject(e);
				});
		});
	}

	cancelSubscription() {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/subscription`)
			.then(r => {
				resolve(r);
			})
			.catch(e => {
				reject(e);
			});
		});
	}

	resumeSubscription() {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/subscription`)
				.then(r => {
					resolve(r);
				})
				.catch(e => {
					reject(e);
				});
		});
	}

	getCustomer() {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/subscription`)
				.then(r => {
					resolve(r.data);
				})
				.catch(e => {
					reject(e);
				});
		});
	}

	getInvoices(start_invoice_id) {
		return new Promise((resolve, reject) => {
			axios.get(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/subscription/invoice`, {params: {start_invoice_id: start_invoice_id || ""}})
				.then(r => {
					resolve(r.data || []);
				})
				.catch(e => {
					reject(e);
				});
		});
	}

	updateCustomer(email, name, country, address_line_1, address_line_2, postal_code, city, phone, tax_number) {
		let data = {
			email,
			name,
			country,
			address_line_1,
			address_line_2,
			postal_code,
			city,
			phone,
			tax_number
		};

		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/subscription/customer`, data)
				.then(r => {
					resolve(r.data);
				})
				.catch(e => {
					reject(e);
				});
		});
	}

	updatePaymentMethod(payment_method_id) {
		return new Promise((resolve, reject) => {
			axios.post(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/subscription/payment`, {payment_method_id})
				.then(r => {
					resolve(r.data);
				})
				.catch(e => {
					reject(e);
				});
		});
	}

	changePlan(interval) {
		return new Promise((resolve, reject) => {
			axios.put(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/subscription/plan`, {interval})
				.then(r => {
					resolve(r.data);
				})
				.catch(e => {
					reject(e);
				});
		});
	}

	removePaymentIntentClientSecret() {
		return new Promise((resolve, reject) => {
			axios.delete(`${EMVI_WIKI_BACKEND_HOST}/api/v1/organization/subscription/payment`)
				.then(r => {
					resolve(r.data);
				})
				.catch(e => {
					reject(e);
				});
		});
	}
};
