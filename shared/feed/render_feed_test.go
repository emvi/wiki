package feed

import (
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"github.com/emvi/null"
	"testing"
)

func TestRenderFeedVariable(t *testing.T) {
	orga := &model.Organization{NameNormalized: "test"}
	text := "{{.Feed.Reason}}"
	feed := &model.Feed{Reason: "feed_reason"}
	expected := "feed_reason"
	is := RenderFeed(orga, text, FeedText, "en", feed)

	if is != expected {
		t.Fatalf("Expected rendered feed text '%v' to be '%v'", is, expected)
	}
}

func TestRenderFeedHashID(t *testing.T) {
	orga := &model.Organization{NameNormalized: "test"}
	text := "{{.Feed.ID | IdToString}}"
	feed := &model.Feed{BaseEntity: db.BaseEntity{ID: 123}, Reason: "hash_id"}
	expected := "beJarVNaQM"
	is := RenderFeed(orga, text, FeedText, "en", feed)

	if is != expected {
		t.Fatalf("Expected rendered feed text '%v' to be '%v'", is, expected)
	}
}

func TestRenderFeedSlices(t *testing.T) {
	orga := &model.Organization{NameNormalized: "test"}
	text := "{{range $i, $e := .Articles}}{{if $i}}, {{end}}{{$e.LatestArticleContent.Title}}{{end}}"
	article1 := &model.Article{LatestArticleContent: &model.ArticleContent{Title: "Title 1"}}
	article2 := &model.Article{LatestArticleContent: &model.ArticleContent{Title: "Title 2"}}
	refs := []model.FeedRef{
		{ArticleID: 123, Article: article1},
		{ArticleID: 321, Article: article2},
	}
	feed := &model.Feed{
		Reason:   "slices",
		FeedRefs: refs,
	}
	expected := "Title 1, Title 2"
	is := RenderFeed(orga, text, FeedText, "en", feed)

	if is != expected {
		t.Fatalf("Expected rendered feed text '%v' to be '%v'", is, expected)
	}
}

func TestRenderFeedTwice(t *testing.T) {
	orga := &model.Organization{NameNormalized: "test"}
	text := "this is a test"
	feed := &model.Feed{BaseEntity: db.BaseEntity{ID: 123}, Reason: "twice"}
	expected := "this is a test"
	is := RenderFeed(orga, text, FeedText, "en", feed)

	if is != expected {
		t.Fatalf("Expected rendered feed text '%v' to be '%v'", is, expected)
	}

	is = RenderFeed(orga, text, FeedText, "en", feed)

	if is != expected {
		t.Fatalf("Expected rendered feed text '%v' to be '%v'", is, expected)
	}
}

func TestRenderFeedKeyValue(t *testing.T) {
	orga := &model.Organization{NameNormalized: "test"}
	text := `{{index .Vars "notfound"}} {{index .Vars "key"}}`
	refs := []model.FeedRef{{Key: null.NewString("key", true), Value: null.NewString("value", true)}}
	feed := &model.Feed{
		Reason:   "keyvalue",
		FeedRefs: refs,
	}
	expected := " value"
	is := RenderFeed(orga, text, FeedText, "en", feed)

	if is != expected {
		t.Fatalf("Expected rendered feed text '%v' to be '%v'", is, expected)
	}
}

func TestRenderFeedReasons(t *testing.T) {
	orga := &model.Organization{NameNormalized: "test"}
	user := &model.User{BaseEntity: db.BaseEntity{ID: 1}, OrganizationMember: &model.OrganizationMember{}}
	group := &model.UserGroup{BaseEntity: db.BaseEntity{ID: 2}}
	article := &model.Article{BaseEntity: db.BaseEntity{ID: 3}, LatestArticleContent: &model.ArticleContent{}}
	content := &model.ArticleContent{BaseEntity: db.BaseEntity{ID: 4}}
	list := &model.ArticleList{BaseEntity: db.BaseEntity{ID: 5}, Name: &model.ArticleListName{}}
	refs := []model.FeedRef{
		{UserID: user.ID, User: user},
		{UserGroupID: group.ID, UserGroup: group},
		{ArticleID: article.ID, Article: article},
		{ArticleContentID: content.ID, ArticleContent: content},
		{ArticleListID: list.ID, ArticleList: list},
	}
	feed := &model.Feed{FeedRefs: refs}
	textTypes := []string{FeedText, NotificationText}

	for _, textType := range textTypes {
		for langCode, reasons := range Reasons {
			for reason, text := range reasons {
				feed.Reason = reason
				var feedText string

				if textType == FeedText {
					feedText = text.Feed
				} else {
					feedText = text.Notification
				}

				if feedText == "" {
					continue
				}

				if RenderFeed(orga, feedText, textType, langCode, feed) == "" {
					t.Fatalf("Feed reason '%v' must have been rendered", reason)
				}
			}
		}
	}
}
