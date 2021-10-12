const LOG_LEVEL = process.env.EMVI_WIKI_LOG_LEVEL || "debug";
const LOG_DISABLE_TIMESTAMP = process.env.EMVI_WIKI_LOG_DISABLE_TIMESTAMP || "false";
console.log(`[INFO] Log level: ${LOG_LEVEL}`);
console.log(`[INFO] Log disable timestamp: ${LOG_DISABLE_TIMESTAMP}`);

const logLevels = {
	debug: 0,
	info: 1,
	warn: 2,
	error: 3
};

class Logger {
	constructor(minLevel, disableTimestamp) {
		if(minLevel === null || minLevel === undefined){
			minLevel = logLevels.info;
		}

		this.minLevel = minLevel;
		this.disableTimestamp = disableTimestamp;
	}

	debug(message, params) {
		this._log(logLevels.debug, message, params);
	}

	info(message, params) {
		this._log(logLevels.info, message, params);
	}

	warn(message, params) {
		this._log(logLevels.warn, message, params);
	}

	error(message, params) {
		this._log(logLevels.error, message, params);
	}

	_log(level, message, params) {
		if(level >= this.minLevel) {
			let levelName = this._levelToString(level);
			let paramsFlat = Logger._paramsToString(params);

			if(this.disableTimestamp) {
				console.log(`[${levelName}] ${message}\t\t\t${paramsFlat}`);
			}
			else {
				let now = Logger._now();
				console.log(`${now} [${levelName}] ${message}\t\t\t${paramsFlat}`);
			}
		}
	}

	static _now() {
		return new Date().toISOString();
	}

	_levelToString(level) {
		let str = Object.keys(logLevels).find(key => logLevels[key] === level);
		return str.substr(0, 4).toUpperCase();
	}

	static _paramsToString(params) {
		if(typeof params !== "object"){
			return "";
		}

		return JSON.stringify(params);
	}
}

module.exports.logger = new Logger(logLevels[LOG_LEVEL], LOG_DISABLE_TIMESTAMP.toLowerCase() === "true");
