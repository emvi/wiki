package article

import (
	"emviwiki/backend/article/schema"
	articleutil "emviwiki/backend/article/util"
	"emviwiki/backend/bookmark"
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/backend/observe"
	"emviwiki/backend/perm"
	"emviwiki/backend/prosemirror"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"strings"
	"time"
)

const (
	mentionArticleType = "article"
	mentionGroupType   = "group"
	mentionListType    = "list"
	mentionTagType     = "tag"
	mentionTitleAttr   = "title"
	mentionNotFoundEN  = "[Not found]"
	mentionNotFoundDE  = "[Nicht gefunden]"
	mentionNoAccessEN  = "[No access]"
	mentionNoAccessDE  = "[Kein Zugriff]"
)

type ArticleResult struct {
	Article         *model.Article
	Content         *model.ArticleContent
	Authors         []model.User
	WriteAccess     bool
	IsObserved      bool
	IsBookmarked    bool
	Recommendations []model.ArticleRecommendation
}

// ReadArticle reads an article and renders its content if so desired.
// In addition to the article, all relevant meta data is returned.
// The format can be either HTML or Markdown.
func ReadArticle(ctx context.EmviContext, articleId, langId hide.ID, version int, renderContent bool, format string) (ArticleResult, error) {
	article, err := articleutil.GetArticleWithAccess(nil, ctx, articleId, true)

	if err != nil {
		return ArticleResult{}, err
	}

	content := GetArticleContent(ctx.Organization.ID, ctx.UserId, articleId, langId, version)

	if content == nil {
		return ArticleResult{}, errs.ArticleContentVersionNotFound
	}

	if err := checkVersionAllowed(ctx.Organization, content, version); err != nil {
		return ArticleResult{}, err
	}

	if ctx.IsClient() {
		articleutil.RemoveNonPublicInformationFromContent(ctx, content)
	}

	if !ctx.IsClient() || ctx.HasScopes(client.Scopes["tags"]) {
		article.Tags = model.FindTagByOrganizationIdAndUserIdAndArticleId(ctx.Organization.ID, ctx.UserId, articleId)
	}

	isObserved := false
	isBookmarked := false
	writeAccess := false

	if ctx.IsUser() {
		article.Access = model.FindArticleAccessByOrganizationIdAndArticleId(ctx.Organization.ID, article.ID)
		isObserved = observe.IsObserved(ctx.UserId, article.ID, 0, 0)
		isBookmarked = bookmark.IsBookmarked(ctx.UserId, articleId, 0)
		writeAccess = hasWriteAccess(article, ctx.UserId)
	}

	if (renderContent || ctx.IsClient()) && content.Content != "" {
		format = strings.ToLower(format)

		if format == formatMarkdown {
			content.Content, err = renderArticleContent(ctx, ctx.Organization.ID, ctx.UserId, content, schema.GetMarkdownSchema(ctx.Organization))
		} else {
			content.Content, err = renderArticleContent(ctx, ctx.Organization.ID, ctx.UserId, content, schema.HTMLSchema)
		}

		if err != nil {
			// return technical error here because this is an "exception" state
			return ArticleResult{}, err
		}
	}

	if ctx.IsUser() {
		updateArticleViews(article, ctx.UserId)
	}

	return ArticleResult{
		article,
		content,
		getAuthors(ctx, articleId),
		writeAccess,
		isObserved,
		isBookmarked,
		getRecommendations(articleId, ctx.UserId),
	}, nil
}

func checkVersionAllowed(orga *model.Organization, content *model.ArticleContent, version int) error {
	if version == 0 || orga.Expert {
		return nil
	}

	lastContent := model.GetArticleContentLastByArticleIdAndLanguageIdAndWIP(content.ArticleId, content.LanguageId, false)

	if lastContent != nil {
		if err := articleutil.CheckContentVersionRequiresExpert(orga.Expert, content.Version, lastContent.Version); err != nil {
			return err
		}
	}

	return nil
}

