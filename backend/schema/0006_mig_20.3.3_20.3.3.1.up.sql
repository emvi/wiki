BEGIN;

ALTER TABLE "user" ALTER COLUMN picture TYPE varchar(2000);

COMMIT;