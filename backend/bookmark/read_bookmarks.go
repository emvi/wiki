package bookmark

import (
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
)

const (
	maxArticles = 10
	maxLists    = 10
)

func ReadBookmarks(orga *model.Organization, userId hide.ID, readArticles, readLists bool, offsetArticles, offsetLists int) ([]model.Bookmark, []model.Bookmark) {
	langId := util.DetermineLang(nil, orga.ID, userId, 0).ID
	var articles []model.Bookmark
	var lists []model.Bookmark

	if readArticles {
		articles = model.FindBookmarkByOrganizationIdAndUserIdAndLanguageIdArticleIdSetWithLimit(orga.ID, userId, langId, offsetArticles, maxArticles)
	}

	if readLists {
		lists = model.FindBookmarkByOrganizationIdAndUserIdAndLanguageIdArticleListIdSetWithLimit(orga.ID, userId, langId, offsetLists, maxLists)
	}

	return articles, lists
}
