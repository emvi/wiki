package articlelist

import (
	articleutil "emviwiki/backend/article/util"
	"emviwiki/backend/bookmark"
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/backend/observe"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
)

const (
	privateListsLimit       = 10
	articlelistEntriesLimit = 20
	articlelistMemberLimit  = 20
)

func ReadPrivateArticleLists(orga *model.Organization, userId hide.ID, offset int) []model.ArticleList {
	langId := util.DetermineLang(nil, orga.ID, userId, 0).ID
	return model.FindArticleListByOrganizationIdAndUserIdAndLanguageIdAndPrivateWithLimit(orga.ID, userId, langId, offset, privateListsLimit)
}

func ReadArticleList(ctx context.EmviContext, langId, listId hide.ID) (*model.ArticleList, bool, bool, bool, error) {
	list := model.GetArticleListByOrganizationIdAndUserIdAndId(ctx.Organization.ID, ctx.UserId, listId)

	if list == nil {
		return nil, false, false, false, errs.ArticleListNotFound
	}

	if err := checkUserOrClientListAccess(ctx.UserId, list); err != nil {
		return nil, false, false, false, err
	}

	list.Names = model.FindArticleListNamesByArticleListId(listId)
	langId = util.DetermineLang(nil, ctx.Organization.ID, ctx.UserId, langId).ID
	list.Name = determineName(list.Names, ctx.Organization.ID, langId)
	isMod := len(model.FindArticleListMemberModeratorByArticleListIdAndUserId(listId, ctx.UserId)) != 0
	isObserved := observe.IsObserved(ctx.UserId, 0, listId, 0)
	isBookmarked := bookmark.IsBookmarked(ctx.UserId, 0, listId)
	return list, isMod, isObserved, isBookmarked, nil
}

func ReadArticleListEntries(ctx context.EmviContext, langId, listId hide.ID, filter *model.SearchArticleListEntryFilter) ([]model.Article, int, int, error) {
	list := model.GetArticleListByOrganizationIdAndId(ctx.Organization.ID, listId)

	if list == nil {
		return nil, 0, 0, errs.ArticleListNotFound
	}

	if err := checkUserOrClientListAccess(ctx.UserId, list); err != nil {
		return nil, 0, 0, err
	}

	if filter == nil {
		filter = new(model.SearchArticleListEntryFilter)
	}

	if filter.Limit > articlelistEntriesLimit {
		filter.Limit = articlelistEntriesLimit
	}

	pos := readArticleListSetCenter(list.ID, filter)
	filter.ClientAccess = ctx.IsClient()
	langId = util.DetermineLang(nil, ctx.Organization.ID, ctx.UserId, langId).ID
	results := model.FindArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdLimit(ctx.Organization.ID, ctx.UserId, langId, listId, filter)
	resultCount := model.CountArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdLimit(ctx.Organization.ID, ctx.UserId, langId, listId, filter)
	articleutil.RemoveAuthorsOrAuthorMails(ctx, results)
	return results, resultCount, pos, nil
}

func readArticleListSetCenter(listId hide.ID, filter *model.SearchArticleListEntryFilter) int {
	if filter.CenterArticleId != 0 {
		entry := model.GetArticleListEntryByArticleListIdAndArticleId(listId, filter.CenterArticleId)

		if entry == nil {
			return 1
		}

		pos := model.CountArticleListEntryByArticleListIdAndPositionBefore(listId, entry.Position)

		if filter.CenterBefore < 0 {
			filter.CenterBefore = 0
		}

		filter.Offset = pos - filter.CenterBefore
		offsetLimit := pos - articlelistEntriesLimit/2

		if offsetLimit < 0 {
			offsetLimit = 0
		}

		if filter.Offset < offsetLimit {
			filter.Offset = offsetLimit
		}

		if filter.Limit > articlelistEntriesLimit/2 {
			filter.Limit = articlelistEntriesLimit / 2
		}

		return filter.Offset + 1
	}

	return 1
}

func ReadArticleListMember(orga *model.Organization, userId, listId hide.ID, filter *model.SearchArticleListMemberFilter) ([]model.ArticleListMember, int, error) {
	list := model.GetArticleListByOrganizationIdAndId(orga.ID, listId)

	if list == nil {
		return nil, 0, errs.ArticleListNotFound
	}

	if !list.Public {
		if err := perm.CheckUserListAccess(nil, listId, userId); err != nil {
			return nil, 0, err
		}
	}

	if filter == nil {
		filter = new(model.SearchArticleListMemberFilter)
	}

	filter.Limit = articlelistMemberLimit
	return model.FindArticleListMemberByOrganizationIdAndArticleListIdAndFilterLimit(orga.ID, listId, filter),
		model.CountArticleListMemberByOrganizationIdAndArticleListIdAndFilterLimit(orga.ID, listId, filter), nil
}

func checkUserOrClientListAccess(userId hide.ID, list *model.ArticleList) error {
	if userId == 0 && !list.ClientAccess {
		return errs.PermissionDenied
	} else if userId != 0 && !list.Public {
		if err := perm.CheckUserListAccess(nil, list.ID, userId); err != nil {
			return err
		}
	}

	return nil
}

func determineName(names []model.ArticleListName, orgaId, langId hide.ID) *model.ArticleListName {
	for _, name := range names {
		if name.LanguageId == langId {
			return &name
		}
	}

	defaultLangId := model.GetDefaultLanguageByOrganizationId(orgaId).ID

	for _, name := range names {
		if name.LanguageId == defaultLangId {
			return &name
		}
	}

	if len(names) > 0 {
		return &names[0]
	}

	return nil
}
