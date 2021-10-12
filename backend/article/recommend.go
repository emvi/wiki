package article

import (
	"bytes"
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/backend/mailtpl"
	"emviwiki/backend/perm"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"html"
	templ "html/template"
	"strings"
	"unicode/utf8"
)

const (
	recommendArticleFeed = "recommend_article"
	recommendMailSubject = "recommend_article"
	recommendPath        = "read"

	// if you change this, make sure you update the validation in the frontend too!
	maxMessageLen = 500
)

var recommendMailI18n = i18n.Translation{
	"en": {
		"title-1":  "recommended the article",
		"title-2":  "to you on Emvi.",
		"text":     "You can read it by following the link below.",
		"action":   "Open Article",
		"link":     "Or paste this link into your browser",
		"greeting": "You have an article recommendation!",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title-1":  "hat Dir den Artikel",
		"title-2":  "auf Emvi empfohlen.",
		"text":     "Klicke auf den folgenden Link um ihn zu lesen.",
		"action":   "Artikel öffnen",
		"link":     "Oder kopiere diesen Link in deinen Browser",
		"greeting": "Du hast eine Artikelempfehlung!",
		"goodbye":  "Dein Emvi Team",
	},
}

func RecommendArticle(orga *model.Organization, userId, articleId hide.ID, notifyUserIds, notifyGroupIds []hide.ID, message string, receiveReadConfirmation bool, mailer mail.Sender) error {
	message = strings.TrimSpace(message)

	if err := checkMessageLen(message); err != nil {
		return err
	}

	article, err := checkUserReadAccess(orga.ID, userId, articleId)

	if err != nil {
		return err
	}

	if err := joinLatestArticleContent(orga.ID, userId, article); err != nil {
		return err
	}

	user := model.GetUserById(userId)

	if user == nil {
		return errs.UserNotFound
	}

	notifyUsers, err := checkAllUserExist(orga.ID, notifyUserIds)

	if err != nil {
		return err
	}

	notifyUsers = appendGroupMembers(orga.ID, notifyUsers, notifyGroupIds)
	createRecommendInviteArticleFeed(orga, userId, notifyUsers, article, "", message, recommendArticleFeed)
	sendRecommendInviteArticleMail(orga, article, "", message, user, notifyUsers, mailtpl.RecommendMailTemplate, recommendMailSubject, recommendPath, recommendMailI18n, mailer)

	if receiveReadConfirmation {
		createRecommendationConfirmations(article.ID, userId, notifyUsers)
	}

	return nil
}

func checkMessageLen(message string) error {
	if utf8.RuneCountInString(message) > maxMessageLen {
		return errs.MessageTooLong
	}

	return nil
}

func checkUserReadAccess(orgaId, userId, articleId hide.ID) (*model.Article, error) {
	article := model.GetArticleByOrganizationIdAndId(orgaId, articleId)

	if article == nil {
		return nil, errs.ArticleNotFound
	}

	if !article.ReadEveryone && !perm.CheckUserReadOrWriteAccess(articleId, userId) {
		return nil, errs.PermissionDenied
	}

	return article, nil
}

func joinLatestArticleContent(orgaId, userId hide.ID, article *model.Article) error {
	userLang := util.DetermineLang(nil, orgaId, userId, 0)
	article.LatestArticleContent = model.GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageId(orgaId, article.ID, userLang.ID, false)

	if article.LatestArticleContent == nil {
		logbuch.Error("Latest article content for article not found", logbuch.Fields{"article_id": article.ID, "language_id": userLang.ID})
		return errs.FindingLatestArticleContent
	}

	return nil
}

func checkAllUserExist(orgaId hide.ID, userIds []hide.ID) ([]model.User, error) {
	user := make([]model.User, len(userIds))

	for i, notifyUserId := range userIds {
		notifyUser := model.GetUserWithOrganizationMemberByOrganizationIdAndId(orgaId, notifyUserId)

		if notifyUser == nil {
			return nil, errs.UserNotFound
		}

		user[i] = *notifyUser
	}

	return user, nil
}

func appendGroupMembers(orgaId hide.ID, user []model.User, notifyGroupIds []hide.ID) []model.User {
	for _, groupId := range notifyGroupIds {
		user = append(user, model.FindUserGroupMemberUserByOrganizationIdAndUserGroupId(orgaId, groupId)...)
	}

	return removeDuplicateUser(user)
}

func removeDuplicateUser(user []model.User) []model.User {
	encountered := make(map[hide.ID]bool)
	out := make([]model.User, 0)

	for _, entry := range user {
		if _, value := encountered[entry.ID]; !value {
			encountered[entry.ID] = true
			out = append(out, entry)
		}
	}

	return out
}

