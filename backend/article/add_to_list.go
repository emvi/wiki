package article

import (
	"emviwiki/backend/articlelist"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func AddArticleToLists(orga *model.Organization, userId, articleId hide.ID, listIds []hide.ID) error {
	if _, err := checkUserReadAccess(orga.ID, userId, articleId); err != nil {
		return err
	}

	for _, listId := range listIds {
		if _, err := articlelist.AddArticleListEntry(orga, userId, listId, []hide.ID{articleId}); err != nil {
			logbuch.Warn("Error adding article to article list", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "article_id": articleId, "list_id": listId})
			return err
		}
	}

	return nil
}
