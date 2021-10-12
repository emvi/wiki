BEGIN;

ALTER TABLE "newsletter" RENAME TO "newsletter_subscription";
ALTER SEQUENCE "newsletter_id_seq" RENAME TO "newsletter_subscription_id_seq";
ALTER INDEX "newsletter_pkey" RENAME TO "newsletter_subscription_pkey";
ALTER INDEX "newsletter_email_list_unique" RENAME TO "newsletter_subscription_email_list_unique";
ALTER INDEX "newsletter_code_unique" RENAME TO "newsletter_subscription_code_unique";
ALTER TRIGGER "update_newsletter_mod_time" ON "newsletter_subscription" RENAME TO "update_newsletter_subscription_mod_time";

COMMIT;
