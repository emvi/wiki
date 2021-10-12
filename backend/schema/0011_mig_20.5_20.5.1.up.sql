BEGIN;

ALTER TABLE "organization" ADD COLUMN "invitation_code" varchar(20);
ALTER TABLE "organization" ADD COLUMN "invitation_read_only" boolean DEFAULT FALSE;

COMMIT;
