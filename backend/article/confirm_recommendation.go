package article

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

func ConfirmRecommendation(orga *model.Organization, userId, articleId hide.ID, sendNotificationUserIds []hide.ID) error {
	if _, err := checkUserReadAccess(orga.ID, userId, articleId); err != nil {
		return err
	}

	recommendations := model.FindArticleRecommendationByArticleIdAndRecommendedTo(articleId, userId)

	if len(recommendations) == 0 {
		return errs.RecommendationsNotFound
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to confirm recommendations", logbuch.Fields{"err": err, "orga_id": orga.ID, "article_id": articleId, "user_id": userId})
		return errs.TxBegin
	}

	notifyUsers := make([]hide.ID, 0, len(recommendations))

	for _, recommendation := range recommendations {
		if err := model.DeleteArticleRecommendationById(tx, recommendation.ID); err != nil {
			db.Rollback(tx)
			logbuch.Error("Error deleting article recommendation", logbuch.Fields{"err": err, "orga_id": orga.ID, "article_id": articleId, "user_id": userId, "recommendation_id": recommendation.ID})
			return errs.Saving
		}

		notifyUsers = append(notifyUsers, recommendation.UserId)
	}

	// the user can opt-out to notify the creator of the recommendations
	if err := createConfirmRecommendationReadFeed(tx, orga, userId, articleId, notifyUsers, sendNotificationUserIds); err != nil {
		db.Rollback(tx)
		logbuch.Error("Error creating feed while confirming recommendations", logbuch.Fields{"err": err, "orga_id": orga.ID, "article_id": articleId, "user_id": userId})
		return errs.Saving
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction to confirm recommendations", logbuch.Fields{"err": err, "orga_id": orga.ID, "article_id": articleId, "user_id": userId})
		return errs.TxCommit
	}

	return nil
}

func createConfirmRecommendationReadFeed(tx *sqlx.Tx, orga *model.Organization, recommendedTo, articleId hide.ID, notifyUsers, sendNotificationUserIds []hide.ID) error {
	article := model.GetArticleByOrganizationIdAndIdTx(tx, orga.ID, articleId)

	if article == nil {
		db.Rollback(tx)
		return errs.ArticleNotFound
	}

	langId := util.DetermineLang(tx, orga.ID, recommendedTo, 0).ID
	content := model.GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageIdTx(tx, orga.ID, articleId, langId, false)

	if content == nil {
		db.Rollback(tx)
		return errs.FindingLatestArticleContent
	}

	notifyUserIds := util.IntersectIds(notifyUsers, sendNotificationUserIds)

	if len(notifyUserIds) > 0 {
		feedData := &feed.CreateFeedData{Tx: tx, Organization: orga,
			UserId: recommendedTo,
			Reason: "recommendation_confirmation",
			Public: false,
			Access: []hide.ID{},
			Notify: notifyUserIds,
			Refs:   []interface{}{article, content}}

		if err := feed.CreateFeed(feedData); err != nil {
			return err
		}
	}

	return nil
}
