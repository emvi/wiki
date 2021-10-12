package article

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
	"time"
)

func TestNotifyMentionedUsers(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := model.GetArticleContentLastByArticleIdAndLanguageIdAndWIP(article.ID, lang.ID, false)
	content.Content = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hallo "},{"type":"mention","attrs":{"id":"testuser1","type":"user","title":"Marvin Blum","time":"2019-08-10T16:24:22.127Z"}},{"type":"text","text":", wie geht es dir? "},{"type":"mention","attrs":{"id":"Bme1VxaXyK","type":"list","title":"Testliste","time":"2019-08-10T16:25:22.532Z"}},{"type":"text","text":" bitte lesen, weil "},{"type":"mention","attrs":{"id":"test","type":"tag","title":"test","time":"2019-08-10T16:25:28.173Z"}},{"type":"text","text":" und an "},{"type":"mention","attrs":{"id":"MJnglZRd0L","type":"group","title":"Administrators","time":"2019-08-10T16:25:33.100Z"}},{"type":"text","text":" weiterleiten. "},{"type":"mention","attrs":{"id":"bnBaMwN1GV","type":"article","title":"Weiterer Test","time":"2019-08-10T16:24:17.660Z"}},{"type":"text","text":" Der hier "},{"type":"mention","attrs":{"id":"max","type":"user","title":"Max Mustermann","time":"2019-08-10T16:28:17.624Z"}},{"type":"text","text":" darf nicht benachrichtigt werden!"}]}]}`
	lastVersion, _ := time.Parse("2006-01-02", "2019-08-09")

	if err := notifyMentionedUsers(orga, user2.ID, lastVersion, article, content); err != nil {
		t.Fatalf("Mentioned users must be found, but was: %v", err)
	}

	testutil.AssertFeedCreated(t, orga, "mentioned")
}
