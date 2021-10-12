package feed

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/db"
	"emviwiki/shared/feed"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"strings"
)

// RoomID for new articles. This can be used to reference a newly created article without an ID yet.
type RoomID string

// KeyValue is a key value pair that can be used to store generic data for a feed text.
// It will be stored as a FeedRef object.
type KeyValue struct {
	Key   string
	Value string
}

// CreateFeedData is used to create a new feed entry and notifications.
// Access and Notify refer to user IDs.
// Refs must be one of the types specified in model.Feed.
// Access must not explicitly granted to users that are notified. They are added to the access list by CreateFeed.
type CreateFeedData struct {
	Tx           *sqlx.Tx
	Organization *model.Organization
	UserId       hide.ID // creating user ID
	Reason       string
	Public       bool
	Access       []hide.ID
	Notify       []hide.ID

	// List of referenced objects or a key value pair, like an article for example.
	Refs []interface{}
}

// CreateFeed creates a new feed entry and/or notifications for given data.
func CreateFeed(data *CreateFeedData) error {
	if !data.Public && len(data.Access)+len(data.Notify) == 0 {
		if data.Tx != nil {
			db.Rollback(data.Tx)
		}

		return errs.NonPublicFeedWithoutAccess
	}

	if !feed.CheckReasonExists(data.Reason) {
		if data.Tx != nil {
			db.Rollback(data.Tx)
		}

		return errs.ReasonNotFound
	}

	// create new transaction or use the one provided
	var tx *sqlx.Tx

	if data.Tx == nil {
		var err error
		tx, err = model.GetConnection().Beginx()

		if err != nil {
			logbuch.Error("Error starting transaction to create feed", logbuch.Fields{"err": err})
			return err
		}
	} else {
		tx = data.Tx
	}

	newFeed := &model.Feed{OrganizationId: data.Organization.ID,
		TriggeredByUserId: data.UserId,
		Reason:            data.Reason,
		Public:            data.Public}

	if err := setRefs(newFeed, data.Refs); err != nil {
		db.Rollback(tx)
		return err
	}

	if err := model.SaveFeed(tx, newFeed); err != nil {
		return err
	}

	for _, ref := range newFeed.FeedRefs {
		ref.FeedId = newFeed.ID

		if err := model.SaveFeedRef(tx, &ref); err != nil {
			return err
		}
	}

	if err := createAccessAndNotifications(tx, newFeed, data); err != nil {
		return err
	}

	// only commit if the transaction was created here
	if data.Tx == nil {
		if err := tx.Commit(); err != nil {
			logbuch.Error("Error committing transaction to create feed", logbuch.Fields{"err": err})
			return err
		}
	}

	return nil
}

func setRefs(feed *model.Feed, refs []interface{}) error {
	for _, ref := range refs {
		switch t := ref.(type) {
		case *model.User:
			if t == nil {
				continue
			}

			feed.FeedRefs = append(feed.FeedRefs, model.FeedRef{UserID: t.ID, User: t})
		case *model.UserGroup:
			if t == nil {
				continue
			}

			feed.FeedRefs = append(feed.FeedRefs, model.FeedRef{UserGroupID: t.ID, UserGroup: t})
		case *model.Article:
			if t == nil {
				continue
			}

			feed.FeedRefs = append(feed.FeedRefs, model.FeedRef{ArticleID: t.ID, Article: t})
		case *model.ArticleContent:
			if t == nil {
				continue
			}

			feed.FeedRefs = append(feed.FeedRefs, model.FeedRef{ArticleContentID: t.ID, ArticleContent: t})
		case *model.ArticleList:
			if t == nil {
				continue
			}

			feed.FeedRefs = append(feed.FeedRefs, model.FeedRef{ArticleListID: t.ID, ArticleList: t})
		case RoomID:
			feed.RoomID.SetValid(string(t))
		case KeyValue:
			t.Key = strings.TrimSpace(t.Key)
			t.Value = strings.TrimSpace(t.Value)

			if t.Key == "" || t.Value == "" {
				continue
			}

			feed.FeedRefs = append(feed.FeedRefs, model.FeedRef{Key: null.NewString(t.Key, true), Value: null.NewString(t.Value, true)})
		default:
			return errs.RefObjectUnknown
		}
	}

	return nil
}

func createAccessAndNotifications(tx *sqlx.Tx, feed *model.Feed, data *CreateFeedData) error {
	access := make(map[hide.ID]model.FeedAccess)
	access = appendAccess(feed, data.UserId, data.Notify, access, true)

	if !data.Public {
		// grant access if this is not pure notification
		if len(data.Access) != 0 {
			data.Access = append(data.Access, data.UserId)
		}

		access = appendAccess(feed, data.UserId, data.Access, access, false)
	}

	// TODO optimize (bulk insert)
	for _, a := range access {
		if err := model.SaveFeedAccess(tx, &a); err != nil {
			return err
		}
	}

	return nil
}

func appendAccess(feed *model.Feed, creatingUserId hide.ID, add []hide.ID, list map[hide.ID]model.FeedAccess, notify bool) map[hide.ID]model.FeedAccess {
	for _, user := range add {
		_, found := list[user]

		if !found {
			list[user] = model.FeedAccess{
				UserId:       user,
				FeedId:       feed.ID,
				Notification: notify && user != creatingUserId,
				Read:         user == creatingUserId || !notify,
			}
		}
	}

	return list
}
