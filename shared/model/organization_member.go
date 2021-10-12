package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"time"
)

type OrganizationMember struct {
	db.BaseEntity

	OrganizationId            hide.ID     `db:"organization_id" json:"organization_id"`
	UserId                    hide.ID     `db:"user_id" json:"user_id"`
	LanguageId                hide.ID     `db:"language_id" json:"language_id"`
	Username                  string      `json:"username"`
	Phone                     null.String `json:"phone"`
	Mobile                    null.String `json:"mobile"`
	Info                      null.String `json:"info"`
	IsModerator               bool        `db:"is_moderator" json:"is_moderator"`
	IsAdmin                   bool        `db:"is_admin" json:"is_admin"`
	ReadOnly                  bool        `db:"read_only" json:"read_only"`
	Active                    bool        `json:"-"`
	LastSeen                  time.Time   `db:"last_seen" json:"-"`
	SendNotificationsInterval uint        `db:"send_notifications_interval" json:"send_notifications_interval"`
	DesktopNotifications      bool        `db:"desktop_notifications" json:"desktop_notifications"`
	NextNotificationMail      time.Time   `db:"next_notification_mail" json:"-"`
	RecommendationMail        bool        `db:"recommendation_mail" json:"recommendation_mail"`
	ShowCreateButton          bool        `db:"show_create_button" json:"show_create_button"`
	ShowNavigation            bool        `db:"show_navigation" json:"show_navigation"`
	ShowActionButtons         bool        `db:"show_action_buttons" json:"show_action_buttons"`

	User *User `db:"user" json:"user"`
}

func GetOrganizationMemberByUsername(name string) *OrganizationMember {
	entity := new(OrganizationMember)

	if err := connection.Get(entity, `SELECT * FROM "organization_member"
		WHERE LOWER(username) = LOWER($1)
		AND active IS TRUE`, name); err != nil {
		logbuch.Debug("Organization member by username not found", logbuch.Fields{"err": err, "name": name})
		return nil
	}

	return entity
}

func GetOrganizationMemberByOrganizationIdAndUsername(orgaId hide.ID, name string) *OrganizationMember {
	entity := new(OrganizationMember)

	if err := connection.Get(entity, `SELECT * FROM "organization_member"
		WHERE organization_id = $1
		AND LOWER(username) = LOWER($2)
		AND active IS TRUE`, orgaId, name); err != nil {
		logbuch.Debug("Organization member by organization id and username not found", logbuch.Fields{"err": err, "orga_id": orgaId, "name": name})
		return nil
	}

	return entity
}

