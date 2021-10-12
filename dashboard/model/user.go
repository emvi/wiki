package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/logbuch"
)

type User struct {
	db.BaseEntity

	Email        string `json:"email"`
	Password     string `db:"password" json:"-"`
	PasswordSalt string `db:"password_salt" json:"-"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
}

func GetUserByEmail(email string) *User {
	entity := new(User)

	if err := dashboardDB.Get(entity, `SELECT * FROM "user" WHERE LOWER(email) = LOWER($1)`, email); err != nil {
		logbuch.Debug("User by email not found", logbuch.Fields{"err": err, "email": email})
		return nil
	}

	return entity
}
