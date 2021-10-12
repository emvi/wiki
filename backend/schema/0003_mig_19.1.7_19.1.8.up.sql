BEGIN;

ALTER TABLE "user" DROP CONSTRAINT user_email_unique;

COMMIT;
