package articlelist

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func DeleteArticleList(organization *model.Organization, userId, listId hide.ID) error {
	if _, err := checkListExists(organization, listId); err != nil {
		return err
	}

	if err := checkUserModAccess(listId, userId); err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error beginning transaction when deleting article list", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	// create feed before deleting list
	if err := createDeletedListFeed(tx, organization, userId, listId); err != nil {
		logbuch.Error("Error creating feed when deleting article list", logbuch.Fields{"err": err})
		return errs.Saving
	}

	if err := feed.DeleteFeed(tx, &feed.DeleteFeedData{ListId: listId}); err != nil {
		logbuch.Error("Error deleting article list feed", logbuch.Fields{"err": err, "list_id": listId, "user_id": userId})
		return errs.Saving
	}

	if err := model.DeleteArticleListById(tx, listId); err != nil {
		return errs.Saving
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when deleting article list", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func checkListExists(organization *model.Organization, listId hide.ID) (*model.ArticleList, error) {
	list := model.GetArticleListByOrganizationIdAndId(organization.ID, listId)

	if list == nil {
		return nil, errs.ArticleListNotFound
	}

	return list, nil
}

func createDeletedListFeed(tx *sqlx.Tx, orga *model.Organization, userId, listId hide.ID) error {
	list := model.GetArticleListByOrganizationIdAndIdTx(tx, orga.ID, listId)

	if list == nil {
		db.Rollback(tx)
		return errs.ArticleListNotFound
	}

	langId := util.DetermineLang(tx, orga.ID, userId, 0).ID
	name := model.GetArticleListNameByOrganizationIdAndArticleListIdAndLangIdTx(tx, orga.ID, listId, langId)
	refs := make([]interface{}, 1)
	refs[0] = feed.KeyValue{"name", name.Name}
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "delete_articlelist",
		Public:       list.Public,
		Access:       model.FindArticleListMemberUserIdByArticleListIdTx(tx, listId),
		Notify:       model.FindObservedObjectUserIdByArticleListIdTx(tx, listId),
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
