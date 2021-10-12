BEGIN;

ALTER TABLE "user" ADD COLUMN "introduction" boolean DEFAULT TRUE;

COMMIT;
