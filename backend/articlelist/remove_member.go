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

func RemoveArticleListMember(organization *model.Organization, userId, listId hide.ID, memberIds []hide.ID) error {
	if err := checkUserModAccess(listId, userId); err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to remove article list members", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	removedMembers, err := removeMember(tx, userId, listId, memberIds)

	if err != nil {
		return err
	}

	// check someone has moderator access after operation
	if err := checkModAccessRemains(tx, listId); err != nil {
		return err
	}

	if err := createRemoveMemberFeed(tx, organization, userId, listId, removedMembers); err != nil {
		logbuch.Error("Error creating feed when removing member from article list", logbuch.Fields{"err": err})
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when removing article list members", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func removeMember(tx *sqlx.Tx, userId, listId hide.ID, memberIds []hide.ID) ([]model.ArticleListMember, error) {
	removedMembers := make([]model.ArticleListMember, 0, len(memberIds))

	for _, memberId := range memberIds {
		member := model.GetArticleListMemberByArticleListIdAndIdTx(tx, listId, memberId)

		if member == nil {
			db.Rollback(tx)
			return nil, errs.MemberNotFound
		}

		if member.UserId != 0 && userId == member.UserId {
			db.Rollback(tx)
			return nil, errs.RemoveYourself
		}

		if err := model.DeleteArticleListMemberByArticleListIdAndId(tx, listId, memberId); err != nil {
			return nil, errs.Saving
		}

		removedMembers = append(removedMembers, *member)
	}

	return removedMembers, nil
}

func checkModAccessRemains(tx *sqlx.Tx, listId hide.ID) error {
	if len(model.FindArticleListMemberModeratorByArticleListIdTx(tx, listId)) == 0 {
		db.Rollback(tx)
		return errs.RemoveModeratorAccess
	}

	return nil
}

func createRemoveMemberFeed(tx *sqlx.Tx, orga *model.Organization, userId, listId hide.ID, removedMembers []model.ArticleListMember) error {
	list := model.GetArticleListByOrganizationIdAndIdTx(tx, orga.ID, listId)

	if list == nil {
		db.Rollback(tx)
		return errs.ArticleListNotFound
	}

	refs := make([]interface{}, 2)
	refs[0] = list

	for i, member := range removedMembers {
		if member.UserId != 0 {
			memberUser := model.GetUserByOrganizationIdAndIdTx(tx, orga.ID, member.UserId)

			if memberUser == nil {
				logbuch.Error("User not found to create remove article list member feed", logbuch.Fields{"list_id": list.ID, "user_id": member.UserId})
				db.Rollback(tx)
				return errs.UserNotFound
			}

			removedMembers[i].User = memberUser
			refs[1] = memberUser
		} else {
			memberGroup := model.GetUserGroupByOrganizationIdAndIdTx(tx, orga.ID, member.UserGroupId)

			if memberGroup == nil {
				logbuch.Error("User group not found to create remove article list member feed", logbuch.Fields{"list_id": list.ID, "user_id": member.UserGroupId})
				db.Rollback(tx)
				return errs.UserNotFound
			}

			removedMembers[i].UserGroup = memberGroup
			refs[1] = memberGroup
		}
	}

	existingMemberIds := model.FindArticleListMemberUserIdByArticleListIdTx(tx, listId)
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "remove_article_list_member",
		Public:       list.Public,
		Access:       model.FindArticleListMemberUserIdByArticleListIdTx(tx, listId),
		Notify:       util.RemoveIds(getMemberUserIds(tx, orga.ID, removedMembers), existingMemberIds),
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
