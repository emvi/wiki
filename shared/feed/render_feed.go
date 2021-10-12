package feed

import (
	"bytes"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"fmt"
	"github.com/emvi/logbuch"
	"sync"
	"text/template"
)

const (
	FeedText         = "f_"
	NotificationText = "n_"
)

var (
	feedTemplates      = make(map[string]*template.Template)
	feedTemplatesMutex sync.RWMutex
)

// RenderFeed renders given text and returns it as string.
// The type is used to distinguish feed and notification texts.
func RenderFeed(orga *model.Organization, text, textType, langCode string, feed *model.Feed) string {
	tpl := getFeedTpl(fmt.Sprintf("%s%s_%s", textType, langCode, feed.Reason), text)

	if tpl == nil {
		return ""
	}

	var buffer bytes.Buffer
	user, groups, articles, content, lists, kvs := extractFeedRefLists(feed)
	data := struct {
		Feed         *model.Feed
		User         []model.User
		Groups       []model.UserGroup
		Articles     []model.Article
		Content      []model.ArticleContent
		Lists        []model.ArticleList
		Vars         map[string]string
		WebsiteHost  string
		AuthHost     string
		FrontendHost string
	}{
		feed,
		user,
		groups,
		articles,
		content,
		lists,
		kvs,
		websiteHost,
		authHost,
		util.InjectSubdomain(frontendHost, orga.NameNormalized),
	}

	if err := tpl.Execute(&buffer, data); err != nil {
		logbuch.Debug("Error executing feed reason template", logbuch.Fields{"err": err, "reason": feed.Reason})
		return ""
	}

	return buffer.String()
}

func getFeedTpl(reason, text string) *template.Template {
	feedTemplatesMutex.RLock()
	tpl := feedTemplates[reason]
	feedTemplatesMutex.RUnlock()

	if tpl == nil {
		feedTemplatesMutex.Lock()
		defer feedTemplatesMutex.Unlock()
		var err error
		tpl, err = template.New("feed").Funcs(renderFuncs).Parse(text)

		if err != nil {
			logbuch.Debug("Error parsing feed reason template", logbuch.Fields{"err": err})
			return nil
		}

		feedTemplates[reason] = tpl
	}

	return tpl
}

func extractFeedRefLists(feed *model.Feed) ([]model.User, []model.UserGroup, []model.Article, []model.ArticleContent, []model.ArticleList, map[string]string) {
	var user []model.User
	var groups []model.UserGroup
	var articles []model.Article
	var content []model.ArticleContent
	var lists []model.ArticleList
	kvs := make(map[string]string)

	for _, ref := range feed.FeedRefs {
		if ref.UserID != 0 {
			user = append(user, *ref.User)
		} else if ref.UserGroupID != 0 {
			groups = append(groups, *ref.UserGroup)
		} else if ref.ArticleID != 0 {
			articles = append(articles, *ref.Article)
		} else if ref.ArticleContentID != 0 {
			content = append(content, *ref.ArticleContent)
		} else if ref.ArticleListID != 0 {
			lists = append(lists, *ref.ArticleList)
		} else if ref.Key.Valid {
			kvs[ref.Key.String] = ref.Value.String
		}
	}

	return user, groups, articles, content, lists, kvs
}
