package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
)

type Organization struct {
	db.BaseEntity

	Name                            string      `json:"name"`
	NameNormalized                  string      `db:"name_normalized" json:"name_normalized"`
	Picture                         null.String `json:"picture"`
	Expert                          bool        `json:"expert"`
	MaxStorageGB                    int64       `db:"max_storage_gb" json:"-"`
	CreateGroupAdmin                bool        `db:"create_group_admin" json:"create_group_admin"`
	CreateGroupMod                  bool        `db:"create_group_mod" json:"create_group_mod"`
	InvitationCode                  null.String `db:"invitation_code" json:"invitation_code"`
	InvitationReadOnly              bool        `db:"invitation_read_only" json:"invitation_read_only"`
	StripeCustomerID                null.String `db:"stripe_customer_id" json:"-"`
	StripeSubscriptionID            null.String `db:"stripe_subscription_id" json:"-"`
	StripePaymentMethodID           null.String `db:"stripe_payment_method_id" json:"stripe_payment_method_id"`
	StripePaymentIntentClientSecret null.String `db:"stripe_payment_intent_client_secret" json:"stripe_payment_intent_client_secret"`
	SubscriptionPlan                null.String `db:"subscription_plan" json:"subscription_plan"`
	SubscriptionCancelled           bool        `db:"subscription_cancelled" json:"subscription_cancelled"`
	SubscriptionCycle               null.Time   `db:"subscription_cycle" json:"-"`

	OwnerUserId  hide.ID `db:"owner_user_id" json:"-"`
	IsAdmin      bool    `db:"is_admin" json:"is_admin"` // used to let client know if user is admin for this organization
	MemberCount  int     `db:"member_count" json:"member_count"`
	ArticleCount int     `db:"article_count" json:"article_count"`
}

func GetOrganizationById(id hide.ID) *Organization {
	return GetOrganizationByIdTx(nil, id)
}

