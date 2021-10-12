package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	userBaseQueryFields = `"user".*,
		"organization_member".id "organization_member.id",
		"organization_member".username "organization_member.username",
		"organization_member".info "organization_member.info",
		"organization_member".phone "organization_member.phone",
		"organization_member".mobile "organization_member.mobile",
		"organization_member".is_moderator "organization_member.is_moderator",
		"organization_member".is_admin "organization_member.is_admin",
		"organization_member".read_only "organization_member.read_only",
		"organization_member".recommendation_mail "organization_member.recommendation_mail" `
	userBaseQueryHead = "SELECT " + userBaseQueryFields
	userBaseQuery     = `FROM "user"
		JOIN "organization_member" ON "user".id = "organization_member".user_id AND "organization_member".organization_id = $1
		WHERE organization_id = $1
		AND "organization_member".active IS TRUE `
)

type User struct {
	db.BaseEntity

	Email           string      `json:"email"`
	Firstname       string      `json:"firstname"`
	Lastname        string      `json:"lastname"`
	Language        null.String `json:"language"`
	Info            null.String `json:"info"`
	Picture         null.String `json:"picture"`
	AcceptMarketing bool        `db:"accept_marketing" json:"accept_marketing"`
	ColorMode       int         `db:"color_mode" json:"color_mode"`
	Introduction    bool        `json:"introduction"`

	// IsSSOUser is obtains from auth and not stored in backend database
	IsSSOUser bool `json:"is_sso_user"`

	OrganizationMember *OrganizationMember `db:"organization_member" json:"organization_member"`
}

func GetUserById(id hide.ID) *User {
	return GetUserByIdTx(nil, id)
}