func createRecommendInviteArticleFeed(organization *model.Organization, recommendingUserId hide.ID, notifyUsers []model.User, article *model.Article, roomId, message, reason string) {
	for _, notifyUser := range notifyUsers {
		refs := make([]interface{}, 1)
		refs[0] = &notifyUser

		if message != "" {
			refs = append(refs, feed.KeyValue{"message", message})
		}

		if article != nil {
			refs = append(refs, article)
			refs = append(refs, article.LatestArticleContent)
			langHashId, _ := hide.ToString(article.LatestArticleContent.LanguageId)
			refs = append(refs, feed.KeyValue{"lang_id", langHashId})
		} else {
			refs = append(refs, feed.RoomID(roomId))
		}

		feedData := &feed.CreateFeedData{Organization: organization,
			UserId: recommendingUserId,
			Reason: reason,
			Public: false,
			Access: []hide.ID{},
			Notify: []hide.ID{notifyUser.ID},
			Refs:   refs}

		if err := feed.CreateFeed(feedData); err != nil {
			var articleId hide.ID

			if article != nil {
				articleId = article.ID
			}

			logbuch.Error("Error creating feed when recommending article", logbuch.Fields{"err": err, "user_id": notifyUser.ID, "article_id": articleId})
		}
	}
}

func sendRecommendInviteArticleMail(orga *model.Organization, article *model.Article, roomId, message string, recommendingUser *model.User, notifyUsers []model.User, template, subject, path string, vars i18n.Translation, mailer mail.Sender) []error {
	var articleHashId, langHashId string

	if article != nil {
		var err error
		articleHashId, err = hide.ToString(article.ID)

		if err != nil {
			logbuch.Error("Error encoding article id to hash", logbuch.Fields{"err": err, "article_id": article.ID})
		}

		langHashId, err = hide.ToString(article.LatestArticleContent.LanguageId)

		if err != nil {
			logbuch.Error("Error encoding language id to hash", logbuch.Fields{"err": err, "lang_id": article.LatestArticleContent.LanguageId})
		}
	}

	var articleId hide.ID

	if article != nil {
		articleId = article.ID
	}

	e := make([]error, 0)

	for _, notifyUser := range notifyUsers {
		if notifyUser.OrganizationMember.RecommendationMail {
			userLang := util.DetermineSystemSupportedLangCode(orga.ID, notifyUser.ID)
			tpl := mailtpl.Cache.Get()
			var buffer bytes.Buffer
			data := struct {
				User      *model.User
				Article   *model.Article
				ArticleId string
				LangId    string
				RoomId    string
				OrgaURL   string
				Path      string
				Message   templ.HTML
				EndVars   map[string]templ.HTML
				Vars      map[string]templ.HTML
			}{
				recommendingUser,
				article,
				articleHashId,
				langHashId,
				roomId,
				util.InjectSubdomain(frontendHost, orga.NameNormalized),
				path,
				templ.HTML(strings.ReplaceAll(html.EscapeString(message), "\n", "<br />")),
				i18n.GetMailEndI18n(userLang),
				i18n.GetVars(userLang, vars),
			}

			if err := tpl.ExecuteTemplate(&buffer, template, &data); err != nil {
				logbuch.Error("Error rendering mail to recommend article", logbuch.Fields{"err": err, "article_id": articleId, "room_id": roomId})
				e = append(e, err)
				continue
			}

			subject := i18n.GetMailTitle(userLang)[subject]

			if err := mailer(subject, buffer.String(), notifyUser.Email); err != nil {
				logbuch.Error("Error sending mail to recommend article", logbuch.Fields{"err": err, "article_id": articleId, "room_id": roomId, "user_id": notifyUser.ID})
				e = append(e, err)
			}
		}
	}

	return e
}

func createRecommendationConfirmations(articleId, userId hide.ID, notifyUsers []model.User) {
	for _, notifyUser := range notifyUsers {
		if model.GetArticleRecommendationByArticleIdAndUserIdAndRecommendedTo(articleId, userId, notifyUser.ID) != nil {
			continue
		}

		recommendation := &model.ArticleRecommendation{ArticleId: articleId, UserId: userId, RecommendedTo: notifyUser.ID}

		if err := model.SaveArticleRecommendation(nil, recommendation); err != nil {
			logbuch.Error("Error saving article recommendation", logbuch.Fields{"err": err, "article_id": articleId, "user_id": userId, "recommend_to": notifyUser.ID})
		}
	}
}
