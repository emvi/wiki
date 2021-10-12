package usergroup

import (
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/jmoiron/sqlx"
)

func GetAllGroup(tx *sqlx.Tx, orgaId hide.ID) *model.UserGroup {
	return model.GetUserGroupByOrganizationIdAndNameTx(tx, orgaId, constants.GroupAllName)
}

func GetAdminGroup(tx *sqlx.Tx, orgaId hide.ID) *model.UserGroup {
	return model.GetUserGroupByOrganizationIdAndNameTx(tx, orgaId, constants.GroupAdminName)
}

func GetModGroup(tx *sqlx.Tx, orgaId hide.ID) *model.UserGroup {
	return model.GetUserGroupByOrganizationIdAndNameTx(tx, orgaId, constants.GroupModName)
}

func GetReadOnlyGroup(tx *sqlx.Tx, orgaId hide.ID) *model.UserGroup {
	return model.GetUserGroupByOrganizationIdAndNameTx(tx, orgaId, constants.GroupReadOnlyName)
}