func GetOrganizationMemberByOrganizationIdAndUserIdAndIsAdmin(orgaId, userId hide.ID) *OrganizationMember {
	entity := new(OrganizationMember)

	if err := connection.Get(entity, `SELECT * FROM "organization_member"
		WHERE organization_id = $1
		AND user_id = $2
		AND is_admin IS TRUE
		AND active IS TRUE`, orgaId, userId); err != nil {
		logbuch.Debug("Organization member by organization id and user id and is admin not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return nil
	}

	return entity
}

func GetOrganizationMemberByOrganizationIdAndUserIdAndIsMod(orgaId, userId hide.ID) *OrganizationMember {
	entity := new(OrganizationMember)

	if err := connection.Get(entity, `SELECT * FROM "organization_member"
		WHERE organization_id = $1
		AND user_id = $2
		AND is_moderator IS TRUE
		AND active IS TRUE`, orgaId, userId); err != nil {
		logbuch.Debug("Organization member by organization id and user id and is moderator not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return nil
	}

	return entity
}

func GetOrganizationMemberByOrganizationIdAndUserId(orgaId, userId hide.ID) *OrganizationMember {
	return GetOrganizationMemberByOrganizationIdAndUserIdTx(nil, orgaId, userId)
}

func GetOrganizationMemberByOrganizationIdAndUserIdTx(tx *sqlx.Tx, orgaId, userId hide.ID) *OrganizationMember {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(OrganizationMember)

	if err := tx.Get(entity, `SELECT "organization_member".*,
		"user".id "user.id",
		"user".email "user.email",
		"user".firstname "user.firstname",
		"user".lastname "user.lastname",
		"user".language "user.language",
		"user".info "user.info",
		"user".picture "user.picture",
		"user".accept_marketing "user.accept_marketing"
		FROM "organization_member"
		JOIN "user" ON "organization_member".user_id = "user".id
		WHERE organization_id = $1
		AND user_id = $2
		AND active IS TRUE`, orgaId, userId); err != nil {
		logbuch.Debug("Organization member by organization id and user id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return nil
	}

	return entity
}

func GetOrganizationMemberByOrganizationIdAndUserIdAndLastSeenBeforeToday(orgaId, userId hide.ID) *OrganizationMember {
	entity := new(OrganizationMember)

	if err := connection.Get(entity, `SELECT * FROM "organization_member"
		WHERE organization_id = $1
		AND user_id = $2
		AND last_seen < CURRENT_DATE`, orgaId, userId); err != nil {
		return nil
	}

	return entity
}

func GetOrganizationMemberByOrganizationIdAndUserIdIgnoreActiveTx(tx *sqlx.Tx, orgaId, userId hide.ID) *OrganizationMember {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(OrganizationMember)

	if err := tx.Get(entity, `SELECT * FROM "organization_member"
		WHERE organization_id = $1
		AND user_id = $2`, orgaId, userId); err != nil {
		logbuch.Debug("Organization member by organization id and user id ignoring activity not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return nil
	}

	return entity
}

func FindOrganizationMemberByOrganizationId(orgaId hide.ID) []OrganizationMember {
	query := `SELECT "organization_member".*,
		"user".id "user.id",
		"user".email "user.email",
		"user".firstname "user.firstname",
		"user".lastname "user.lastname",
		"user".info "user.info",
		"user".picture "user.picture"
		FROM "organization_member"
		JOIN "user" ON "organization_member".user_id = "user".id
		WHERE organization_id = $1
		AND active IS TRUE
		ORDER BY "user".lastname, "user".firstname, username ASC`
	var entities []OrganizationMember

	if err := connection.Select(&entities, query, orgaId); err != nil {
		logbuch.Error("Error reading member by organization id", logbuch.Fields{"err": err, "orga_id": orgaId})
		return nil
	}

	return entities
}

func FindOrganizationMemberWithNextNotificationMailReachedCursor() (*sqlx.Rows, error) {
	rows, err := connection.Queryx(`SELECT "organization_member".*,
		"user".id "user.id",
		"user".email "user.email",
		"user".firstname "user.firstname",
		"user".lastname "user.lastname",
		"user"."language" "user.language",
		"user".info "user.info",
		"user".picture "user.picture",
		"user".accept_marketing "user.accept_marketing"
		FROM "organization_member"
		JOIN "user" ON "organization_member".user_id = "user".id
		WHERE active IS TRUE
		AND send_notifications_interval > 0
		AND next_notification_mail < NOW()`)

	if err != nil {
		logbuch.Error("Error reading user with next notification interval reached", logbuch.Fields{"err": err})
		return nil, err
	}

	return rows, nil
}

func CountOrganizationMemberWithNextNotificationMailReached() int {
	query := `SELECT COUNT(1)
		FROM "organization_member"
		JOIN "user" ON "organization_member".user_id = "user".id
		WHERE active IS TRUE
		AND send_notifications_interval > 0
		AND next_notification_mail < NOW()`
	var count int

	if err := connection.Get(&count, query); err != nil {
		logbuch.Error("Error counting user with next notification interval reached", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func CountOrganizationMemberByOrganizationIdAndActive(orgaId hide.ID) int {
	query := `SELECT COUNT(1)
		FROM "organization_member"
		WHERE active IS TRUE
		AND organization_id = $1`
	var count int

	if err := connection.Get(&count, query, orgaId); err != nil {
		logbuch.Error("Error counting organization member by organization id and active", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func CountOrganizationMemberByOrganizationIdAndActiveAndNotReadOnly(orgaId hide.ID) int {
	query := `SELECT COUNT(1)
		FROM "organization_member"
		WHERE active IS TRUE
		AND read_only IS FALSE
		AND organization_id = $1`
	var count int

	if err := connection.Get(&count, query, orgaId); err != nil {
		logbuch.Error("Error counting organization member by organization id and active and not read only", logbuch.Fields{"err": err, "orga_id": orgaId})
		return 0
	}

	return count
}

func CountOrganizationMemberByOrganizationIdAndLastSeenAfter(orgaId hide.ID, lastSeen time.Time) int {
	query := `SELECT COUNT(1)
		FROM "organization_member"
		WHERE active IS TRUE
		AND read_only IS FALSE
		AND organization_id = $1
		AND last_seen >= $2`
	var count int

	if err := connection.Get(&count, query, orgaId, lastSeen); err != nil {
		logbuch.Error("Error counting organization member by organization id and active and not read only and last seen after", logbuch.Fields{"err": err, "orga_id": orgaId, "last_seen": lastSeen})
		return 0
	}

	return count
}

func UpdateOrganizationMemberLastSeenById(id hide.ID) error {
	query := `UPDATE "organization_member" SET last_seen = CURRENT_DATE WHERE id = $1`

	if _, err := connection.DB.Exec(query, id); err != nil {
		logbuch.Error("Error updating article list entries by article list id", logbuch.Fields{"err": err, "id": id})
		return err
	}

	return nil
}

func SaveOrganizationMember(tx *sqlx.Tx, entity *OrganizationMember) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "organization_member" (organization_id,
			user_id,
			language_id,
			username,
			phone,
			mobile,
			info,
			is_moderator,
			is_admin,
			read_only,
			active,
			last_seen,
			send_notifications_interval,
			desktop_notifications,
			next_notification_mail,
			recommendation_mail,
			show_create_button,
			show_navigation,
			show_action_buttons)
			VALUES (:organization_id,
			:user_id,
			:language_id,
			:username,
			:phone,
			:mobile,
			:info,
			:is_moderator,
			:is_admin,
			:read_only,
			:active,
			:last_seen,
			:send_notifications_interval,
			:desktop_notifications,
			:next_notification_mail,
			:recommendation_mail,
			:show_create_button,
			:show_navigation,
			:show_action_buttons) RETURNING id`,
		`UPDATE "organization_member" SET organization_id = :organization_id,
			user_id = :user_id,
			language_id = :language_id,
			username = :username,
			phone = :phone,
			mobile = :mobile,
			info = :info,
			is_moderator = :is_moderator,
			is_admin = :is_admin,
			read_only = :read_only,
			active = :active,
			last_seen = :last_seen,
			send_notifications_interval = :send_notifications_interval,
			desktop_notifications = :desktop_notifications,
			next_notification_mail = :next_notification_mail,
			recommendation_mail = :recommendation_mail,
			show_create_button = :show_create_button,
			show_navigation = :show_navigation,
			show_action_buttons = :show_action_buttons
			WHERE id = :id`)
}
