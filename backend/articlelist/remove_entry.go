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

func RemoveArticleListEntry(organization *model.Organization, userId, listId hide.ID, articleIds []hide.ID) error {
	_, err := checkListExists(organization, listId)

	if err != nil {
		return err
	}

	if err := checkUserModAccess(listId, userId); err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to remove article list entries", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	for _, articleId := range articleIds {
		entry := model.GetArticleListEntryByArticleListIdAndArticleIdTx(tx, listId, articleId)

		if entry == nil {
			continue
		}

		if _, err := checkUserArticleAccess(tx, organization.ID, entry.ArticleId, userId); err != nil {
			db.Rollback(tx)
			return err
		}

		if err := model.DeleteArticleListEntryById(tx, entry.ID); err != nil {
			return errs.Saving
		}

		if err := createRemoveArticleFeed(tx, organization, userId, entry.ArticleListId, entry.ArticleId); err != nil {
			logbuch.Error("Error creating feed when removing entry from article list", logbuch.Fields{"err": err})
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when removing article list entries", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func createRemoveArticleFeed(tx *sqlx.Tx, orga *model.Organization, userId, listId, articleId hide.ID) error {
	list := model.GetArticleListByOrganizationIdAndIdTx(tx, orga.ID, listId)

	if list == nil {
		db.Rollback(tx)
		return errs.ArticleListNotFound
	}

	article := model.GetArticleByOrganizationIdAndIdIgnoreArchivedTx(tx, orga.ID, articleId)

	if article == nil {
		db.Rollback(tx)
		return errs.ArticleNotFound
	}

	reason := "remove_article_list_entry"
	langId := util.DetermineLang(tx, orga.ID, userId, 0).ID
	refs := make([]interface{}, 1)
	refs[0] = list

	if article.ReadEveryone {
		refs = append(refs, article, model.GetArticleContentLatestByArticleIdAndLanguageIdTx(tx, article.ID, langId, false))
	} else {
		reason = "remove_protected_article_list_entry"
	}

	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       reason,
		Public:       list.Public,
		Access:       model.FindArticleListMemberUserIdByArticleListIdTx(tx, listId),
		Notify:       model.FindObservedObjectUserIdByArticleListIdTx(tx, listId),
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
