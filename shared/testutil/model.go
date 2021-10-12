package testutil

import (
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"strconv"
	"testing"
	"time"
)

func CreateOrgaAndUser(t *testing.T) (*model.Organization, *model.User) {
	user := &model.User{Firstname: "Firstname",
		Lastname:     "Lastname",
		Email:        "test@user.com",
		Language:     null.NewString("en", true),
		Introduction: true}
	user.ID = 123

	if err := model.SaveUser(nil, user, true); err != nil {
		t.Fatal(err)
	}

	orga := &model.Organization{Name: "test",
		NameNormalized: "test",
		OwnerUserId:    user.ID,
		Expert:         true,
		MaxStorageGB:   5}

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	lang := CreateLang(t, orga, "ja", "Japanese", false)

	member := &model.OrganizationMember{OrganizationId: orga.ID,
		UserId:      user.ID,
		LanguageId:  lang.ID,
		Username:    "testuser1",
		IsModerator: true,
		IsAdmin:     true,
		Active:      true,
		User:        user}

	if err := model.SaveOrganizationMember(nil, member); err != nil {
		t.Fatal(err)
	}

	user.OrganizationMember = member

	all := CreateUserGroup(t, orga, constants.GroupAllName)
	admin := CreateUserGroup(t, orga, constants.GroupAdminName)
	mod := CreateUserGroup(t, orga, constants.GroupModName)
	readonly := CreateUserGroup(t, orga, constants.GroupReadOnlyName)
	CreateUserGroupMember(t, all, user, false)
	CreateUserGroupMember(t, admin, user, false)
	CreateUserGroupMember(t, mod, user, false)
	CreateUserGroupMember(t, readonly, user, false)

	return orga, user
}

func CreateOrga(t *testing.T, owner *model.User, name string) (*model.Organization, *model.OrganizationMember) {
	orga := &model.Organization{Name: name,
		NameNormalized: name,
		OwnerUserId:    owner.ID}

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	lang := CreateLang(t, orga, "jp", "Japanese", false)

	member := &model.OrganizationMember{OrganizationId: orga.ID,
		UserId:     owner.ID,
		LanguageId: lang.ID,
		Username:   "testuser1",
		IsAdmin:    true,
		Active:     true}

	if err := model.SaveOrganizationMember(nil, member); err != nil {
		t.Fatal(err)
	}

	return orga, member
}

func CreateUserWithoutOrganization(t *testing.T, id hide.ID, email string) *model.User {
	user := &model.User{Email: email, Language: null.NewString("en", true)}
	user.ID = id

	if err := model.SaveUser(nil, user, true); err != nil {
		t.Fatal(err)
	}

	return user
}

func CreateUser(t *testing.T, orga *model.Organization, id hide.ID, email string) *model.User {
	user := &model.User{Email: email, Language: null.NewString("en", true)}
	user.ID = id

	if err := model.SaveUser(nil, user, true); err != nil {
		t.Fatal(err)
	}

	lang := CreateLang(t, orga, "ru", "Russian", false)

	// find unique username
	usernameNumber := 2 // 1 is used for the default user

	for model.GetOrganizationMemberByUsername("testuser"+strconv.Itoa(usernameNumber)) != nil {
		usernameNumber++
	}

	username := "testuser" + strconv.Itoa(usernameNumber)
	member := &model.OrganizationMember{OrganizationId: orga.ID,
		UserId:             user.ID,
		LanguageId:         lang.ID,
		Username:           username,
		Active:             true,
		User:               user,
		RecommendationMail: true}

	if err := model.SaveOrganizationMember(nil, member); err != nil {
		t.Fatal(err)
	}

	user.OrganizationMember = member
	return user
}

func CreateUserGroup(t *testing.T, orga *model.Organization, name string) *model.UserGroup {
	group := &model.UserGroup{OrganizationId: orga.ID, Name: name}

	if err := model.SaveUserGroup(nil, group); err != nil {
		t.Fatal(err)
	}

	return group
}