func GetUserByIdTx(tx *sqlx.Tx, id hide.ID) *User {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(User)

	if err := tx.Get(entity, `SELECT * FROM "user" WHERE id = $1`, id); err != nil {
		logbuch.Debug("User by id not found", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return entity
}

func GetUserByOrganizationIdAndId(orgaId, id hide.ID) *User {
	return GetUserByOrganizationIdAndIdTx(nil, orgaId, id)
}

func GetUserByOrganizationIdAndIdTx(tx *sqlx.Tx, orgaId, id hide.ID) *User {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(User)

	if err := tx.Get(entity, `SELECT "user".* FROM "user"
		JOIN "organization_member" ON "user".id = "organization_member".user_id AND "organization_member".organization_id = $2
		WHERE "user".id = $1`, id, orgaId); err != nil {
		logbuch.Debug("User by organization id and id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func GetUserWithOrganizationMemberByOrganizationIdAndId(orgaId, id hide.ID) *User {
	return GetUserWithOrganizationMemberByOrganizationIdAndIdTx(nil, orgaId, id)
}

func GetUserWithOrganizationMemberByOrganizationIdAndIdTx(tx *sqlx.Tx, orgaId, id hide.ID) *User {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := userBaseQueryHead + userBaseQuery + `AND "user".id = $2`
	entity := new(User)

	if err := tx.Get(entity, query, orgaId, id); err != nil {
		logbuch.Debug("User with organization member by organization id and id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func GetUserWithOrganizationMemberByOrganizationIdAndUsername(orgaId hide.ID, username string) *User {
	query := userBaseQueryHead + userBaseQuery + `AND LOWER("organization_member".username) = LOWER($2)`
	entity := new(User)

	if err := connection.Get(entity, query, orgaId, username); err != nil {
		logbuch.Debug("User with organization member by organization id and username not found", logbuch.Fields{"err": err, "orga_id": orgaId, "username": username})
		return nil
	}

	return entity
}

func FindUserByOrganizationIdAndUsernameOrFirstnameOrLastnameOrEmail(orgaId hide.ID, keywords string, filter *SearchUserFilter) []User {
	query, params := buildUserByOrganizationIdAndUsernameOrFirstnameOrLastnameOrEmailQuery(orgaId, keywords, filter, false)
	var entities []User

	if err := connection.Select(&entities, query, params...); err != nil {
		logbuch.Error("Error reading user by organization id and query string", logbuch.Fields{"err": err, "organization_id": orgaId, "keywords": keywords})
		return nil
	}

	return entities
}

func CountUserByOrganizationIdAndUsernameOrFirstnameOrLastnameOrEmail(orgaId hide.ID, keywords string, filter *SearchUserFilter) int {
	query, params := buildUserByOrganizationIdAndUsernameOrFirstnameOrLastnameOrEmailQuery(orgaId, keywords, filter, true)
	var count int

	if err := connection.Get(&count, query, params...); err != nil {
		logbuch.Error("Error counting user by organization id and query string", logbuch.Fields{"err": err, "organization_id": orgaId, "keywords": keywords})
		return 0
	}

	return count
}

func buildUserByOrganizationIdAndUsernameOrFirstnameOrLastnameOrEmailQuery(orgaId hide.ID, keywords string, filter *SearchUserFilter, count bool) (string, []interface{}) {
	params := make([]interface{}, 1)
	params[0] = orgaId
	selectOrCount := userBaseQueryHead

	if count {
		selectOrCount = "SELECT COUNT(1) "
	}

	var sb strings.Builder
	sb.WriteString(selectOrCount)
	sb.WriteString(userBaseQuery)

	if keywords != "" {
		params = append(params, keywords)
		sb.WriteString(`AND (SIMILARITY(username, $2) > 0.2
			OR LOWER(username) LIKE LOWER('%'||$2||'%')
			OR SIMILARITY(firstname, $2) > 0.2
			OR LOWER(firstname) LIKE LOWER('%'||$2||'%')
			OR SIMILARITY(lastname, $2) > 0.2
			OR LOWER(lastname) LIKE LOWER('%'||$2||'%')
			OR SIMILARITY(email, $2) > 0.2
			OR LOWER(email) LIKE LOWER('%'||$2||'%')) `)
	}

	// add field filter (joined with "AND")
	index := len(params) + 1
	fieldFilter, index, params := filter.addFieldFilter("", index, params, []string{filter.Username, filter.Firstname, filter.Lastname, filter.Email}, "username", "firstname", "lastname", "email")
	sb.WriteString(fieldFilter)

	// add date filter
	dateFilter, _, params := filter.addDateFilter("user", index, params)
	sb.WriteString(dateFilter)

	if !count {
		// sorting
		defaultFields := []SortValue{{`"user".id`, sortDirectionASC}}
		sb.WriteString(filter.addSorting("user", defaultFields, SortValue{"username", filter.SortUsername}, SortValue{"email", filter.SortEmail}, SortValue{"firstname", filter.SortFirstname}, SortValue{"lastname", filter.SortLastname}))

		// set limit
		var limit string
		limit, _, params = filter.addLimit(index, params)
		sb.WriteString(limit)
	}

	return sb.String(), params
}

func CountNewsletterSubscriptions() int {
	query := `SELECT COUNT(DISTINCT email) FROM (
			SELECT email FROM "user" WHERE accept_marketing IS TRUE
			UNION
			SELECT email FROM "newsletter_subscription" WHERE list IS NULL
		) AS email`
	var count int

	if err := connection.Get(&count, query); err != nil {
		logbuch.Error("Error counting newsletter subscriptions", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func FindNewsletterSubscriptionEmails() (*sqlx.Rows, error) {
	return connection.Queryx(`SELECT DISTINCT ON (email) email, newsletter_subscription_code FROM (
			SELECT email, NULL AS newsletter_subscription_code FROM "user" WHERE accept_marketing IS TRUE
			UNION
			SELECT email, code AS newsletter_subscription_code FROM "newsletter_subscription" WHERE list IS NULL
		) AS email`)
}

// SaveUser saves the given user. If create is set to true, a new user will be created.
// User is a special case because the ID is managed by auth and in a different schema,
// so that the generic save cannot be used here.
func SaveUser(tx *sqlx.Tx, entity *User, create bool) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var err error

	if create {
		_, err = tx.NamedExec(`INSERT INTO "user" (id,
			email,
			firstname,
			lastname,
			language,
			info,
			picture,
			accept_marketing,
			color_mode,
			introduction)
			VALUES (:id,
			:email,
			:firstname,
			:lastname,
			:language,
			:info,
			:picture,
			:accept_marketing,
			:color_mode,
			:introduction)`, entity)
	} else {
		_, err = tx.NamedExec(`UPDATE "user" SET email = :email,
			firstname = :firstname,
			lastname = :lastname,
			language = :language,
			info = :info,
			picture = :picture,
			accept_marketing = :accept_marketing,
			color_mode = :color_mode,
			introduction = :introduction
			WHERE id = :id`, entity)
	}

	if err != nil {
		logbuch.Error("Error saving user", logbuch.Fields{"err": err, "entity": entity})
		db.Rollback(tx)
		return err
	}

	return nil
}
