package articlelist

import (
	"emviwiki/backend/observe"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"strings"
	"unicode/utf8"

	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
)

const (
	nameMaxLen = 40
	infoMaxLen = 100
	maxLists   = 10
)

type SaveArticleListData struct {
	Id           hide.ID                   `json:"id"`
	Names        []SaveArticleListNameData `json:"names"`
	Public       bool                      `json:"public"`
	ClientAccess bool                      `json:"client_access"`
}

type SaveArticleListNameData struct {
	LanguageId hide.ID `json:"language_id"`
	Name       string  `json:"name"`
	Info       string  `json:"info"`
}

func (data *SaveArticleListData) validate(orgaId hide.ID) []error {
	if len(data.Names) == 0 {
		return []error{errs.NoNamesProvided}
	}

	defaultLangId := model.GetDefaultLanguageByOrganizationId(orgaId).ID
	names := make([]SaveArticleListNameData, 0)
	err := make([]error, 0)

	for _, name := range data.Names {
		name.Name = strings.TrimSpace(name.Name)
		name.Info = strings.TrimSpace(name.Info)
		langId, _ := hide.ToString(name.LanguageId) // set name in front of field name to distinguish between name and info errors for the same language

		// if not default language and empty, this name is optional
		if name.LanguageId != defaultLangId && name.Name == "" && name.Info == "" {
			continue
		}

		if len(name.Name) == 0 {
			err = append(err, rest.NewApiError(errs.NameTooShort.Message, "name."+langId))
		}

		if utf8.RuneCountInString(name.Name) > nameMaxLen {
			err = append(err, rest.NewApiError(errs.NameTooLong.Message, "name."+langId))
		}

		if utf8.RuneCountInString(name.Info) > infoMaxLen {
			err = append(err, rest.NewApiError(errs.InfoTooLong.Message, "info."+langId))
		}

		if model.GetLanguageByOrganizationIdAndId(orgaId, name.LanguageId) == nil {
			err = append(err, rest.NewApiError(errs.LanguageNotFound.Message, langId))
		}

		names = append(names, name)
	}

	data.Names = names

	if len(err) == 0 {
		return nil
	}

	return err
}

func SaveArticleList(organization *model.Organization, userId hide.ID, data SaveArticleListData) (hide.ID, []error) {
	if !organization.Expert {
		if model.CountArticleListByOrganizationId(organization.ID) >= maxLists {
			return 0, []error{errs.MaxListsReached}
		}
	}

	_, err := checkListExists(organization, data.Id)

	if data.Id != 0 && err != nil {
		return 0, []error{err}
	}

	err = checkUserModAccess(data.Id, userId)

	if data.Id != 0 && err != nil {
		return 0, []error{err}
	}

	if err := data.validate(organization.ID); err != nil {
		return 0, err
	}

	// save article list, names and first member (if new)
	isNew := data.Id == 0
	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to save article list", logbuch.Fields{"err": err})
		return 0, []error{errs.TxBegin}
	}

	list, err := saveArticleList(tx, organization.ID, data.Id, data.Public, data.ClientAccess)

	if err != nil {
		return 0, []error{err}
	}

	if err := saveArticleListNames(tx, list.ID, data.Names); err != nil {
		return 0, []error{err}
	}

	if isNew {
		if err := saveFirstArticleListMember(tx, list.ID, userId); err != nil {
			return 0, []error{err}
		}
	}

	if err := createSaveArticleListFeed(tx, organization, userId, list, isNew); err != nil {
		logbuch.Error("Error creating feed when saving article list", logbuch.Fields{"err": err})
		return 0, []error{err}
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when saving article list", logbuch.Fields{"err": err})
		return 0, []error{errs.TxCommit}
	}

	if isNew {
		if err := observe.ObserveObject(organization, userId, 0, list.ID, 0); err != nil {
			logbuch.Error("Error observing new article list", logbuch.Fields{"err": err})
		}
	}

	return list.ID, nil
}

// Checks the given user is moderator of given article list.
// Users that have access due to a user group having access to the list, cannot be moderators.
func checkUserModAccess(listId, userId hide.ID) error {
	if len(model.FindArticleListMemberModeratorByArticleListIdAndUserId(listId, userId)) == 0 {
		return errs.PermissionDenied
	}

	return nil
}

func saveArticleList(tx *sqlx.Tx, orgaId, listId hide.ID, public, clientAccess bool) (*model.ArticleList, error) {
	var list *model.ArticleList

	if listId == 0 {
		list = &model.ArticleList{}
	} else {
		list = model.GetArticleListByOrganizationIdAndIdTx(tx, orgaId, listId)

		if list == nil {
			db.Rollback(tx)
			return nil, errs.ArticleListNotFound
		}
	}

	list.OrganizationId = orgaId
	list.Public = public || clientAccess
	list.ClientAccess = clientAccess

	if err := model.SaveArticleList(tx, list); err != nil {
		return nil, errs.Saving
	}

	return list, nil
}

func saveArticleListNames(tx *sqlx.Tx, listId hide.ID, names []SaveArticleListNameData) error {
	if err := model.DeleteArticleListNameByArticleListId(tx, listId); err != nil {
		return errs.Saving
	}

	for _, name := range names {
		listName := &model.ArticleListName{ArticleListId: listId,
			LanguageId: name.LanguageId,
			Name:       name.Name,
			Info:       null.NewString(name.Info, name.Info != "")}

		if err := model.SaveArticleListName(tx, listName); err != nil {
			return errs.Saving
		}
	}

	return nil
}

func saveFirstArticleListMember(tx *sqlx.Tx, listId, userId hide.ID) error {
	member := &model.ArticleListMember{ArticleListId: listId,
		UserId:      userId,
		IsModerator: true}

	if err := model.SaveArticleListMember(tx, member); err != nil {
		return errs.Saving
	}

	return nil
}

func createSaveArticleListFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, list *model.ArticleList, isNew bool) error {
	var access, notify []hide.ID
	reason := "create_article_list"

	if !isNew {
		access = model.FindArticleListMemberUserIdByArticleListIdTx(tx, list.ID)
		notify = model.FindObservedObjectUserIdByArticleListIdTx(tx, list.ID)
		reason = "update_article_list"
	} else {
		access = []hide.ID{userId}
	}

	refs := make([]interface{}, 1)
	refs[0] = list
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       reason,
		Public:       list.Public,
		Access:       access,
		Notify:       notify,
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