func CreateUserGroupMember(t *testing.T, group *model.UserGroup, user *model.User, mod bool) *model.UserGroupMember {
	member := &model.UserGroupMember{IsModerator: mod,
		UserGroupId: group.ID,
		UserId:      user.ID}

	if err := model.SaveUserGroupMember(nil, member); err != nil {
		t.Fatal(err)
	}

	return member
}

func CreateLang(t *testing.T, orga *model.Organization, code, name string, isdefault bool) *model.Language {
	lang := &model.Language{Code: code,
		Name:           name,
		Default:        isdefault,
		OrganizationId: orga.ID}

	if err := model.SaveLanguage(nil, lang); err != nil {
		t.Fatal(err)
	}

	return lang
}

func CreateArticle(t *testing.T, orga *model.Organization, user *model.User, lang *model.Language, r, w bool) *model.Article {
	article := &model.Article{OrganizationId: orga.ID,
		Views:         54,
		ReadEveryone:  r,
		WriteEveryone: w,
		WIP:           -1,
		Published:     null.NewTime(time.Now(), true)}

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	tags := []string{"test", "article", "foo", "bar"}

	for _, name := range tags {
		tag := &model.Tag{OrganizationId: orga.ID, Name: name}

		if err := model.SaveTag(nil, tag); err != nil {
			t.Fatal(err)
		}

		articleTag := &model.ArticleTag{ArticleId: article.ID, TagId: tag.ID}

		if err := model.SaveArticleTag(nil, articleTag); err != nil {
			t.Fatal(err)
		}
	}

	content := []model.ArticleContent{{
		Title:           "title 1",
		Content:         "content 1",
		TitleTsvector:   "title 1",
		ContentTsvector: "content 1",
		Version:         1,
		Commit:          null.NewString("First commit", true),
		ArticleId:       article.ID,
		LanguageId:      lang.ID,
		UserId:          user.ID,
		User:            user,
		SchemaVersion:   constants.LatestSchemaVersion,
	}, {
		Title:           "title 2",
		Content:         "content 2",
		TitleTsvector:   "title 2",
		ContentTsvector: "content 2",
		Version:         2,
		Commit:          null.NewString("Second commit", true),
		ArticleId:       article.ID,
		LanguageId:      lang.ID,
		UserId:          user.ID,
		User:            user,
		SchemaVersion:   constants.LatestSchemaVersion,
	}, {
		Title:           "title 2",
		Content:         "content 2",
		TitleTsvector:   "title 2",
		ContentTsvector: "content 2",
		Version:         0,
		Commit:          null.NewString("Third commit", true),
		ArticleId:       article.ID,
		LanguageId:      lang.ID,
		UserId:          user.ID,
		User:            user,
		SchemaVersion:   constants.LatestSchemaVersion,
	}}

	for i, c := range content {
		if err := model.SaveArticleContent(nil, &c); err != nil {
			t.Fatal(err)
		}

		author := &model.ArticleContentAuthor{ArticleContentId: c.ID,
			UserId: user.ID}

		if err := model.SaveArticleContentAuthor(nil, author); err != nil {
			t.Fatal(err)
		}

		content[i].Authors = []model.User{*user}
	}

	access := &model.ArticleAccess{Write: true,
		UserId:    user.ID,
		ArticleId: article.ID}

	if err := model.SaveArticleAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	article.Access = []model.ArticleAccess{*access}
	article.LatestArticleContent = &content[2]
	return article
}

