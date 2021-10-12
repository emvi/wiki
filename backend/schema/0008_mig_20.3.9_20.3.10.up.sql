BEGIN;

DROP TABLE article_claps;
ALTER TABLE article DROP COLUMN claps;

COMMIT;
