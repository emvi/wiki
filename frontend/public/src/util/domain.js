export function getSubdomain() {
	let host = window.location.host;
	let parts = host.split(".");
	return parts[0];
}
