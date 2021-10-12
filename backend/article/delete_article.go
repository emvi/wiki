package article

import (
	"emviwiki/backend/content"
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/backend/perm"
	"emviwiki/backend/tag"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func DeleteArticle(orga *model.Organization, userId, articleId hide.ID) error {
	article, err := checkArchiveDeleteArticleAccess(orga, userId, articleId)

	if err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error beginning transaction when deleting article", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	content.DeleteFileForArticle(tx, orga, userId, articleId, nil)

	// create feed if public before deleting article
	if article.ReadEveryone || article.WriteEveryone {
		if err := createDeletedArticleFeed(tx, orga, userId, article); err != nil {
			logbuch.Error("Error creating feed when deleting article", logbuch.Fields{"err": err})
			return errs.Saving
		}
	}

	if err := feed.DeleteFeed(tx, &feed.DeleteFeedData{ArticleId: articleId}); err != nil {
		logbuch.Error("Error deleting article feed", logbuch.Fields{"err": err, "article_id": article.ID, "user_id": userId})
		return errs.Saving
	}

	if err := model.DeleteArticleById(tx, articleId); err != nil {
		logbuch.Error("Error deleting article", logbuch.Fields{"err": err, "article_id": article.ID, "user_id": userId})
		return errs.Saving
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when deleting article", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	tag.CleanupUnusedTags(orga.ID)
	return nil
}

func createDeletedArticleFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, article *model.Article) error {
	langId := util.DetermineLang(tx, orga.ID, userId, 0).ID
	latestContent := model.GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageIdTx(tx, orga.ID, article.ID, langId, false)

	if latestContent == nil {
		db.Rollback(tx)
		return errs.FindingLatestArticleContent
	}

	refs := make([]interface{}, 1)
	refs[0] = feed.KeyValue{"name", latestContent.Title}
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "delete_article",
		Public:       true,
		Access:       perm.GetUserIdsFromAccess(tx, article.Access),
		Notify:       model.FindObservedObjectUserIdByArticleIdOrArticleListIdTx(tx, article.ID, 0),
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
