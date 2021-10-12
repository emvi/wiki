let assert = require("assert");
let {nextRoom, getRoomId, findRoomIdByArticleId} = require("../src/rooms.js");
let {Article} = require("../src/article.js");

describe("nextRoom", function() {
    it("returns the next room ID", function() {
        assert.equal(nextRoom(), 1);
        assert.equal(nextRoom(), 2);
        assert.equal(nextRoom(), 3);
    });
});

describe("getRoomId", function() {
    it("returns the next room ID with organization name", function() {
        assert.ok(getRoomId("orga").includes("orga"));
    });

    it("returns the provided ID", function() {
        assert.equal(getRoomId("orga", "id"), "id");
    });
});

describe("findRoomIdByArticleId", function() {
    let article1 = new Article();
    article1.id = "id1";
    article1.lang = "lang1";
    let article2 = new Article();
    article2.id = "id2";
    article2.lang = "lang2";
    let article3 = new Article();
    article3.id = "id3";
    article3.lang = "lang3";
    let rooms = new Map([
        ["room1", article1],
        ["room2", article2],
        ["room3", article3]
    ]);

    it("finds an existing room", () => {
        assert.equal(findRoomIdByArticleId(rooms, "id2", "lang2"), "room2");
    });

    it("not finds a non existing room", () => {
        assert.equal(findRoomIdByArticleId(rooms, "id2", "unknown"), null);
        assert.equal(findRoomIdByArticleId(rooms, "unknown", "lang2"), null);
    });
});
