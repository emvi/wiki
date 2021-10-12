package util

import (
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/jmoiron/sqlx"
)

const (
	maxCommitMsgLen            = 100
	ArticleMaxHistoryNonExpert = 3
)

// GetArticleWithAccess returns an article with access permissions for given context and ID.
// If ignoreArchived is set to true it will look for archived articles too.
func GetArticleWithAccess(tx *sqlx.Tx, ctx context.EmviContext, articleId hide.ID, ignoreArchived bool) (*model.Article, error) {
	var article *model.Article

	if ignoreArchived {
		article = model.GetArticleByOrganizationIdAndIdIgnoreArchivedTx(tx, ctx.Organization.ID, articleId)
	} else {
		article = model.GetArticleByOrganizationIdAndIdTx(tx, ctx.Organization.ID, articleId)
	}

	if article == nil {
		return nil, errs.ArticleNotFound
	}

	// cancel early...
	if ctx.IsClient() && !article.ClientAccess {
		// don't return nil in case access is denied
		return article, errs.PermissionDenied
	}

	access := model.FindArticleAccessByArticleIdAndUserIdTx(tx, articleId, ctx.UserId)

	if !article.ReadEveryone && access == nil {
		// don't return nil in case access is denied
		return article, errs.PermissionDenied
	}

	article.Access = access
	return article, nil
}

// CheckCommitMsg checks the length of a commit message.
func CheckCommitMsg(commit string) error {
	if len(commit) > maxCommitMsgLen {
		return errs.CommitMsgLen
	}

	return nil
}

// CheckContentVersionRequiresExpert checks if the content version requires an expert organization to be read/modified or not.
func CheckContentVersionRequiresExpert(isExpert bool, contentVersion, lastContentVersion int) error {
	if !isExpert && contentVersion <= lastContentVersion-ArticleMaxHistoryNonExpert {
		return errs.RequiresExpertVersion
	}

	return nil
}
