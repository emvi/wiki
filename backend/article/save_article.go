package article

import (
	articleutil "emviwiki/backend/article/util"
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/backend/observe"
	"emviwiki/backend/perm"
	"emviwiki/backend/pinned"
	"emviwiki/backend/prosemirror"
	"emviwiki/backend/tag"
	"emviwiki/shared/constants"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxTitleLen = 100
	maxArticles = 100
)

type SaveArticleData struct {
	Organization  *model.Organization
	UserId        hide.ID                  `json:"user_id"`
	Id            hide.ID                  `json:"id"`
	RoomId        string                   `json:"room_id"`
	LanguageId    hide.ID                  `json:"language_id"`
	Authors       []hide.ID                `json:"authors"`
	CommitMsg     string                   `json:"message"`
	Wip           bool                     `json:"wip"`
	ReadEveryone  bool                     `json:"read_everyone"`
	WriteEveryone bool                     `json:"write_everyone"`
	Private       bool                     `json:"private"`
	ClientAccess  bool                     `json:"client_access"`
	Access        []perm.SaveArticleAccess `json:"access"`
	Title         string                   `json:"title"`
	Content       string                   `json:"content"`
	RTL           bool                     `json:"rtl"`
	Tags          []string                 `json:"tags"`
}

func (data *SaveArticleData) validate() []error {
	data.RoomId = strings.TrimSpace(data.RoomId)
	data.Title = strings.TrimSpace(data.Title)
	data.Content = strings.TrimSpace(data.Content)
	data.CommitMsg = strings.TrimSpace(data.CommitMsg)
	err := make([]error, 0)

	if data.Title == "" {
		err = append(err, errs.NoTitle)
	} else if utf8.RuneCountInString(data.Title) > maxTitleLen {
		err = append(err, errs.TitleLen)
	}

	if e := articleutil.CheckCommitMsg(data.CommitMsg); e != nil {
		err = append(err, e)
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func SaveArticle(data SaveArticleData) (hide.ID, []error) {
	savingStartTime := time.Now()
	logbuch.Debug("Started saving article", logbuch.Fields{"data": data})

	if !data.Organization.Expert {
		if model.CountArticleByOrganizationId(data.Organization.ID) >= maxArticles {
			return 0, []error{errs.MaxArticlesReached}
		}
	}

	data.LanguageId = util.DetermineLang(nil, data.Organization.ID, data.UserId, data.LanguageId).ID
	data.Tags = cleanTags(data.Tags)
	article, err := getArticleOrCreateNew(&data)

	if err != nil {
		return 0, []error{err}
	}

	if model.GetLanguageByOrganizationIdAndId(data.Organization.ID, data.LanguageId) == nil {
		return 0, []error{errs.LanguageNotFound}
	}

	// validate write access and authors
	if data.Id != 0 && !article.WriteEveryone && !perm.CheckUserWriteAccess(data.Id, data.UserId) {
		return 0, []error{errs.PermissionDenied}
	}

	validateAuthors(&data)

	if err := data.validate(); err != nil {
		return 0, err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to save article", logbuch.Fields{"err": err})
		return 0, []error{errs.TxBegin}
	}

	if !data.Wip {
		if err := deleteWIP(tx, article.ID); err != nil {
			return 0, []error{err}
		}
	}

	lastCommit := model.GetArticleContentLastByArticleIdAndLanguageIdAndWIPTx(tx, article.ID, data.LanguageId, true)

	if err := saveArticle(tx, data.Organization, article, &data, lastCommit); err != nil {
		return 0, []error{err}
	}

	logbuch.Debug("Saving content", logbuch.Fields{"id": article.ID})
	content, err := saveContent(tx, article.ID, &data, lastCommit)

	if err != nil {
		return 0, []error{err}
	}

	if err := saveAccess(tx, article.ID, &data); err != nil {
		return 0, []error{err}
	}

	if err := updateAttachments(tx, data.Organization.ID, article.ID, content.LanguageId, data.RoomId); err != nil {
		return 0, []error{err}
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when saving article", logbuch.Fields{"err": err})
		return 0, []error{errs.TxCommit}
	}

	// IGNORE ERRORS FOR SUBSEQUENT ACTIONS
	// since these steps are not essential to save an article
	subsequentActions(&data, article, content, lastCommit)

	savingTime := time.Now().Sub(savingStartTime)
	logbuch.Debug("Article saved", logbuch.Fields{"id": article.ID, "time": savingTime})
	return article.ID, nil
}

func subsequentActions(data *SaveArticleData, article *model.Article, content, lastCommit *model.ArticleContent) {
	saveTags(data.Organization, article.ID, data.Tags)
	removePinnedWhenPrivate(data.Organization, article, data)

	if !data.Wip {
		createSaveArticleFeed(data, article, content)
	}

	// observe if this is a new article
	if data.Id == 0 {
		observeArticle(article.ID, data.Authors)
	}

	// cleanup attachments that were uploaded but are not used in new content anymore
	cleanupDefTime := time.Time{}

	if lastCommit != nil {
		cleanupDefTime = lastCommit.DefTime
	}

	go func() {
		if err := cleanupAttachments(data.Organization.ID, data.UserId, article.ID, cleanupDefTime, content); err != nil {
			logbuch.Error("Error cleaning up attachments when saving article", logbuch.Fields{"err": err, "article_id": article.ID, "content_id": content.ID})
		}
	}()

	// notify mentioned users
	go func() {
		if err := notifyMentionedUsers(data.Organization, data.UserId, cleanupDefTime, article, content); err != nil {
			logbuch.Error("Error notifying mentioned users when saving article", logbuch.Fields{"err": err, "article_id": article.ID, "content_id": content.ID})
		}
	}()
}

func cleanTags(tags []string) []string {
	clean := make([]string, 0)

	for _, t := range tags {
		if strings.TrimSpace(t) != "" {
			clean = append(clean, t)
		}
	}

	return clean
}

func getArticleOrCreateNew(data *SaveArticleData) (*model.Article, error) {
	article := &model.Article{OrganizationId: data.Organization.ID, WIP: -1}

	if data.Id != 0 {
		article = model.GetArticleByOrganizationIdAndId(data.Organization.ID, data.Id)

		if article == nil {
			logbuch.Error("Error finding article to save", logbuch.Fields{"article_id": data.Id, "organization": data.Organization.ID, "user_id": data.UserId})
			return nil, errs.ArticleNotFound
		}
	}

	return article, nil
}

// Validates authors exist.
// Also ensures the user saving the article is within the list.
func validateAuthors(data *SaveArticleData) {
	authors := make([]hide.ID, 0)

	for _, author := range data.Authors {
		if author != data.UserId &&
			model.GetUserByOrganizationIdAndId(data.Organization.ID, author) != nil {
			authors = append(authors, author)
		}
	}

	authors = append(authors, data.UserId)
	data.Authors = authors
}

func saveArticle(tx *sqlx.Tx, orga *model.Organization, article *model.Article, data *SaveArticleData, lastCommit *model.ArticleContent) error {
	wipVersion := article.WIP

	// set WIP version if saved as WIP and not set yet
	if data.Wip && wipVersion == -1 {
		if lastCommit != nil {
			wipVersion = lastCommit.Version
		} else {
			wipVersion = 1
		}
	} else if !data.Wip {
		wipVersion = -1
	}

	// tags are added elsewhere, using the tag module
	article.WIP = wipVersion
	article.ReadEveryone = data.ReadEveryone || data.WriteEveryone || data.ClientAccess
	article.WriteEveryone = data.WriteEveryone
	article.Private = data.Private && !data.ReadEveryone && !data.WriteEveryone && !data.ClientAccess
	article.ClientAccess = data.ClientAccess && orga.Expert

	// set published date when article gets published for the first time not WIP
	if article.WIP == -1 && !article.Published.Valid {
		article.Published.SetValid(time.Now())
	}

	if err := model.SaveArticle(tx, article); err != nil {
		logbuch.Error("Error saving article", logbuch.Fields{"err": err, "article_id": article.ID, "user_id": data.UserId})
		return errs.Saving
	}

	return nil
}

func saveContent(tx *sqlx.Tx, articleId hide.ID, data *SaveArticleData, lastCommit *model.ArticleContent) (*model.ArticleContent, error) {
	doc, err := prosemirror.ParseDoc(data.Content)

	if err != nil {
		db.Rollback(tx)
		logbuch.Error("Error parsing content while saving article", logbuch.Fields{"err": err})
		return nil, errs.Saving
	}

	// create new commit based on last commit
	version := 1

	if lastCommit != nil {
		version = lastCommit.Version + 1
	}

	textContent := extractTextFromContent(doc)
	newContent := &model.ArticleContent{ArticleId: articleId,
		LanguageId:      data.LanguageId,
		UserId:          data.UserId,
		Title:           data.Title,
		Content:         data.Content,
		Version:         version,
		Commit:          null.NewString(data.CommitMsg, data.CommitMsg != ""),
		WIP:             data.Wip,
		TitleTsvector:   data.Title,
		ContentTsvector: textContent,
		ReadingTime:     calculateReadingTimeSeconds(textContent),
		SchemaVersion:   constants.LatestSchemaVersion,
		RTL:             data.RTL}

	if err := model.SaveArticleContent(tx, newContent); err != nil {
		logbuch.Error("Error saving new article content when saving article", logbuch.Fields{"err": err, "article_id": articleId})
		return nil, errs.Saving
	}

	if err := saveAuthors(tx, newContent.ID, data.Authors); err != nil {
		logbuch.Error("Error saving authors for new article content when saving article", logbuch.Fields{"err": err, "article_id": articleId})
		return nil, err
	}

	// don't update latest content if existing article and WIP
	if data.Wip && lastCommit != nil {
		return newContent, nil
	}

	// set empty content if this is a new article and WIP
	if data.Wip && lastCommit == nil {
		data.Content = ""
	}

	// find latest commit
	// if it does not exist, this is either a new article or the content for the language does not exist yet
	// if it does for the given language, update it
	latestContent := model.GetArticleContentLatestByArticleIdAndLanguageIdTx(tx, articleId, data.LanguageId, false)

	if latestContent != nil {
		// just set it here, the current content is build when the article is opened for editing
		latestContent.Title = data.Title
		latestContent.Content = data.Content
		latestContent.WIP = data.Wip
		latestContent.TitleTsvector = data.Title
		latestContent.ContentTsvector = textContent
		latestContent.ReadingTime = calculateReadingTimeSeconds(textContent)
		latestContent.RTL = data.RTL
	} else {
		latestContent = &model.ArticleContent{ArticleId: articleId,
			LanguageId:      data.LanguageId,
			UserId:          data.UserId,
			Title:           data.Title,
			Content:         data.Content,
			WIP:             data.Wip,
			Version:         0, // latest is always marked as 0
			TitleTsvector:   data.Title,
			ContentTsvector: textContent,
			ReadingTime:     calculateReadingTimeSeconds(textContent),
			SchemaVersion:   constants.LatestSchemaVersion,
			RTL:             data.RTL}
	}

	if err := model.SaveArticleContent(tx, latestContent); err != nil {
		logbuch.Error("Error saving latest article content when saving article", logbuch.Fields{"err": err, "article_id": articleId})
		return nil, errs.Saving
	}

	return newContent, nil
}

func saveAuthors(tx *sqlx.Tx, articleContentId hide.ID, authors []hide.ID) error {
	for _, author := range authors {
		contentAuthor := &model.ArticleContentAuthor{ArticleContentId: articleContentId, UserId: author}

		if err := model.SaveArticleContentAuthor(tx, contentAuthor); err != nil {
			logbuch.Error("Error saving article content author when saving article", logbuch.Fields{"err": err, "author": author, "article_content_id": articleContentId})
			return errs.Saving
		}
	}

	return nil
}

func saveAccess(tx *sqlx.Tx, articleId hide.ID, data *SaveArticleData) error {
	// make sure all current authors have write access, except this article is private
	if data.Private {
		data.Access = make([]perm.SaveArticleAccess, 1)
		data.Access[0] = perm.SaveArticleAccess{UserId: data.UserId, Write: true}
	} else {
		for _, author := range data.Authors {
			data.Access = append(data.Access, perm.SaveArticleAccess{UserId: author, Write: true})
		}
	}

	// filter list and check all users/groups exist
	access := perm.FilterAccessList(data.Access)

	if err := perm.CheckAccessList(tx, data.Organization, access); err != nil {
		db.Rollback(tx)
		return err
	}

	// save new permission
	if err := perm.DeleteAccess(tx, articleId); err != nil {
		return err
	}

	if err := perm.SaveAccessList(tx, data.Organization, access, articleId); err != nil {
		return err
	}

	return nil
}

func deleteWIP(tx *sqlx.Tx, articleId hide.ID) error {
	logbuch.Debug("Deleting WIP saves", logbuch.Fields{"id": articleId})
	wipContentIds := model.FindArticleContentIdByArticleIdAndWIPTx(tx, articleId)

	for _, id := range wipContentIds {
		if err := feed.DeleteFeed(tx, &feed.DeleteFeedData{ArticleContentId: id}); err != nil {
			return errs.Saving
		}
	}

	if err := model.DeleteArticleContentByArticleIdAndWIP(tx, articleId); err != nil {
		return errs.Saving
	}

	return nil
}

func saveTags(orga *model.Organization, articleId hide.ID, tags []string) {
	if len(tags) > tag.MaxTagsPerArticle {
		tags = tags[:tag.MaxTagsPerArticle]
	}

	for _, t := range tags {
		if err := tag.AddTag(orga, tag.AddTagData{ArticleId: articleId, Tag: t}); err != nil && err != errs.TagExistsAlready {
			logbuch.Warn("Error adding tag when saving", logbuch.Fields{"err": err, "article_id": articleId, "tag": t})
		}
	}
}

func removePinnedWhenPrivate(orga *model.Organization, article *model.Article, data *SaveArticleData) {
	if article.Private && pinned.IsPinned(orga.ID, article.ID, 0) {
		if err := pinned.PinObject(data.Organization, data.UserId, article.ID, 0); err != nil {
			logbuch.Error("Error pinning article when saving article", logbuch.Fields{"err": err, "article_id": article.ID})
		}
	}
}

// don't use a transaction here to make sure saving does not fail due to it!
func createSaveArticleFeed(data *SaveArticleData, article *model.Article, content *model.ArticleContent) {
	reason := "create_article"

	if data.Id != 0 {
		reason = "update_article"
	}

	refs := make([]interface{}, 2)
	refs[0] = article
	refs[1] = content
	feedData := &feed.CreateFeedData{Organization: data.Organization,
		UserId: data.UserId,
		Reason: reason,
		Public: data.ReadEveryone || data.WriteEveryone,
		Access: perm.GetUserIdsFromArticleAccess(data.Access),
		Notify: model.FindObservedObjectUserIdByArticleIdOrArticleListId(article.ID, 0),
		Refs:   refs}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when creating/updating article", logbuch.Fields{"err": err})
	}
}

func observeArticle(articleId hide.ID, authors []hide.ID) {
	for _, author := range authors {
		if !observe.IsObserved(author, articleId, 0, 0) {
			observed := &model.ObservedObject{UserId: author,
				ArticleId: articleId}

			if err := model.SaveObservedObject(nil, observed); err != nil {
				logbuch.Error("Error author observing article when saving article", logbuch.Fields{"err": err})
			}
		}
	}
}
