BEGIN;

ALTER TABLE "invitation" DROP COLUMN "user_id";

COMMIT;