func CreateArticleWithoutContent(t *testing.T, orga *model.Organization, user *model.User, lang *model.Language, r, w bool) *model.Article {
	article := &model.Article{OrganizationId: orga.ID,
		Views:         54,
		ReadEveryone:  r,
		WriteEveryone: w,
		WIP:           -1,
		Published:     null.NewTime(time.Now(), true)}

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	tags := []string{"test", "article", "foo", "bar"}

	for _, name := range tags {
		tag := &model.Tag{OrganizationId: orga.ID, Name: name}

		if err := model.SaveTag(nil, tag); err != nil {
			t.Fatal(err)
		}

		articleTag := &model.ArticleTag{ArticleId: article.ID, TagId: tag.ID}

		if err := model.SaveArticleTag(nil, articleTag); err != nil {
			t.Fatal(err)
		}
	}

	access := &model.ArticleAccess{Write: true,
		UserId:    user.ID,
		ArticleId: article.ID}

	if err := model.SaveArticleAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	article.Access = []model.ArticleAccess{*access}
	return article
}

func CreateArticleContent(t *testing.T, user *model.User, article *model.Article, lang *model.Language, version int) *model.ArticleContent {
	content := &model.ArticleContent{Title: "content title",
		Content:       "content",
		Version:       version,
		Commit:        null.NewString("content commit", true),
		ArticleId:     article.ID,
		LanguageId:    lang.ID,
		UserId:        user.ID,
		SchemaVersion: constants.LatestSchemaVersion}

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	return content
}

func CreateArticleContentAuthor(t *testing.T, user *model.User, content *model.ArticleContent) *model.ArticleContentAuthor {
	author := &model.ArticleContentAuthor{ArticleContentId: content.ID, UserId: user.ID}

	if err := model.SaveArticleContentAuthor(nil, author); err != nil {
		t.Fatal(err)
	}

	return author
}

func CreateArticleAccess(t *testing.T, article *model.Article, user *model.User, group *model.UserGroup, write bool) *model.ArticleAccess {
	userId := hide.ID(0)
	groupId := hide.ID(0)

	if user != nil {
		userId = user.ID
	}

	if group != nil {
		groupId = group.ID
	}

	access := &model.ArticleAccess{Write: write,
		UserId:      userId,
		UserGroupId: groupId,
		ArticleId:   article.ID}

	if err := model.SaveArticleAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	return access
}

func CreateFeed(t *testing.T, orga *model.Organization, user *model.User, lang *model.Language, notification bool) *model.Feed {
	group := CreateUserGroup(t, orga, "group")
	article := CreateArticle(t, orga, user, lang, true, true)
	list, _ := CreateArticleList(t, orga, user, lang, true)
	feed := &model.Feed{OrganizationId: orga.ID,
		Public:            false,
		Reason:            "joined_organization",
		TriggeredByUserId: user.ID}

	if err := model.SaveFeed(nil, feed); err != nil {
		t.Fatal(err)
	}

	userRef := &model.FeedRef{FeedId: feed.ID, UserID: user.ID, User: user}
	groupRef := &model.FeedRef{FeedId: feed.ID, UserGroupID: group.ID, UserGroup: group}
	articleRef := &model.FeedRef{FeedId: feed.ID, ArticleID: article.ID, Article: article}
	listRef := &model.FeedRef{FeedId: feed.ID, ArticleListID: list.ID, ArticleList: list}

	if err := model.SaveFeedRef(nil, userRef); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveFeedRef(nil, groupRef); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveFeedRef(nil, articleRef); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveFeedRef(nil, listRef); err != nil {
		t.Fatal(err)
	}

	access := &model.FeedAccess{UserId: user.ID, FeedId: feed.ID, Notification: notification, Read: false}

	if err := model.SaveFeedAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	return feed
}

func CreateFeedForObject(t *testing.T, orga *model.Organization, user *model.User, article *model.Article, list *model.ArticleList, group *model.UserGroup) *model.Feed {
	feed := &model.Feed{OrganizationId: orga.ID,
		Public:            true,
		Reason:            "reason_name",
		TriggeredByUserId: user.ID}

	if err := model.SaveFeed(nil, feed); err != nil {
		t.Fatal(err)
	}

	access := &model.FeedAccess{FeedId: feed.ID, UserId: user.ID}

	if err := model.SaveFeedAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	articleId := hide.ID(0)
	listId := hide.ID(0)
	groupId := hide.ID(0)

	if article != nil {
		articleId = article.ID
	} else if list != nil {
		listId = list.ID
	} else if group != nil {
		groupId = group.ID
	}

	ref := &model.FeedRef{FeedId: feed.ID,
		ArticleID:     articleId,
		ArticleListID: listId,
		UserGroupID:   groupId}

	if err := model.SaveFeedRef(nil, ref); err != nil {
		t.Fatal(err)
	}

	return feed
}