// GetArticleContent returns the article content for given article, language and version with a best effort principle.
// In case the content cannot be found, an empty content will be returned.
// The content is migrated to the latest schema version if required.
func GetArticleContent(orgaId, userId, articleId, langId hide.ID, version int) *model.ArticleContent {
	var content *model.ArticleContent

	if version == 0 {
		if langId == 0 {
			langId = util.DetermineLang(nil, orgaId, userId, 0).ID
			content = model.GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageId(orgaId, articleId, langId, true)
		} else {
			content = model.GetArticleContentLatestByArticleIdAndLanguageId(articleId, langId, true)
		}
	} else {
		content = model.GetArticleContentByArticleIdAndLanguageIdAndMaxVersion(articleId, langId, version)
	}

	if content == nil {
		content = &model.ArticleContent{LanguageId: langId, SchemaVersion: constants.LatestSchemaVersion}
	}

	if err := schema.Migrate(content); err != nil {
		// there is not much we can do about it...
		logbuch.Fatal("Error migration article content", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "article_id": articleId, "lang_id": langId, "version": version, "content_id": content.ID})
	}

	return content
}

func getAuthors(ctx context.EmviContext, articleId hide.ID) []model.User {
	var authors []model.User

	if !ctx.IsClient() || ctx.HasScopes(client.Scopes["article_authors"]) {
		authors = model.FindArticleContentAuthorUserByArticleId(articleId)
		removeEmail := ctx.IsClient() && !ctx.HasScopes(client.Scopes["article_authors_mails"])

		for i := range authors {
			authors[i].AcceptMarketing = false

			if ctx.IsClient() {
				authors[i].OrganizationMember = nil
			}

			if removeEmail {
				authors[i].Email = ""
			}
		}
	}

	return authors
}

func getRecommendations(articleId, userId hide.ID) []model.ArticleRecommendation {
	return model.FindArticleRecommendationByArticleIdAndRecommendedToWithUser(articleId, userId)
}

func hasWriteAccess(article *model.Article, userId hide.ID) bool {
	if article.WriteEveryone {
		return true
	}

	return perm.CheckUserWriteAccess(article.ID, userId)
}

func updateArticleViews(article *model.Article, userId hide.ID) {
	if model.GetArticleVisitByArticleIdAndUserIdAndDefTimeAfter(article.ID, userId, time.Now().Add(-time.Hour*24)) == nil {
		visit := &model.ArticleVisit{ArticleId: article.ID, UserId: userId}

		if err := model.SaveArticleVisit(nil, visit); err != nil {
			logbuch.Error("Error saving article visit", logbuch.Fields{"err": err, "article_id": article.ID, "user_id": userId})
		}

		article.Views++

		if err := model.SaveArticle(nil, article); err != nil {
			logbuch.Error("Error saving article when updating views", logbuch.Fields{"err": err, "article_id": article.ID, "user_id": userId})
		}
	}
}

func renderArticleContent(ctx context.EmviContext, orgaId, userId hide.ID, content *model.ArticleContent, schema *prosemirror.Schema) (string, error) {
	if content.Content == "" {
		return "", nil
	}

	doc, err := prosemirror.ParseDoc(content.Content)

	if err != nil {
		logbuch.Warn("Error parsing article content to prosemirror document", logbuch.Fields{"err": err, "article_content_id": content.ID})
		return "", err
	}

	return RenderDocument(ctx, orgaId, userId, content.LanguageId, doc, schema)
}

// RenderDocument renders given node for given schema.
func RenderDocument(ctx context.EmviContext, orgaId, userId, langId hide.ID, doc *prosemirror.Node, schema *prosemirror.Schema) (string, error) {
	findAndReplaceMentions(ctx, doc, orgaId, userId, langId)
	out, err := prosemirror.RenderDoc(schema, doc)

	if err != nil {
		logbuch.Warn("Error rendering article content", logbuch.Fields{"err": err})
		return "", err
	}

	return out, nil
}

