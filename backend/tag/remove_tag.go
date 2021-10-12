package tag

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func RemoveTag(organization *model.Organization, userId, articleId hide.ID, tag string) error {
	article := model.GetArticleByOrganizationIdAndIdIgnoreArchived(organization.ID, articleId)

	if article == nil {
		return errs.ArticleNotFound
	}

	access := model.FindArticleAccessByArticleIdAndUserId(articleId, userId)

	if !article.ReadEveryone && access == nil {
		return errs.PermissionDenied
	}

	articleTag := model.GetArticleTagByOrganizationIdAndArticleIdAndName(organization.ID, article.ID, tag)

	if articleTag == nil {
		return nil
	}

	if err := perm.CheckUserTagAccess(organization.ID, userId, articleTag.TagId); err != nil {
		return err
	}

	if err := model.DeleteArticleTagByOrganizationIdAndId(nil, organization.ID, articleTag.ID); err != nil {
		logbuch.Error("Error deleting article tag when removing tag from article", logbuch.Fields{"err": err, "article_tag_id": articleTag.ID})
		return errs.Saving
	}

	go CleanupUnusedTags(organization.ID)
	return nil
}

func CleanupUnusedTags(orgaId hide.ID) {
	if err := model.DeleteTagUnusedByOrganizationId(nil, orgaId); err != nil {
		logbuch.Error("Error deleting unused tags", logbuch.Fields{"err": err, "orga_id": orgaId})
	}
}
