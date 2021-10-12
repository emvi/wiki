package feed

import (
	rendering "emviwiki/shared/feed"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
)

const (
	notificationsMinLimit = 5
	notificationsMaxLimit = 20
	defaultSupportedLang  = "en"
)

var (
	// list of all supported feed translations (ISO-639-1 code)
	supportedLangs = []string{"en", "de"}
)

func GetFilteredFeed(organization *model.Organization, userId hide.ID, filter *model.SearchFeedFilter) ([]model.Feed, int) {
	// set default filter
	if filter == nil {
		filter = new(model.SearchFeedFilter)
	}

	if filter.SortCreated == "" {
		filter.SortCreated = "desc"
	}

	language := util.DetermineLang(nil, organization.ID, userId, 0)
	var feed []model.Feed
	var notificationCount int

	if filter.Notifications {
		filter.Limit = validateLimit(filter.Limit)
		notificationCount = model.CountFeedAccessByOrganizationIdAndUserIdAndNotificationAndRead(organization.ID, userId, true, false)
	}

	feed = model.FindFeedByOrganizationIdAndUserIdAndLanguageIdAndFilterLimit(organization.ID, userId, language.ID, filter)
	feedLangCode := util.DetermineSystemSupportedLangCode(organization.ID, userId)

	for i := range feed {
		reason, ok := rendering.Reasons[feedLangCode][feed[i].Reason]

		if ok {
			feed[i].Feed = rendering.RenderFeed(organization, reason.Feed, rendering.FeedText, feedLangCode, &feed[i])

			if feed[i].Notification == "" {
				feed[i].Notification = feed[i].Feed
			} else {
				feed[i].Notification = rendering.RenderFeed(organization, reason.Notification, rendering.NotificationText, feedLangCode, &feed[i])
			}
		}
	}

	return feed, notificationCount
}

func validateLimit(limit int) int {
	if limit < notificationsMinLimit {
		return notificationsMinLimit
	} else if limit > notificationsMaxLimit {
		return notificationsMaxLimit
	}

	return limit
}

func getSupportedLanguage(orgaId, userId hide.ID) string {
	orgaLang := model.GetDefaultLanguageByOrganizationId(orgaId).Code
	user := model.GetUserByOrganizationIdAndId(orgaId, userId)
	var userLang, supportedOrgaLang string

	if user != nil && user.Language.Valid {
		userLang = user.Language.String
	}

	for _, code := range supportedLangs {
		if code == userLang {
			return code
		}

		if code == orgaLang {
			supportedOrgaLang = code
		}
	}

	if supportedOrgaLang != "" {
		return supportedOrgaLang
	}

	return defaultSupportedLang
}
