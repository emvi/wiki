package article

import (
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
)

const (
	maxPrivateArticles = 10
	maxDrafts          = 20
)

func ReadPrivateArticles(orga *model.Organization, userId hide.ID, offset int) []model.Article {
	langId := util.DetermineLang(nil, orga.ID, userId, 0).ID
	return model.FindArticleByOrganizationIdAndUserIdAndLanguageIdAndPrivateWithLimit(orga.ID, userId, langId, offset, maxPrivateArticles)
}

func ReadDrafts(orga *model.Organization, userId hide.ID, offset int) []model.Article {
	langId := util.DetermineLang(nil, orga.ID, userId, 0).ID
	return model.FindArticleByOrganizationIdAndUserIdAndLanguageIdAndUnpublishedWithLimit(orga.ID, userId, langId, offset, maxDrafts)
}
