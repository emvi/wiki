BEGIN;

ALTER TABLE "article_content" ADD COLUMN "schema_version" integer NOT NULL DEFAULT 1;

COMMIT;
