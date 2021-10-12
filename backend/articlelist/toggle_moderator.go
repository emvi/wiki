package articlelist

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func ToggleArticleListModerator(organization *model.Organization, userId, listId, memberId hide.ID) error {
	if err := checkUserModAccess(listId, userId); err != nil {
		return err
	}

	member := model.GetArticleListMemberByArticleListIdAndId(listId, memberId)

	if member == nil {
		return errs.MemberNotFound
	}

	if member.UserId != 0 && userId == member.UserId {
		return errs.ModeratorYourself
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to toggle article list moderator", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	member.IsModerator = !member.IsModerator

	if err := model.SaveArticleListMember(tx, member); err != nil {
		return errs.Saving
	}

	if err := createToggleModeratorFeed(tx, organization, userId, member); err != nil {
		logbuch.Error("Error creating feed when toggling article list moderator", logbuch.Fields{"err": err})
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when toggling article list moderator", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func createToggleModeratorFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, member *model.ArticleListMember) error {
	reason := "set_article_list_moderator"

	if !member.IsModerator {
		reason = "remove_article_list_moderator"
	}

	list := model.GetArticleListByOrganizationIdAndIdTx(tx, orga.ID, member.ArticleListId)

	if list == nil {
		db.Rollback(tx)
		return errs.ArticleListNotFound
	}

	memberUser := model.GetUserByOrganizationIdAndIdTx(tx, orga.ID, member.UserId)

	if memberUser == nil {
		logbuch.Error("User not found to create toggle article list moderator", logbuch.Fields{"list_id": list.ID, "user_id": member.UserId})
		db.Rollback(tx)
		return errs.UserNotFound
	}

	refs := make([]interface{}, 2)
	refs[0] = list
	refs[1] = memberUser

	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       reason,
		Public:       false,
		Notify:       []hide.ID{memberUser.ID},
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
