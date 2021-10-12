package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type User struct {
	db.BaseEntity

	Email                 string
	Password              null.String
	PasswordSalt          null.String `db:"password_salt"`
	ResetPassword         bool        `db:"reset_password"`
	LastPasswordReset     pq.NullTime `db:"last_password_reset"`
	Firstname             null.String
	Lastname              null.String
	Language              null.String
	PictureURL            null.String `db:"picture_url"`
	AcceptMarketing       bool        `db:"accept_marketing"`
	Active                bool
	LastLogin             time.Time   `db:"last_login"` // default: now
	RegistrationCode      null.String `db:"registration_code"`
	RegistrationStep      int         `db:"registration_step"`
	RegistrationMailsSend int         `db:"registration_mails_send"`
	NewEmail              null.String `db:"new_email"`
	NewEmailCode          null.String `db:"new_email_code"`
	LoginAttempts         int         `db:"login_attempts"`
	LastLoginAttempt      time.Time   `db:"last_login_attempt"`
	AuthProvider          string      `db:"auth_provider" json:"-"`
	AuthProviderUserId    null.String `db:"auth_provider_user_id" json:"-"`
}

func GetUserById(id hide.ID) *User {
	user := new(User)

	if err := connection.Get(user, `SELECT * FROM "user" WHERE id = $1 AND active = TRUE`, id); err != nil {
		logbuch.Debug("User by ID not found", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return user
}

func GetUserByEmailAndPassword(email, password string) *User {
	user := new(User)

	if err := connection.Get(user, `SELECT * FROM "user" WHERE LOWER(email) = LOWER($1) AND password = $2 AND active = TRUE AND auth_provider = 'emvi'`, email, password); err != nil {
		logbuch.Debug("User by email and password not found", logbuch.Fields{"err": err, "email": email})
		return nil
	}

	return user
}

func GetUserByIdAndPassword(id hide.ID, password string) *User {
	user := new(User)

	if err := connection.Get(user, `SELECT * FROM "user" WHERE id = $1 AND password = $2 AND active = TRUE`, id, password); err != nil {
		logbuch.Debug("User by id and password not found", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return user
}

func GetUserByEmailIgnoreActive(email string) *User {
	return GetUserByEmailIgnoreActiveTx(nil, email)
}

func GetUserByEmailIgnoreActiveTx(tx *sqlx.Tx, email string) *User {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	user := new(User)

	if err := tx.Get(user, `SELECT * FROM "user" WHERE LOWER(email) = LOWER($1) AND auth_provider = 'emvi'`, email); err != nil {
		logbuch.Debug("User by email not found", logbuch.Fields{"err": err, "email": email})
		return nil
	}

	return user
}

func GetUserByEmail(email string) *User {
	return GetUserByEmailTx(nil, email)
}

func GetUserByEmailTx(tx *sqlx.Tx, email string) *User {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	user := new(User)

	if err := tx.Get(user, `SELECT * FROM "user" WHERE LOWER("email") = LOWER($1) AND "active" = TRUE AND auth_provider = 'emvi'`, email); err != nil {
		logbuch.Debug("User by email not found", logbuch.Fields{"err": err, "email": email})
		return nil
	}

	return user
}

func GetUserByAuthProviderAndUserId(authProvider string, id string) *User {
	user := new(User)

	if err := connection.Get(user, `SELECT * FROM "user" WHERE auth_provider = $1 AND auth_provider_user_id = $2`, authProvider, id); err != nil {
		logbuch.Debug("User by auth provider and id not found", logbuch.Fields{"err": err, "auth_provider": authProvider, "id": id})
		return nil
	}

	return user
}

func GetUserByRegistrationCodeAndInactive(tx *sqlx.Tx, code string) *User {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	user := new(User)

	if err := tx.Get(user, `SELECT * FROM "user" WHERE registration_code = $1 AND active = FALSE`, code); err != nil {
		logbuch.Debug("User by registration code not found", logbuch.Fields{"err": err, "code": code})
		return nil
	}

	return user
}

func GetUserByNewEmailAndNewEmailCode(email, code string) *User {
	user := new(User)

	if err := connection.Get(user, `SELECT * FROM "user" WHERE LOWER(new_email) = LOWER($1) AND new_email_code = $2 AND active = TRUE`, email, code); err != nil {
		logbuch.Debug("User by new email and new email code not found", logbuch.Fields{"err": err, "email": email, "code": code})
		return nil
	}

	return user
}

func DeleteUserByInactiveAndRegistrationStepAndDefTime(tx *sqlx.Tx) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := connection.Exec(tx, `DELETE FROM "user" WHERE active IS FALSE AND registration_step < 4 AND def_time < NOW() - INTERVAL '1 month'`); err != nil {
		return err
	}

	return nil
}

func DeleteUserById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := connection.Exec(tx, `DELETE FROM "user" WHERE id = $1`, id); err != nil {
		return err
	}

	return nil
}

func SaveUser(tx *sqlx.Tx, entity *User) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "user" (email,
			password,
			password_salt,
			reset_password,
			last_password_reset,
			firstname,
			lastname,
			language,
			picture_url,
			accept_marketing,
			active,
			last_login,
			registration_code,
			registration_step,
			registration_mails_send,
			new_email,
			new_email_code,
			login_attempts,
			last_login_attempt,
			auth_provider,
			auth_provider_user_id)
			VALUES (:email,
			:password,
			:password_salt,
			:reset_password,
			:last_password_reset,
			:firstname,
			:lastname,
			:language,
			:picture_url,
			:accept_marketing,
			:active,
			:last_login,
			:registration_code,
			:registration_step,
			:registration_mails_send,
			:new_email,
			:new_email_code,
			:login_attempts,
			:last_login_attempt,
			:auth_provider,
			:auth_provider_user_id) RETURNING id`,
		`UPDATE "user" SET email = :email,
			password = :password,
			password_salt = :password_salt,
			reset_password = :reset_password,
			last_password_reset = :last_password_reset,
			firstname = :firstname,
			lastname = :lastname,
			language = :language,
			picture_url = :picture_url,
			accept_marketing = :accept_marketing,
			active = :active,
			last_login = :last_login,
			registration_code = :registration_code,
			registration_step = :registration_step,
			registration_mails_send = :registration_mails_send,
			new_email = :new_email,
			new_email_code = :new_email_code,
			login_attempts = :login_attempts,
			last_login_attempt = :last_login_attempt,
			auth_provider = :auth_provider,
			auth_provider_user_id = :auth_provider_user_id
			WHERE id = :id`)
}
