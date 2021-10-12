package articlelist

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func SwapArticleListEntries(organization *model.Organization, userId, listId, articleIdA, articleIdB hide.ID) error {
	if articleIdA == articleIdB || articleIdA == 0 || articleIdB == 0 {
		return errs.InvalidPos
	}

	if _, err := checkListExists(organization, listId); err != nil {
		return err
	}

	if err := checkUserModAccess(listId, userId); err != nil {
		return err
	}

	entryA := model.GetArticleListEntryByArticleListIdAndArticleId(listId, articleIdA)
	entryB := model.GetArticleListEntryByArticleListIdAndArticleId(listId, articleIdB)

	if entryA == nil || entryB == nil {
		return errs.PosNotFound
	}

	// swap positions
	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to swap article list entries", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	if err := swapEntries(tx, entryA, entryB); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when swapping article list entries", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func SortArticleListEntry(organization *model.Organization, userId, listId, articleId hide.ID, direction int) error {
	if direction == 0 {
		return nil
	}

	if _, err := checkListExists(organization, listId); err != nil {
		return err
	}

	if err := checkUserModAccess(listId, userId); err != nil {
		return err
	}

	entryCount := model.CountArticleListEntryByArticleListId(listId)

	if entryCount <= 1 {
		return nil
	}

	entry := model.GetArticleListEntryByArticleListIdAndArticleId(listId, articleId)

	if entry == nil {
		return errs.PosNotFound
	}

	var swapWith *model.ArticleListEntry

	if direction < 0 {
		swapWith = model.GetArticleListEntryByArticleListIdAndPosition(listId, entry.Position-1)
	} else {
		swapWith = model.GetArticleListEntryByArticleListIdAndPosition(listId, entry.Position+1)
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to sort article list entry", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	if swapWith != nil {
		if err := swapEntries(tx, entry, swapWith); err != nil {
			return err
		}
	} else {
		if direction < 0 {
			entry.Position = uint(entryCount + 1)
		} else {
			entry.Position = 0
		}

		if err := model.SaveArticleListEntry(tx, entry); err != nil {
			return errs.Saving
		}

		if err := model.UpdateArticleListEntryPositionByArticleListIdTx(tx, listId); err != nil {
			return errs.Saving
		}
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when sorting article list entry", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func swapEntries(tx *sqlx.Tx, entryA, entryB *model.ArticleListEntry) error {
	logbuch.Debug("Swapping list entries", logbuch.Fields{"a": entryA.Position, "b": entryB.Position})
	entryA.Position, entryB.Position = entryB.Position, entryA.Position

	if err := model.SaveArticleListEntry(tx, entryA); err != nil {
		return errs.Saving
	}

	if err := model.SaveArticleListEntry(tx, entryB); err != nil {
		return errs.Saving
	}

	return nil
}
