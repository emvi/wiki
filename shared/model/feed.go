package model

import (
	"emviwiki/shared/db"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

const (
	feedBaseQuery = `SELECT "feed".*,
		CASE WHEN "feed_access".read IS NULL THEN TRUE ELSE "feed_access".read END,
		"triggeredby".id "triggered_by_user.id",
		"triggeredby".email "triggered_by_user.email",
		"triggeredby".firstname "triggered_by_user.firstname",
		"triggeredby".lastname "triggered_by_user.lastname",
		"triggeredby".picture "triggered_by_user.picture",
		"triggeredbymember".username "triggered_by_user.organization_member.username"
		FROM "feed"
		LEFT JOIN "feed_access" ON "feed_access".feed_id = "feed".id AND "feed_access".user_id = $2
		JOIN "user" "triggeredby" ON "feed".triggered_by_user_id = "triggeredby".id
		JOIN "organization_member" "triggeredbymember" ON "triggeredby".id = "triggeredbymember".user_id AND "triggeredbymember".organization_id = "feed".organization_id `
)

type Feed struct {
	db.BaseEntity

	Public            bool        `json:"public"`
	Reason            string      `json:"reason"`
	OrganizationId    hide.ID     `db:"organization_id" json:"organization_id"`
	TriggeredByUserId hide.ID     `db:"triggered_by_user_id" json:"triggered_by_user_id"`
	RoomID            null.String `db:"room_id" json:"room_id"`

	TriggeredByUser *User     `db:"triggered_by_user" json:"triggered_by_user"`
	FeedRefs        []FeedRef `db:"-" json:"refs"`

	// feed and notification text from config
	Feed         string `db:"-" json:"feed"`
	Notification string `db:"-" json:"notification"`

	// is read from feed access
	Read bool `json:"read"`
}

func FindFeedByOrganizationIdAndReason(orgaId hide.ID, reason string) []Feed {
	var entities []Feed

	if err := connection.Select(&entities, `SELECT * FROM "feed" WHERE organization_id = $1 AND reason = $2`, orgaId, reason); err != nil {
		logbuch.Error("Error reading feed by organization id and reason", logbuch.Fields{"err": err, "organization_id": orgaId, "reason": reason})
		return nil
	}

	return entities
}

func FindFeedByOrganizationIdAndUserIdAndLanguageIdAndFilterLimit(orgaId, userId, langId hide.ID, filter *SearchFeedFilter) []Feed {
	query, params := buildFeedByOrganizationIdAndUserIdAndLanguageIdAndFilterLimitQuery(orgaId, userId, filter)
	var entities []Feed

	if err := connection.Select(&entities, query, params...); err != nil {
		logbuch.Error("Error reading feed by organization id and user id and language id and filter", logbuch.Fields{"err": err, "organization_id": orgaId, "user_id": userId, "lang_id": langId, "filter": filter})
		return nil
	}

	return joinFeedReferences(entities, orgaId, langId)
}

func buildFeedByOrganizationIdAndUserIdAndLanguageIdAndFilterLimitQuery(orgaId, userId hide.ID, filter *SearchFeedFilter) (string, []interface{}) {
	params := make([]interface{}, 2)
	params[0] = orgaId
	params[1] = userId

	var sb strings.Builder
	sb.WriteString(feedBaseQuery)
	sb.WriteString(`WHERE "feed".organization_id = $1
		AND ("feed".public IS TRUE OR EXISTS
			(SELECT 1 FROM "feed_access" WHERE feed_id = "feed".id AND user_id = $2)) `)

	if filter.Notifications {
		if filter.Unread {
			sb.WriteString(`AND EXISTS (SELECT 1 FROM "feed_access" WHERE feed_id = "feed".id AND user_id = $2 AND notification IS TRUE AND read IS FALSE) `)
		} else {
			sb.WriteString(`AND EXISTS (SELECT 1 FROM "feed_access" WHERE feed_id = "feed".id AND user_id = $2 AND notification IS TRUE) `)
		}
	}

	// add field filter (joined with "AND")
	index := 3

	if len(filter.UserIds) > 0 {
		var userFilter string
		userFilter, index, params = filter.filterInIds(filter.UserIds, index, params)
		sb.WriteString(fmt.Sprintf(`AND "triggeredbymember".user_id %v`, userFilter))
	}

	if len(filter.Reasons) > 0 {
		var reasonsFilter string
		reasonsFilter, index, params = filter.filterInStrings(filter.Reasons, index, params)
		sb.WriteString(fmt.Sprintf(`AND "feed".reason %v`, reasonsFilter))
	}

	// add date filter
	dateFilter, index, params := filter.addDateFilter("feed", index, params)
	sb.WriteString(dateFilter)

	// sorting
	defaultFields := []SortValue{{"def_time", sortDirectionDESC}}
	sb.WriteString(filter.addSorting("feed", defaultFields))

	// set limit
	limit, _, params := filter.addLimit(index, params)
	sb.WriteString(limit)

	return sb.String(), params
}

func FindNotificationByOrganizationIdAndUserIdAndLanguageIdAndAfterDefTimeUnread(orgaId, userId, langId hide.ID, minDefTime time.Time) []Feed {
	var entities []Feed

	query := feedBaseQuery + `WHERE "feed".organization_id = $1
		AND EXISTS (SELECT 1 FROM "feed_access" WHERE feed_id = "feed".id AND user_id = $2 AND notification IS TRUE AND read IS FALSE)
		AND "feed".def_time > $3`

	if err := connection.Select(&entities, query, orgaId, userId, minDefTime); err != nil {
		logbuch.Error("Error reading unread feed by organization id and user id", logbuch.Fields{"err": err, "organization_id": orgaId, "user_id": userId})
		return nil
	}

	return joinFeedReferences(entities, orgaId, langId)
}

func joinFeedReferences(feed []Feed, orgaId, langId hide.ID) []Feed {
	for i := range feed {
		feed[i].FeedRefs = FindFeedRefByOrganizationIdAndLanguageIdAndFeedId(orgaId, langId, feed[i].ID)
	}

	return feed
}

func DeleteFeedByIds(tx *sqlx.Tx, ids []hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query, args, err := sqlx.In(`DELETE FROM "feed" WHERE id IN (?)`, ids)

	if err != nil {
		logbuch.Error("Error deleting feed by ids", logbuch.Fields{"err": err})
		db.Rollback(tx)
		return err
	}

	query = connection.Rebind(query)
	_, err = tx.Exec(query, args...)

	if err != nil {
		logbuch.Error("Error deleting feed by ids", logbuch.Fields{"err": err})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveFeed(tx *sqlx.Tx, entity *Feed) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "feed" (public,
			reason,
			organization_id,
			triggered_by_user_id,
			room_id)
			VALUES (:public,
			:reason,
			:organization_id,
			:triggered_by_user_id,
			:room_id) RETURNING id`,
		`UPDATE "feed" SET public = :public,
			reason = :reason,
			organization_id = :organization_id,
			triggered_by_user_id = :triggered_by_user_id,
			room_id = :room_id
			WHERE id = :id`)
}
