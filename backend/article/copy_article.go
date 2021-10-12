package article

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/backend/perm"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

var copyTitle = map[string]string{
	"en": "copy",
	"de": "Kopie",
}

func CopyArticle(orga *model.Organization, userId, articleId, langId hide.ID) (hide.ID, error) {
	article := model.GetArticleByOrganizationIdAndId(orga.ID, articleId)

	if article == nil {
		return 0, errs.ArticleNotFound
	}

	if !article.ReadEveryone && !perm.CheckUserReadOrWriteAccess(articleId, userId) {
		return 0, errs.PermissionDenied
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to copy article", logbuch.Fields{"err": err})
		return 0, errs.TxBegin
	}

	newArticle, err := copyArticle(tx, orga.ID, articleId)

	if err != nil {
		return 0, err
	}

	latestContent, err := copyArticleContent(tx, orga.ID, userId, articleId, langId, newArticle.ID)

	if err != nil {
		return 0, err
	}

	if err := copyArticleAccess(tx, orga.ID, article, newArticle); err != nil {
		return 0, err
	}

	if err := copyArticleTags(tx, orga.ID, userId, articleId, newArticle.ID); err != nil {
		return 0, err
	}

	if err := createCopiedArticleFeed(tx, orga, userId, article, latestContent); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error when committing transaction to copy article", logbuch.Fields{"err": err})
		return 0, errs.TxCommit
	}

	return newArticle.ID, nil
}

func copyArticle(tx *sqlx.Tx, orgaId, articleId hide.ID) (*model.Article, error) {
	article := model.GetArticleByOrganizationIdAndIdTx(tx, orgaId, articleId)

	if article == nil {
		db.Rollback(tx)
		return nil, errs.ArticleNotFound
	}

	article.ID = 0
	article.Pinned = false
	article.Views = 0

	if err := model.SaveArticle(tx, article); err != nil {
		logbuch.Error("Error saving article when copying article", logbuch.Fields{"err": err, "orga_id": orgaId, "article_id": articleId})
		return nil, errs.Saving
	}

	return article, nil
}

func copyArticleContent(tx *sqlx.Tx, orgaId, userId, articleId, langId, newArticleId hide.ID) (*model.ArticleContent, error) {
	langId = util.DetermineLang(tx, orgaId, userId, langId).ID
	content := model.FindArticleContentByArticleIdTx(tx, articleId)
	var latestContent model.ArticleContent // create explicit copy

	for _, c := range content {
		// determine latest article content of the current article for feed
		if c.Version == 0 && c.LanguageId == langId {
			latestContent = c
		} else if c.Version == 0 && latestContent.ID == 0 {
			latestContent = c
		}

		if err := copySingleArticleContent(tx, orgaId, articleId, newArticleId, &c); err != nil {
			return nil, err
		}
	}

	return &latestContent, nil
}

func copySingleArticleContent(tx *sqlx.Tx, orgaId, articleId, newArticleId hide.ID, content *model.ArticleContent) error {
	contentId := content.ID
	content.ID = 0
	content.ArticleId = newArticleId

	// append copy to title for latest version
	if content.Version == 0 {
		lang := model.GetLanguageByOrganizationIdAndIdTx(tx, orgaId, content.LanguageId)
		langCode := util.GetSystemSupportedLangCode(lang.Code)
		content.Title += fmt.Sprintf(" (%s)", copyTitle[langCode])
	}

	if err := model.SaveArticleContent(tx, content); err != nil {
		logbuch.Error("Error saving article content when copying article", logbuch.Fields{"err": err, "article_id": articleId})
		return errs.Saving
	}

	if err := copyArticleContentAuthor(tx, contentId, content.ID); err != nil {
		return err
	}

	return nil
}

func copyArticleContentAuthor(tx *sqlx.Tx, contentId, newContentId hide.ID) error {
	authors := model.FindArticleContentAuthorByArticleContentIdTx(tx, contentId)

	for _, a := range authors {
		a.ID = 0
		a.ArticleContentId = newContentId

		if err := model.SaveArticleContentAuthor(tx, &a); err != nil {
			logbuch.Error("Error saving article content author when copying article", logbuch.Fields{"err": err, "content_id": contentId})
			return errs.Saving
		}
	}

	return nil
}

func copyArticleAccess(tx *sqlx.Tx, orgaId hide.ID, article, newArticle *model.Article) error {
	access := model.FindArticleAccessByOrganizationIdAndArticleIdTx(tx, orgaId, article.ID)
	article.Access = access

	for _, a := range access {
		a.ID = 0
		a.ArticleId = newArticle.ID

		if err := model.SaveArticleAccess(tx, &a); err != nil {
			logbuch.Error("Error saving article access when copying article", logbuch.Fields{"err": err, "orga_id": orgaId, "article_id": article.ID})
			return errs.Saving
		}
	}

	return nil
}

func copyArticleTags(tx *sqlx.Tx, orgaId, userId, articleId, newArticleId hide.ID) error {
	tags := model.FindTagByOrganizationIdAndUserIdAndArticleIdTx(tx, orgaId, userId, articleId)

	for _, t := range tags {
		tag := &model.ArticleTag{ArticleId: newArticleId, TagId: t.ID}

		if err := model.SaveArticleTag(tx, tag); err != nil {
			logbuch.Error("Error saving tag when copying article", logbuch.Fields{"err": err, "orga_id": orgaId, "article_id": articleId})
			return errs.Saving
		}
	}

	return nil
}

// This takes the old article and content as parameters! Do not create a feed entry for the new article!
func createCopiedArticleFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, article *model.Article, latestContent *model.ArticleContent) error {
	refs := make([]interface{}, 2)
	refs[0] = article
	refs[1] = latestContent
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "copy_article",
		Public:       article.ReadEveryone || article.WriteEveryone,
		Access:       perm.GetUserIdsFromAccess(tx, article.Access),
		Notify:       model.FindObservedObjectUserIdByArticleIdOrArticleListIdTx(tx, article.ID, 0),
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when copying article", logbuch.Fields{"err": err})
		return errs.Saving
	}

	return nil
}
