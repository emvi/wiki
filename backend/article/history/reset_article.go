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
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"strings"
)

func ResetArticle(orga *model.Organization, userId, articleId, langId hide.ID, version int, commit string) error {
	article, err := util.GetArticleWithAccess(nil, context.NewEmviUserContext(orga, userId), articleId, false)

	if err != nil {
		return err
	}

	if !article.WriteEveryone && !perm.CheckUserWriteAccess(articleId, userId) {
		return errs.PermissionDenied
	}

	lastContent, content, err := validateVersionAndCommitMsg(orga, articleId, langId, version, commit)

	if err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error beginning transaction when resetting article to content version", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "article_id": articleId})
		return errs.TxBegin
	}

	if err := resetContentToVersion(tx, lastContent, content, userId, commit); err != nil {
		return err
	}

	if err := createResetArticleFeed(tx, orga, userId, article, content); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when resetting article to content version", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "article_id": articleId})
		return errs.TxCommit
	}

	return nil
}

func validateVersionAndCommitMsg(orga *model.Organization, articleId, langId hide.ID, version int, commit string) (*model.ArticleContent, *model.ArticleContent, error) {
	if version <= 0 {
		return nil, nil, errs.ArticleContentVersionInvalid
	}

	commit = strings.TrimSpace(commit)

	if err := util.CheckCommitMsg(commit); err != nil {
		return nil, nil, err
	}

	lastContent := model.GetArticleContentLastByArticleIdAndLanguageIdAndWIP(articleId, langId, false)

	if lastContent == nil {
		return nil, nil, errs.FindingLatestArticleContent
	}

	if err := util.CheckContentVersionRequiresExpert(orga.Expert, version, lastContent.Version); err != nil {
		return nil, nil, err
	}

	content := model.GetArticleContentByArticleIdAndLanguageIdAndVersion(articleId, langId, version)

	if content == nil {
		return nil, nil, errs.ArticleContentVersionNotFound
	}

	if content.Version == lastContent.Version {
		return nil, nil, errs.ArticleContentVersionInvalid
	}

	return lastContent, content, nil
}

func resetContentToVersion(tx *sqlx.Tx, lastContent, content *model.ArticleContent, userId hide.ID, commit string) error {
	// create new commit for reset
	content.ID = 0
	content.Version = lastContent.Version + 1
	content.Commit = null.NewString(commit, commit != "")
	content.WIP = false
	content.UserId = userId

	if err := model.SaveArticleContent(tx, content); err != nil {
		return errs.Saving
	}

	contentAuthor := &model.ArticleContentAuthor{ArticleContentId: content.ID, UserId: userId}

	if err := model.SaveArticleContentAuthor(tx, contentAuthor); err != nil {
		return errs.Saving
	}

	// update latest commit
	latestContent := model.GetArticleContentLatestByArticleIdAndLanguageIdTx(tx, content.ArticleId, content.LanguageId, false)

	if latestContent == nil {
		db.Rollback(tx)
		return errs.FindingLatestArticleContent
	}

	latestContent.Title = content.Title
	latestContent.Content = content.Content
	latestContent.Commit = null.NewString(commit, commit != "")
	latestContent.WIP = false
	latestContent.UserId = userId
	latestContent.ReadingTime = content.ReadingTime

	if err := model.SaveArticleContent(tx, latestContent); err != nil {
		return errs.Saving
	}

	return nil
}

func createResetArticleFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, article *model.Article, content *model.ArticleContent) error {
	refs := make([]interface{}, 2)
	refs[0] = article
	refs[1] = content
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "reset_article",
		Public:       article.ReadEveryone || article.WriteEveryone,
		Access:       perm.GetUserIdsFromAccess(tx, article.Access),
		Notify:       model.FindObservedObjectUserIdByArticleIdOrArticleListIdTx(tx, article.ID, 0),
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when resetting article", logbuch.Fields{"err": err})
		return errs.Saving
	}

	return nil
}
