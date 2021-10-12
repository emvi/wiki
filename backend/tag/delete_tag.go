package tag

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func DeleteTag(orga *model.Organization, userId, tagId hide.ID) error {
	if _, err := perm.CheckUserIsAdminOrMod(orga.ID, userId); err != nil {
		return err
	}

	if err := perm.CheckUserTagAccess(orga.ID, userId, tagId); err != nil {
		return err
	}

	tag := model.GetTagByOrganizationIdAndId(orga.ID, tagId)

	if tag == nil {
		return errs.TagNotFound
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error beginning transaction when deleting tag", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	if err := model.DeleteArticleTagByOrganizationIdAndTagId(tx, orga.ID, tagId); err != nil {
		return errs.Saving
	}

	if err := model.DeleteTagByOrganizationIdAndId(tx, orga.ID, tagId); err != nil {
		return errs.Saving
	}

	if err := createDeletedTagFeed(tx, orga, userId, tag.Name); err != nil {
		logbuch.Error("Error creating feed when deleting tag", logbuch.Fields{"err": err})
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when deleting tag", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func createDeletedTagFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, name string) error {
	refs := make([]interface{}, 1)
	refs[0] = feed.KeyValue{"name", name}
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "delete_tag",
		Public:       true,
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
