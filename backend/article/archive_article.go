package article

import (
	articleutil "emviwiki/backend/article/util"
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	archiveMessageMaxLen = 100
)

func ArchiveArticle(organization *model.Organization, userId, articleId hide.ID, message string, delete bool) error {
	article, err := checkArchiveDeleteArticleAccess(organization, userId, articleId)

	if err != nil {
		return err
	}

	if article.Archived.Valid {
		article.Archived.SetNil()
	} else {
		message = strings.TrimSpace(message)

		if !delete {
			if message == "" {
				return errs.MessageTooShort
			}

			if len(message) > archiveMessageMaxLen {
				return errs.MessageTooLong
			}
		}

		article.Archived.SetValid(message)
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error beginning transaction to archive article", logbuch.Fields{"err": err, "orga_id": organization.ID, "user_id": userId, "article_id": articleId})
		return errs.TxBegin
	}

	if err := model.SaveArticle(tx, article); err != nil {
		logbuch.Error("Error archiving article", logbuch.Fields{"err": err, "orga_id": organization.ID, "user_id": userId, "article_id": articleId})
		return errs.Saving
	}

	if err := createArchivedArticleFeed(tx, organization, userId, article); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when archiving article", logbuch.Fields{"err": err, "orga_id": organization.ID, "user_id": userId, "article_id": articleId})
		return errs.TxCommit
	}

	if delete {
		if err := DeleteArticle(organization, userId, articleId); err != nil {
			return err
		}
	}

	return nil
}

func checkArchiveDeleteArticleAccess(orga *model.Organization, userId, articleId hide.ID) (*model.Article, error) {
	article, err := articleutil.GetArticleWithAccess(nil, context.NewEmviUserContext(orga, userId), articleId, true)

	// ignore permission denied here, since we check that in the next step (for admin/mod in addition)
	if err != nil && err != errs.PermissionDenied {
		return nil, err
	}

	if !hasWriteAccess(article, userId) && !checkUserIsModeratorOrAdmin(orga.ID, userId) {
		return nil, errs.PermissionDenied
	}

	return article, nil
}

func checkUserIsModeratorOrAdmin(orgaId, userId hide.ID) bool {
	member := model.GetOrganizationMemberByOrganizationIdAndUserId(orgaId, userId)

	if member == nil {
		logbuch.Error("Error finding organization member to check if he/she is admin/moderator", logbuch.Fields{"orga_id": orgaId, "user_id": userId})
		return false
	}

	return member.IsAdmin || member.IsModerator
}

func createArchivedArticleFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, article *model.Article) error {
	reason := "restored_article"

	if article.Archived.Valid {
		reason = "archived_article"
	}

	langId := util.DetermineLang(tx, orga.ID, userId, 0).ID
	refs := make([]interface{}, 2)
	refs[0] = article
	refs[1] = model.GetArticleContentLatestByArticleIdAndLanguageIdTx(tx, article.ID, langId, false)
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       reason,
		Public:       article.ReadEveryone || article.WriteEveryone,
		Access:       perm.GetUserIdsFromAccess(tx, model.FindArticleAccessByOrganizationIdAndArticleIdTx(tx, orga.ID, article.ID)),
		Notify:       model.FindObservedObjectUserIdByArticleIdOrArticleListIdTx(tx, article.ID, 0),
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when archiving article", logbuch.Fields{"err": err})
		return err
	}

	return nil
}
