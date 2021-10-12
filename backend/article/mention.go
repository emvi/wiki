package article

import (
	"emviwiki/backend/feed"
	"emviwiki/backend/prosemirror"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"time"
)

const (
	mentionTypeName     = "mention"
	mentionUserType     = "user"
	mentionTypeAttr     = "type"
	mentionTimeAttr     = "time"
	mentionUsernameAttr = "id"
	mentionTimeLayout   = time.RFC3339
)

func notifyMentionedUsers(orga *model.Organization, userId hide.ID, lastContentDefTime time.Time, article *model.Article, content *model.ArticleContent) error {
	doc, err := prosemirror.ParseDoc(content.Content)

	if err != nil {
		return err
	}

	notifyUserIds := make([]hide.ID, 0)
	mentions := prosemirror.FindNodes(doc, -1, mentionTypeName)

	for _, mention := range mentions {
		mentionType, mentionTime, mentionUsername := getMentionAttrs(mention)

		if mentionType == mentionUserType && isNewMention(mentionTime, lastContentDefTime) {
			user := model.GetUserWithOrganizationMemberByOrganizationIdAndUsername(orga.ID, mentionUsername)

			if user != nil {
				notifyUserIds = append(notifyUserIds, user.ID)
			}
		}
	}

	if len(notifyUserIds) != 0 {
		createMentionFeed(orga, userId, article, content, notifyUserIds)
	}

	return nil
}

func getMentionAttrs(mention prosemirror.Node) (string, string, string) {
	mentionType, typeOk := mention.Attrs[mentionTypeAttr].(string)
	mentionTime, timeOk := mention.Attrs[mentionTimeAttr].(string)
	mentionId, idOk := mention.Attrs[mentionUsernameAttr].(string)

	if !typeOk || !timeOk || !idOk {
		return "", "", ""
	}

	return mentionType, mentionTime, mentionId
}

func isNewMention(t string, lastContentDefTime time.Time) bool {
	mentioned, err := time.Parse(mentionTimeLayout, t)
	return err == nil && mentioned.After(lastContentDefTime)
}

func createMentionFeed(orga *model.Organization, userId hide.ID, article *model.Article, content *model.ArticleContent, notifyUserIds []hide.ID) {
	refs := make([]interface{}, 2)
	refs[0] = article
	refs[1] = content
	feedData := &feed.CreateFeedData{Organization: orga,
		UserId: userId,
		Reason: "mentioned",
		Public: false,
		Notify: notifyUserIds,
		Refs:   refs}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when mentioning users", logbuch.Fields{"err": err})
	}
}
