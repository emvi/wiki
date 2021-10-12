package history

import (
	"emviwiki/backend/article/util"
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/backend/perm"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func DeleteHistoryEntry(ctx context.EmviContext, articleContentId hide.ID) error {
	content := model.GetArticleContentById(articleContentId)

	if content == nil {
		logbuch.Debug("Content version to delete history entry not found", logbuch.Fields{"id": articleContentId})
		return errs.ArticleContentVersionNotFound
	}

	if !perm.CheckUserWriteAccess(content.ArticleId, ctx.UserId) {
		return errs.PermissionDenied
	}

	// latest version cannot be deleted, but it will be updated if the user tries to delete the latest history entry
	if content.Version == 0 {
		logbuch.Debug("Content version to delete history entry invalid", logbuch.Fields{"id": articleContentId})
		return errs.ArticleContentVersionInvalid
	}

	// if this is the last and only version, do not deleted it
	if model.CountArticleContentVersionByArticleIdAndLanguageIdAndNotWIP(content.ArticleId, content.LanguageId) <= 1 {
		return errs.ArticleContentRemainingVersion
	}

	lastContent := model.GetArticleContentLastByArticleIdAndLanguageIdAndWIP(content.ArticleId, content.LanguageId, false)

	if err := util.CheckContentVersionRequiresExpert(ctx.Organization.Expert, content.Version, lastContent.Version); err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to delete history entry", logbuch.Fields{"err": err, "content_id": content.ID})
		return errs.TxBegin
	}

	if err := feed.DeleteFeed(tx, &feed.DeleteFeedData{ArticleContentId: content.ID}); err != nil {
		db.Rollback(tx)
		logbuch.Error("Error deleting feed entries for history entry", logbuch.Fields{"err": err, "content_id": content.ID})
		return errs.Saving
	}

	if content.Version == lastContent.Version {
		if err := deleteArticleHistoryEntryAndUpdateLatest(tx, content); err != nil {
			db.Rollback(tx)
			return errs.Saving
		}
	} else {
		if err := deleteArticleHistoryEntry(tx, content.ID); err != nil {
			db.Rollback(tx)
			return errs.Saving
		}
	}

	if err := createDeleteArticleHistoryFeed(tx, ctx, content.ArticleId, content.LanguageId); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction to delete history entry", logbuch.Fields{"err": err, "content_id": content.ID})
		return errs.TxCommit
	}

	return nil
}

func deleteArticleHistoryEntryAndUpdateLatest(tx *sqlx.Tx, content *model.ArticleContent) error {
	if err := deleteArticleHistoryEntry(tx, content.ID); err != nil {
		return err
	}

	// set ID of last to ID of latest to overwrite latest with the last version
	lastContent := model.GetArticleContentLastByArticleIdAndLanguageIdAndWIPTx(tx, content.ArticleId, content.LanguageId, false)
	latestContent := model.GetArticleContentLatestByArticleIdAndLanguageIdTx(tx, content.ArticleId, content.LanguageId, false)
	lastContent.ID = latestContent.ID
	lastContent.Version = 0

	if err := model.SaveArticleContent(tx, lastContent); err != nil {
		logbuch.Error("Error saving new latest article content when deleting history entry", logbuch.Fields{"err": err, "content_id": content.ID})
		return err
	}

	return nil
}

func deleteArticleHistoryEntry(tx *sqlx.Tx, contentId hide.ID) error {
	if err := model.DeleteArticleContentById(tx, contentId); err != nil {
		logbuch.Error("Error deleting article content when deleting history entry", logbuch.Fields{"err": err, "content_id": contentId})
		return err
	}

	return nil
}

func createDeleteArticleHistoryFeed(tx *sqlx.Tx, ctx context.EmviContext, articleId, langId hide.ID) error {
	article, err := util.GetArticleWithAccess(tx, ctx, articleId, false)

	if err != nil {
		db.Rollback(tx)
		logbuch.Error("Error reading article for feed when deleting article history entry", logbuch.Fields{"err": err})
		return errs.Saving
	}

	refs := make([]interface{}, 2)
	refs[0] = article
	refs[1] = model.GetArticleContentLatestByArticleIdAndLanguageIdTx(tx, articleId, langId, false)
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: ctx.Organization,
		UserId:       ctx.UserId,
		Reason:       "delete_article_history_entry",
		Public:       article.ReadEveryone || article.WriteEveryone,
		Access:       perm.GetUserIdsFromAccess(tx, article.Access),
		Notify:       model.FindObservedObjectUserIdByArticleIdOrArticleListIdTx(tx, article.ID, 0),
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when deleting article history entry", logbuch.Fields{"err": err})
		return errs.Saving
	}

	return nil
}
