package observe

import (
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
)

const (
	maxArticles = 10
	maxLists    = 10
	maxGroups   = 10
)

func ReadObserved(orga *model.Organization, userId hide.ID, readArticles, readLists, readGroups bool, offsetArticles, offsetLists, offsetGroups int) ([]model.Article, []model.ArticleList, []model.UserGroup) {
	langId := util.DetermineLang(nil, orga.ID, userId, 0).ID
	var articles []model.Article
	var lists []model.ArticleList
	var groups []model.UserGroup

	if readArticles {
		articles = model.FindArticleByOrganizationIdAndUserIdAndLanguageIdAndObservedWithLimit(orga.ID, userId, langId, offsetArticles, maxArticles)
	}

	if readLists {
		lists = model.FindArticleListByOrganizationIdAndUserIdAndLanguageIdAndObservedWithLimit(orga.ID, userId, langId, offsetLists, maxLists)
	}

	if readGroups {
		groups = model.FindUserGroupByOrganizationIdAndUserIdAndObservedWithLimit(orga.ID, userId, offsetGroups, maxGroups)
	}

	return articles, lists, groups
}
