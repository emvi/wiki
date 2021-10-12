const express = require("express");
const cors = require('cors');
const http = require("http");
const socketIo = require("socket.io");

const {logger} = require("./logger.js");
const {acceptConnection, handleConnection, saveArticles, getRoomCount, getConnectionCount} = require("./connection.js");
const {apiClient} = require("./client.js");

const PORT = process.env.EMVI_WIKI_PORT || 4004;
const SHUTDOWN_TIMEOUT = parseInt(process.env.EMVI_WIKI_SHUTDOWN_TIMEOUT) || 30; // seconds

process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);

// build and start server
const app = express();
const server = http.createServer(app);
const io = socketIo(server, {
	path: "/api/v1/collab/ws",
	origins: "*:*",
	cookie: false
});

app.get("/api/v1/collab/health", (req, res) => {
	res.send("Alive!");
});

app.get("/api/v1/collab/stats", (req, res) => {
	res.send({
		rooms: getRoomCount(),
		connections: getConnectionCount()
	});
});

app.use(cors());

server.listen(PORT, () => {
	logger.info("Server started");
});

// handle shutdown
function shutdown() {
	logger.info("Server shutting down...");

	io.close(() => {
		saveArticles(io);
		logger.info("Shutdown completed");
		process.exit(0);
	});

	setTimeout(() => {
		logger.info("Forcefully shutting down");
		process.exit(1);
	}, SHUTDOWN_TIMEOUT*1000);
}

// accept websocket connections
io.on("connection", client => {
	acceptConnection(client)
	.then(r => {
		handleConnection(io,
			apiClient,
			client,
			r.organization,
			r.data.user.id,
			r.read_only);
	})
	.catch(e => {
		if(e.hasOwnProperty("Error") && e.Error) {
			logger.error(e.Error);
		}

		logger.debug("Connection refused", {error: e.message});
		client.disconnect();
	});
});