func CreateObservedObject(t *testing.T, user *model.User, article *model.Article, list *model.ArticleList, group *model.UserGroup) *model.ObservedObject {
	articleId := hide.ID(0)
	listId := hide.ID(0)
	groupId := hide.ID(0)

	if article != nil {
		articleId = article.ID
	} else if list != nil {
		listId = list.ID
	} else if group != nil {
		groupId = group.ID
	}

	observed := &model.ObservedObject{UserId: user.ID,
		ArticleId:     articleId,
		ArticleListId: listId,
		UserGroupId:   groupId}

	if err := model.SaveObservedObject(nil, observed); err != nil {
		t.Fatal(err)
	}

	return observed
}

func CreateArticleList(t *testing.T, orga *model.Organization, user *model.User, lang *model.Language, public bool) (*model.ArticleList, *model.ArticleListMember) {
	list := &model.ArticleList{OrganizationId: orga.ID, Public: public}

	if err := model.SaveArticleList(nil, list); err != nil {
		t.Fatal(err)
	}

	name := &model.ArticleListName{ArticleListId: list.ID,
		LanguageId: lang.ID,
		Name:       "article list name",
		Info:       null.NewString("article list info", true)}

	if err := model.SaveArticleListName(nil, name); err != nil {
		t.Fatal(err)
	}

	var member *model.ArticleListMember

	if user != nil {
		member = &model.ArticleListMember{ArticleListId: list.ID,
			UserId:      user.ID,
			IsModerator: true}

		if err := model.SaveArticleListMember(nil, member); err != nil {
			t.Fatal(err)
		}
	}

	list.Name = name
	return list, member
}

func CreateArticleListEntry(t *testing.T, list *model.ArticleList, article *model.Article, pos uint) *model.ArticleListEntry {
	entry := &model.ArticleListEntry{ArticleListId: list.ID,
		ArticleId: article.ID,
		Position:  pos}

	if err := model.SaveArticleListEntry(nil, entry); err != nil {
		t.Fatal(err)
	}

	return entry
}

func CreateArticleListMember(t *testing.T, list *model.ArticleList, userId, groupId hide.ID, mod bool) *model.ArticleListMember {
	member := &model.ArticleListMember{ArticleListId: list.ID,
		UserId:      userId,
		UserGroupId: groupId,
		IsModerator: mod}

	if err := model.SaveArticleListMember(nil, member); err != nil {
		t.Fatal(err)
	}

	return member
}

func CreateArticleListName(t *testing.T, list *model.ArticleList, lang *model.Language, name, info string) *model.ArticleListName {
	listname := &model.ArticleListName{ArticleListId: list.ID,
		LanguageId: lang.ID,
		Name:       name,
		Info:       null.NewString(info, info != "")}

	if err := model.SaveArticleListName(nil, listname); err != nil {
		t.Fatal(err)
	}

	return listname
}

func CreateTag(t *testing.T, orga *model.Organization, name string) *model.Tag {
	tag := &model.Tag{OrganizationId: orga.ID, Name: name}

	if err := model.SaveTag(nil, tag); err != nil {
		t.Fatal(err)
	}

	return tag
}

func CreateArticleTag(t *testing.T, article *model.Article, tag *model.Tag) *model.ArticleTag {
	articleTag := &model.ArticleTag{ArticleId: article.ID, TagId: tag.ID}

	if err := model.SaveArticleTag(nil, articleTag); err != nil {
		t.Fatal(err)
	}

	return articleTag
}

