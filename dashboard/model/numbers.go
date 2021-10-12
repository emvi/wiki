package model

import (
	"github.com/emvi/logbuch"
	"time"
)

type Statistics struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func CountOrganizations() int {
	var count int

	if err := backendDB.Get(&count, `SELECT COUNT(1) FROM "organization"`); err != nil {
		logbuch.Error("Error counting organizations", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func CountUser() int {
	var count int

	if err := backendDB.Get(&count, `SELECT COUNT(1) FROM "user"`); err != nil {
		logbuch.Error("Error counting organizations", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func CountArticles() int {
	var count int

	if err := backendDB.Get(&count, `SELECT COUNT(1) FROM "article"`); err != nil {
		logbuch.Error("Error counting organizations", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func CountLists() int {
	var count int

	if err := backendDB.Get(&count, `SELECT COUNT(1) FROM "article_list"`); err != nil {
		logbuch.Error("Error counting organizations", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func CountGroups() int {
	var count int

	if err := backendDB.Get(&count, `SELECT COUNT(1) FROM "user_group"`); err != nil {
		logbuch.Error("Error counting organizations", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func CountTags() int {
	var count int

	if err := backendDB.Get(&count, `SELECT COUNT(1) FROM "tag"`); err != nil {
		logbuch.Error("Error counting organizations", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func CountNewsletterWithoutList() int {
	var count int

	if err := backendDB.Get(&count, `SELECT COUNT(1) FROM "newsletter_subscription" WHERE list IS NULL`); err != nil {
		logbuch.Error("Error counting newsletter without list", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func CountNewsletterWithListOnPremise() int {
	var count int

	if err := backendDB.Get(&count, `SELECT COUNT(1) FROM "newsletter_subscription" WHERE list = 'onpremise'`); err != nil {
		logbuch.Error("Error counting newsletter with list on premise", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func CountUserLoggedInAfter(after time.Time) int {
	var count int

	if err := authDB.Get(&count, `SELECT COUNT(DISTINCT user_id) FROM "login" WHERE def_time > $1`, after); err != nil {
		logbuch.Error("Error counting users logged in after", logbuch.Fields{"err": err})
		return 0
	}

	return count
}

func FindLoginStatisticAfter(start time.Time) []Statistics {
	query := `SELECT date("date"), (SELECT count(1) FROM "login" WHERE date(def_time) = "date") "count"
		FROM (SELECT * FROM generate_series(date($1), date(now()), interval '1 day') "date") AS date_series
		ORDER BY "date" ASC`
	var entities []Statistics

	if err := authDB.Select(&entities, query, start); err != nil {
		logbuch.Error("Error reading login statistics", logbuch.Fields{"err": err, "start": start})
		return nil
	}

	return entities
}

func FindRegistrationStatisticAfter(start time.Time) []Statistics {
	query := `SELECT date("date"), (SELECT count(1) FROM "user" WHERE active IS TRUE AND date(def_time) = "date") "count"
		FROM (SELECT * FROM generate_series(date($1), date(now()), interval '1 day') "date") AS date_series
		ORDER BY "date" ASC`
	var entities []Statistics

	if err := authDB.Select(&entities, query, start); err != nil {
		logbuch.Error("Error reading login statistics", logbuch.Fields{"err": err, "start": start})
		return nil
	}

	return entities
}
