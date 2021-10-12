BEGIN;

ALTER TABLE "user" ALTER COLUMN picture_url TYPE varchar(2000);

COMMIT;