func GetOrganizationByIdTx(tx *sqlx.Tx, id hide.ID) *Organization {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(Organization)

	if err := tx.Get(entity, `SELECT * FROM "organization" WHERE id = $1`, id); err != nil {
		logbuch.Debug("Organization by id not found", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return entity
}

func GetOrganizationByName(name string) *Organization {
	entity := new(Organization)

	if err := connection.Get(entity, `SELECT * FROM "organization" WHERE LOWER(name) = LOWER($1)`, name); err != nil {
		logbuch.Debug("Organization by name not found", logbuch.Fields{"err": err, "name": name})
		return nil
	}

	return entity
}

func GetOrganizationByNameNormalized(name string) *Organization {
	entity := new(Organization)

	if err := connection.Get(entity, `SELECT * FROM "organization" WHERE LOWER(name_normalized) = LOWER($1)`, name); err != nil {
		logbuch.Debug("Organization by name normalized not found", logbuch.Fields{"err": err, "name": name})
		return nil
	}

	return entity
}

func GetOrganizationByInvitationCode(code string) *Organization {
	entity := new(Organization)

	if err := connection.Get(entity, `SELECT * FROM "organization" WHERE invitation_code = $1`, code); err != nil {
		logbuch.Debug("Organization by invitation code not found", logbuch.Fields{"err": err, "code": code})
		return nil
	}

	return entity
}

func GetOrganizationByStripeSubscriptionID(subId string) *Organization {
	entity := new(Organization)

	if err := connection.Get(entity, `SELECT * FROM "organization" WHERE stripe_subscription_id = $1`, subId); err != nil {
		logbuch.Debug("Organization by stripe subscription id not found", logbuch.Fields{"err": err, "subscription_id": subId})
		return nil
	}

	return entity
}

func GetOrganizationByStripeCustomerID(customerId string) *Organization {
	entity := new(Organization)

	if err := connection.Get(entity, `SELECT * FROM "organization" WHERE stripe_customer_id = $1`, customerId); err != nil {
		logbuch.Debug("Organization by stripe customer id not found", logbuch.Fields{"err": err, "customer_id": customerId})
		return nil
	}

	return entity
}

func GetOrganizationByUserIdAndNameNormalized(userId hide.ID, name string) *Organization {
	entity := new(Organization)

	if err := connection.Get(entity, `SELECT "organization".* FROM "organization"
		JOIN "organization_member" ON "organization".id = "organization_member".organization_id
		WHERE "organization_member".user_id = $1
		AND name_normalized = $2
		AND "organization_member".active IS TRUE`, userId, name); err != nil {
		logbuch.Debug("Organization by user id and name not found", logbuch.Fields{"err": err, "name": name})
		return nil
	}

	return entity
}

func GetOrganizationByUserIdAndIdAndIsAdmin(orgaId, userId hide.ID) *Organization {
	entity := new(Organization)

	if err := connection.Get(entity, `SELECT "organization".* FROM "organization"
		JOIN "organization_member" ON "organization".id = "organization_member".organization_id
		WHERE "organization_member".user_id = $2
		AND "organization_member".is_admin IS TRUE
		AND "organization".id = $1`, orgaId, userId); err != nil {
		logbuch.Debug("Organization by user id and id and is admin not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return nil
	}

	return entity
}

func FindOrganizationByUserIdAndNotExpert(userId hide.ID) []Organization {
	var entities []Organization

	if err := connection.Select(&entities, `SELECT * FROM "organization" WHERE owner_user_id = $1 AND expert IS FALSE`, userId); err != nil {
		logbuch.Error("Error reading organizations by user id and not expert", logbuch.Fields{"err": err, "userId": userId})
		return nil
	}

	return entities
}

func FindOrganizationsByUserId(userId hide.ID) []Organization {
	var entities []Organization

	if err := connection.Select(&entities, `SELECT "organization".*, "organization_member".is_admin,
		(SELECT COUNT(1) FROM "organization_member" WHERE organization_id = "organization".id AND "organization_member".active IS TRUE) AS "member_count",
		(SELECT COUNT(1) FROM "article" WHERE organization_id = "organization".id AND "article".archived IS NULL) AS "article_count"
		FROM "organization"
		JOIN "organization_member" ON "organization".id = "organization_member".organization_id AND "organization_member".active IS TRUE
		WHERE "organization_member".user_id = $1
		ORDER BY "organization".name ASC`, userId); err != nil {
		logbuch.Error("Error reading organizations by user id", logbuch.Fields{"err": err, "userId": userId})
		return nil
	}

	return entities
}

func FindOrganizationWithSubscriptionCycleReached() (*sqlx.Rows, error) {
	query := `SELECT * FROM "organization"
		WHERE "expert" IS TRUE
		AND "subscription_cycle" + INTERVAL '1 month' <= CURRENT_DATE`
	rows, err := connection.Queryx(query)

	if err != nil {
		logbuch.Error("Error reading organizations with subscription cycle today and expert true", logbuch.Fields{"err": err})
		return nil, err
	}

	return rows, nil
}

func DeleteOrganizationById(tx *sqlx.Tx, orgaId hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := tx.Exec(`DELETE FROM "support_ticket" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting support tickets when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "bookmark" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting bookmarks when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "file" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting tags when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "invitation" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting invitations when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM observed_object
		WHERE user_group_id IN (SELECT id FROM user_group WHERE organization_id = $1)
		OR article_id IN (SELECT id FROM article WHERE organization_id = $1)
		OR article_list_id IN (SELECT id FROM article_list WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting observed objects when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "feed_access"
		WHERE feed_id IN (SELECT id FROM feed WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting feed access when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "feed_ref" WHERE feed_id IN (SELECT id FROM "feed" WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting feed ref when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "feed" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting feed when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article_list_entry"
		WHERE article_list_id IN (SELECT id FROM article_list WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting article list entries when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article_list_member"
		WHERE article_list_id IN (SELECT id FROM article_list WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting article list member when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article_list_name"
		WHERE article_list_id IN (SELECT id FROM article_list WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting article list names when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article_list" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting article lists when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article_tag" WHERE tag_id IN (SELECT id FROM "tag" WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting tags when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "tag" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting tags when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article_content_author"
		WHERE article_content_id IN (SELECT article_content.id FROM article_content JOIN article ON article_content.article_id = article.id WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting article content authors when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article_access"
		WHERE article_id IN (SELECT id FROM article WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting article access when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article_visit"
		WHERE article_id IN (SELECT id FROM article WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting article visit when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article_recommendation"
		WHERE article_id IN (SELECT id FROM article WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting article recommendation when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article_content"
		WHERE article_id IN (SELECT id FROM article WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting article content when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "article" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting articles when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "user_group_member"
		WHERE user_group_id IN (SELECT id FROM user_group WHERE organization_id = $1)`, orgaId); err != nil {
		logbuch.Error("Error deleting user group members when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "user_group" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting user groups when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "organization_member" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting organization members when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "language" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting languages when deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	if _, err := tx.Exec(`DELETE FROM "organization" WHERE id = $1`, orgaId); err != nil {
		logbuch.Error("Error deleting organization by id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveOrganization(tx *sqlx.Tx, entity *Organization) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "organization" (name,
			name_normalized,
			picture,
			expert,
			max_storage_gb,
			create_group_admin,
			create_group_mod,
			owner_user_id,
			invitation_code,
			invitation_read_only,
			stripe_customer_id,
			stripe_subscription_id,
			stripe_payment_method_id,
			stripe_payment_intent_client_secret,
			subscription_plan,
			subscription_cancelled,
			subscription_cycle)
			VALUES (:name,
			:name_normalized,
			:picture,
			:expert,
			:max_storage_gb,
			:create_group_admin,
			:create_group_mod,
			:owner_user_id,
			:invitation_code,
			:invitation_read_only,
			:stripe_customer_id,
			:stripe_subscription_id,
			:stripe_payment_method_id,
			:stripe_payment_intent_client_secret,
			:subscription_plan,
			:subscription_cancelled,
			:subscription_cycle) RETURNING id`,
		`UPDATE "organization" SET name = :name,
			name_normalized = :name_normalized,
			picture = :picture,
			expert = :expert,
			max_storage_gb = :max_storage_gb,
			create_group_admin = :create_group_admin,
			create_group_mod = :create_group_mod,
			owner_user_id = :owner_user_id,
			invitation_code = :invitation_code,
			invitation_read_only = :invitation_read_only,
			stripe_customer_id = :stripe_customer_id,
			stripe_subscription_id = :stripe_subscription_id,
			stripe_payment_method_id = :stripe_payment_method_id,
			stripe_payment_intent_client_secret = :stripe_payment_intent_client_secret,
			subscription_plan = :subscription_plan,
			subscription_cancelled = :subscription_cancelled,
			subscription_cycle = :subscription_cycle
			WHERE id = :id`)
}
