package organization

import (
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

type Statistics struct {
	ArticleCount        int   `json:"article_count"`
	ListCount           int   `json:"list_count"`
	MemberCount         int   `json:"member_count"`
	BillableMemberCount int   `json:"billable_member_count"`
	GroupCount          int   `json:"group_count"`
	TagCount            int   `json:"tag_count"`
	StorageUsage        int64 `json:"storage_usage"`
	MaxStorage          int64 `json:"max_storage"`
}

func ReadOrganizations(userId hide.ID) []model.Organization {
	return model.FindOrganizationsByUserId(userId)
}

func ReadOrganization(ctx context.EmviContext) (*model.Organization, error) {
	var organization *model.Organization

	if ctx.IsUser() {
		organization = model.GetOrganizationByUserIdAndIdAndIsAdmin(ctx.Organization.ID, ctx.UserId)
	} else {
		organization = model.GetOrganizationById(ctx.Organization.ID)
	}

	if organization == nil {
		return nil, errs.OrganizationNotFound
	}

	return organization, nil
}

func GetOrganizationStatistics(orga *model.Organization, userId hide.ID) (*Statistics, error) {
	if _, err := perm.CheckUserIsAdminOrMod(orga.ID, userId); err != nil {
		return nil, err
	}

	return &Statistics{
		model.CountArticleByOrganizationId(orga.ID),
		model.CountArticleListByOrganizationId(orga.ID),
		model.CountOrganizationMemberByOrganizationIdAndActive(orga.ID),
		model.CountOrganizationMemberByOrganizationIdAndActiveAndNotReadOnly(orga.ID),
		model.CountUserGroupByOrganizationId(orga.ID),
		model.CountTagByOrganizationId(orga.ID),
		model.GetFileStorageUsageByOrganizationId(orga.ID),
		orga.MaxStorageGB,
	}, nil
}
