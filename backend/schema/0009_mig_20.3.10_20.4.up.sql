BEGIN;

ALTER TABLE "article_content" ADD COLUMN "rtl" boolean DEFAULT FALSE;

COMMIT;
