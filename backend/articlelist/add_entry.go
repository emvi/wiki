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

func AddArticleListEntry(organization *model.Organization, userId, listId hide.ID, articleIds []hide.ID) ([]model.ArticleListEntry, error) {
	_, err := checkListExists(organization, listId)

	if err != nil {
		return nil, err
	}

	if err := checkUserModAccess(listId, userId); err != nil {
		return nil, err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to add article list entries", logbuch.Fields{"err": err})
		return nil, errs.TxBegin
	}

	pos := model.GetArticleListEntryLastPositionByArticleListIdTx(tx, listId)
	langId := util.DetermineLang(tx, organization.ID, userId, 0).ID
	articles := make([]model.Article, 0, len(articleIds))
	entries := make([]model.ArticleListEntry, 0, len(articleIds))

	for _, articleId := range articleIds {
		if model.GetArticleListEntryByArticleListIdAndArticleIdTx(tx, listId, articleId) != nil {
			continue
		}

		article, err := checkUserArticleAccess(tx, organization.ID, articleId, userId)

		if err != nil {
			db.Rollback(tx)
			return nil, err
		}

		article.LatestArticleContent = model.GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageIdTx(tx, organization.ID, article.ID, langId, false)
		pos++
		entry := &model.ArticleListEntry{ArticleListId: listId,
			ArticleId: articleId,
			Article:   article,
			Position:  pos}

		if err := model.SaveArticleListEntry(tx, entry); err != nil {
			return nil, errs.Saving
		}

		articles = append(articles, *article)
		entries = append(entries, *entry)
	}

	if err := createAddEntryFeed(tx, organization, userId, listId, articles); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when adding article list entries", logbuch.Fields{"err": err})
		return nil, errs.TxCommit
	}

	return entries, nil
}

func checkUserArticleAccess(tx *sqlx.Tx, orgaId, articleId, userId hide.ID) (*model.Article, error) {
	article := model.GetArticleByOrganizationIdAndIdIgnoreArchivedTx(tx, orgaId, articleId)

	if article == nil {
		return nil, errs.ArticlePermissionDenied
	}

	if article.ReadEveryone {
		return article, nil
	} else {
		if model.FindArticleAccessByArticleIdAndUserIdTx(tx, articleId, userId) == nil {
			return nil, errs.ArticlePermissionDenied
		}
	}

	return article, nil
}

func createAddEntryFeed(tx *sqlx.Tx, orga *model.Organization, userId, listId hide.ID, articles []model.Article) error {
	list := model.GetArticleListByOrganizationIdAndIdTx(tx, orga.ID, listId)

	if list == nil {
		db.Rollback(tx)
		return errs.ArticleListNotFound
	}

	refs := make([]interface{}, 0, len(articles)+1)
	refs = append(refs, list)
	reason := "add_article_list_entry"

	for i := range articles {
		if articles[i].ReadEveryone {
			refs = append(refs, &articles[i])
		} else {
			reason = "add_protected_article_list_entry"
		}
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
		logbuch.Error("Error creating feed when adding entry to article list", logbuch.Fields{"err": err})
		return err
	}

	return nil
}