func CreateFile(t *testing.T, orga *model.Organization, user *model.User, article *model.Article, roomId string) *model.File {
	return CreateFileWithMd5(t, orga, user, article, roomId, "md5")
}

func CreateFileWithMd5(t *testing.T, orga *model.Organization, user *model.User, article *model.Article, roomId string, md5 string) *model.File {
	var orgaId hide.ID
	var articleId hide.ID

	if orga != nil {
		orgaId = orga.ID
	}

	if article != nil {
		articleId = article.ID
	}

	uniqueName := "unique_name"
	index := 0

	for model.GetFileByUniqueName(uniqueName+strconv.Itoa(index)) != nil {
		index++
	}

	uniqueName += strconv.Itoa(index)

	file := &model.File{OrganizationId: orgaId,
		UserId:       user.ID,
		ArticleId:    articleId,
		RoomId:       null.NewString(roomId, roomId != ""),
		OriginalName: "original_name",
		UniqueName:   uniqueName,
		Path:         "some/test/file.txt",
		Type:         ".txt",
		Size:         42,
		MD5:          md5}

	if err := model.SaveFile(nil, file); err != nil {
		t.Fatal(err)
	}

	return file
}

func CreateBookmark(t *testing.T, orga *model.Organization, user *model.User, article *model.Article, list *model.ArticleList) *model.Bookmark {
	var articleId null.Int64
	var listId null.Int64

	if article != nil {
		articleId = null.NewInt64(int64(article.ID), true)
	}

	if list != nil {
		listId = null.NewInt64(int64(list.ID), true)
	}

	bookmark := &model.Bookmark{OrganizationId: orga.ID,
		UserId:        user.ID,
		ArticleId:     articleId,
		ArticleListId: listId}

	if err := model.SaveBookmark(nil, bookmark); err != nil {
		t.Fatal(err)
	}

	return bookmark
}

func CreateSupportTicket(t *testing.T, orga *model.Organization, user *model.User) *model.SupportTicket {
	ticket := &model.SupportTicket{OrganizationId: orga.ID,
		UserId:  user.ID,
		Type:    "type",
		Message: "message",
		Status:  "open"}

	if err := model.SaveSupportTicket(nil, ticket); err != nil {
		t.Fatal(err)
	}

	return ticket
}

func CreateClient(t *testing.T, orga *model.Organization, name, id, secret string) *model.Client {
	client := &model.Client{OrganizationId: orga.ID, Name: name, ClientId: id, ClientSecret: secret}

	if err := model.SaveClient(nil, client); err != nil {
		t.Fatal(err)
	}

	return client
}

func CreateClientScope(t *testing.T, client *model.Client, name string, r, w bool) *model.ClientScope {
	scope := &model.ClientScope{ClientId: client.ID, Name: name, Read: r, Write: w}

	if err := model.SaveClientScope(nil, scope); err != nil {
		t.Fatal(err)
	}

	return scope
}

func CreateInvitation(t *testing.T, orga *model.Organization, email, code string, readOnly bool) *model.Invitation {
	inv := &model.Invitation{OrganizationId: orga.ID,
		Email:    email,
		Code:     code,
		ReadOnly: readOnly}

	if err := model.SaveInvitation(nil, inv); err != nil {
		t.Fatal(err)
	}

	return inv
}

func CreateArticleVisit(t *testing.T, article *model.Article, user *model.User) *model.ArticleVisit {
	visit := &model.ArticleVisit{ArticleId: article.ID, UserId: user.ID}

	if err := model.SaveArticleVisit(nil, visit); err != nil {
		t.Fatal(err)
	}

	return visit
}

func CreateArticleRecommendation(t *testing.T, article *model.Article, user, recommendedTo *model.User) *model.ArticleRecommendation {
	entity := &model.ArticleRecommendation{ArticleId: article.ID, UserId: user.ID, RecommendedTo: recommendedTo.ID}

	if err := model.SaveArticleRecommendation(nil, entity); err != nil {
		t.Fatal(err)
	}

	return entity
}
