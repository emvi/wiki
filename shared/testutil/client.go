package testutil

import (
	"emviwiki/shared/model"
	"testing"
)

func SetArticleClientAccess(t *testing.T, article *model.Article) {
	article.ClientAccess = true

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}
}

func SetListClientAccess(t *testing.T, list *model.ArticleList) {
	list.ClientAccess = true

	if err := model.SaveArticleList(nil, list); err != nil {
		t.Fatal(err)
	}
}
