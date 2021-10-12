BEGIN;

ALTER TABLE "feed_ref"
    ADD COLUMN "key" varchar(100),
    ADD COLUMN "value" text;

INSERT INTO feed_ref (feed_id, "key", "value")
	SELECT id "feed_id", 'name' "key", deleted_name "value" FROM feed
	WHERE reason IN ('delete_article', 'delete_articlelist', 'delete_tag', 'delete_usergroup');

ALTER TABLE feed DROP COLUMN deleted_name;

COMMIT;
