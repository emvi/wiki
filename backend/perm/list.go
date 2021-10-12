package perm

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/jmoiron/sqlx"
)

// CheckUserListAccess checks the given user has direct access or indirect access to through a group to given article list.
// Returns an error if he doesn't. The transaction is optional.
// It does not check if the list is public!
func CheckUserListAccess(tx *sqlx.Tx, listId, userId hide.ID) error {
	if len(model.FindArticleListMemberByArticleListIdAndUserIdIncludingUserGroupMember(tx, listId, userId)) == 0 {
		return errs.PermissionDenied
	}

	return nil
}