func findAndReplaceMentions(ctx context.EmviContext, doc *prosemirror.Node, orgaId, userId, langId hide.ID) {
	langCode := ""
	lang := model.GetLanguageByOrganizationIdAndId(orgaId, langId)

	if lang != nil {
		langCode = lang.Code
	}

	isClient := ctx.IsClient()
	clientArticleAccess := ctx.HasScopes(client.Scopes["articles"])
	clientListAccess := ctx.HasScopes(client.Scopes["lists"])
	clientTagAccess := ctx.HasScopes(client.Scopes["tags"])
	noAccessText := getNoAccessText(langCode)
	notFoundText := getNotFoundText(langCode)

	prosemirror.TransformNodes(doc, mentionTypeName, func(node *prosemirror.Node) {
		mentionType, _, mentionId := getMentionAttrs(*node)
		title := ""
		access := false

		if !isClient || checkClientMentionAccess(mentionType, clientArticleAccess, clientListAccess, clientTagAccess) {
			title, access = getMentionTitleAndAccessFromObject(isClient, orgaId, userId, langId, mentionType, mentionId)
		}

		if !access {
			title = noAccessText
		} else if title == "" {
			title = notFoundText
		}

		node.Attrs[mentionTitleAttr] = title
	})
}

func getNoAccessText(langCode string) string {
	if langCode == "de" {
		return mentionNoAccessDE
	}

	return mentionNoAccessEN
}

func getNotFoundText(langCode string) string {
	if langCode == "de" {
		return mentionNotFoundDE
	}

	return mentionNotFoundEN
}

func checkClientMentionAccess(mentionType string, clientArticleAccess, clientListAccess, clientTagAccess bool) bool {
	return mentionType == mentionArticleType && clientArticleAccess ||
		mentionType == mentionListType && clientListAccess ||
		mentionType == mentionTagType && clientTagAccess
}

func getMentionTitleAndAccessFromObject(isClient bool, orgaId, userId, langId hide.ID, mentionType, mentionId string) (string, bool) {
	title := ""
	access := true

	if mentionType == mentionArticleType {
		id := getMentionIdFromString(mentionId, mentionType)

		if id != 0 {
			content := model.GetArticleContentLatestByArticleIdAndLanguageId(id, langId, false)

			if content != nil {
				title = content.Title

				if !isClient {
					_, err := checkUserReadAccess(orgaId, userId, id)
					access = err == nil
				} else {
					article := model.GetArticleByOrganizationIdAndId(orgaId, id)
					access = article.ClientAccess
				}
			}
		}
	} else if mentionType == mentionGroupType {
		id := getMentionIdFromString(mentionId, mentionType)

		if id != 0 {
			group := model.GetUserGroupByOrganizationIdAndId(orgaId, id)

			if group != nil {
				title = group.Name
			}
		}
	} else if mentionType == mentionListType {
		id := getMentionIdFromString(mentionId, mentionType)

		if id != 0 {
			name := model.GetArticleListNameByOrganizationIdAndArticleListIdAndLangId(orgaId, id, langId)

			if name != nil {
				title = name.Name
				list := model.GetArticleListByOrganizationIdAndId(orgaId, id)

				if !isClient {
					access = list.Public || perm.CheckUserListAccess(nil, id, userId) == nil
				} else {
					access = list.ClientAccess
				}
			}
		}
	} else if mentionType == mentionUserType {
		user := model.GetUserWithOrganizationMemberByOrganizationIdAndUsername(orgaId, mentionId)

		if user != nil {
			title = fmt.Sprintf("%v %v", user.Firstname, user.Lastname)
		}
	} else if mentionType == mentionTagType {
		tag := model.GetTagByOrganizationIdAndName(orgaId, mentionId)

		if tag != nil {
			title = tag.Name
		}
	}

	return title, access
}

func getMentionIdFromString(mentionId, mentionType string) hide.ID {
	id, err := hide.FromString(mentionId)

	if err != nil {
		logbuch.Warn("Error parsing ID from mention while replacing title", logbuch.Fields{"err": err, "type": mentionType, "id": mentionId})
		return 0
	}

	return id
}
