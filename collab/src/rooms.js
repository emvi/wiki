let nextRoomID = 1;
let nextRoomIDMax = Number.MAX_SAFE_INTEGER;

function nextRoom() {
	if(nextRoomID === nextRoomIDMax) {
		nextRoomID = 1;
	}

	return nextRoomID++;
}

function getRoomId(organization, id) {
	if(id) {
		return id;
	}

	return organization+nextRoom();
}

// rooms is a map of Article objects with the room ID as key.
function findRoomIdByArticleIdLangId(rooms, article_id, lang_id) {
	for(let [roomID, article] of rooms) {
		if(article.getArticleID() === article_id && (!lang_id || article.getLanguageID() === lang_id)) {
			return roomID;
		}
	}

	return null;
}

module.exports = {nextRoom, getRoomId, findRoomIdByArticleId: findRoomIdByArticleIdLangId};
