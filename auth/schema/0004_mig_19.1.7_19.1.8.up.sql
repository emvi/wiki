BEGIN;

ALTER TABLE "user" DROP CONSTRAINT user_email_unique;
ALTER TABLE "user" DROP CONSTRAINT user_new_email_unique;
ALTER TABLE "user" ADD COLUMN "auth_provider" character varying(10) NOT NULL DEFAULT 'emvi';
ALTER TABLE "user" ADD COLUMN "auth_provider_user_id" character varying(100);
ALTER TABLE "user" ADD COLUMN "picture_url" character varying(500);

COMMIT;
