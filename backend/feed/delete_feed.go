package feed

import (
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/jmoiron/sqlx"
)

type DeleteFeedData struct {
	ArticleId        hide.ID
	ArticleContentId hide.ID
	ListId           hide.ID
	GroupId          hide.ID
}

// DeleteFeed deletes all feeds including all references belonging to a referenced object id (article, list, group, ...).
func DeleteFeed(tx *sqlx.Tx, data *DeleteFeedData) error {
	var feedIds []hide.ID

	if data.ArticleId != 0 {
		feedIds = model.FindFeedRefFeedIdsByArticleId(tx, data.ArticleId)
	} else if data.ArticleContentId != 0 {
		feedIds = model.FindFeedRefFeedIdsByArticleContentId(tx, data.ArticleContentId)
	} else if data.ListId != 0 {
		feedIds = model.FindFeedRefFeedIdsByArticleListId(tx, data.ListId)
	} else if data.GroupId != 0 {
		feedIds = model.FindFeedRefFeedIdsByUserGroupId(tx, data.GroupId)
	}

	return deleteFeeds(tx, feedIds)
}

// no rollback, since all delete functions rollback transactions on error
func deleteFeeds(tx *sqlx.Tx, feedIds []hide.ID) error {
	if len(feedIds) == 0 {
		return nil
	}

	if err := model.DeleteFeedAccessByFeedIds(tx, feedIds); err != nil {
		return err
	}

	if err := model.DeleteFeedRefByFeedIds(tx, feedIds); err != nil {
		return err
	}

	if err := model.DeleteFeedByIds(tx, feedIds); err != nil {
		return err
	}

	return nil
}
