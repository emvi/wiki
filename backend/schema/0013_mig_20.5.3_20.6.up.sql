BEGIN;

ALTER TABLE "organization" ADD COLUMN stripe_customer_id character varying(255);
ALTER TABLE "organization" ADD COLUMN stripe_subscription_id character varying(255);
ALTER TABLE "organization" ADD COLUMN stripe_payment_method_id character varying(255);
ALTER TABLE "organization" ADD COLUMN stripe_payment_intent_client_secret character varying(255);
ALTER TABLE "organization" ADD COLUMN "subscription_plan" character varying(7);
ALTER TABLE "organization" ADD COLUMN "subscription_cancelled" boolean DEFAULT FALSE;
ALTER TABLE "organization" ADD COLUMN "subscription_cycle" date;

ALTER TABLE "organization_member" ADD COLUMN "last_seen" date DEFAULT CURRENT_DATE;
ALTER TABLE "organization_member" ADD COLUMN "show_create_button" boolean DEFAULT TRUE;
ALTER TABLE "organization_member" ADD COLUMN "show_navigation" boolean DEFAULT TRUE;
ALTER TABLE "organization_member" ADD COLUMN "show_action_buttons" boolean DEFAULT TRUE;

COMMIT;
