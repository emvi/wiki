package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type ArticleRecommendation struct {
	db.BaseEntity

	ArticleId     hide.ID `db:"article_id" json:"article_id"`
	UserId        hide.ID `db:"user_id" json:"user_id"`
	RecommendedTo hide.ID `db:"recommended_to" json:"recommended_to"`

	Member *OrganizationMember `db:"organization_member" json:"organization_member"` // the user who recommended the article
}

func GetArticleRecommendationByArticleIdAndUserIdAndRecommendedTo(articleId, userId, recommendTo hide.ID) *ArticleRecommendation {
	entity := new(ArticleRecommendation)

	if err := connection.Get(entity, `SELECT * FROM article_recommendation WHERE article_id = $1 AND user_id = $2 AND recommended_to = $3`, articleId, userId, recommendTo); err != nil {
		logbuch.Debug("Article recommendation by article id and user id and recommend to not found", logbuch.Fields{"err": err, "article_id": articleId, "user_id": userId, "recommend_to": recommendTo})
		return nil
	}

	return entity
}

func FindArticleRecommendationByArticleIdAndRecommendedTo(articleId, recommendTo hide.ID) []ArticleRecommendation {
	query := `SELECT * FROM article_recommendation WHERE article_id = $1 AND recommended_to = $2`
	var entities []ArticleRecommendation

	if err := connection.Select(&entities, query, articleId, recommendTo); err != nil {
		logbuch.Error("Article recommendation by article id and recommend to not found", logbuch.Fields{"err": err, "article_id": articleId, "recommend_to": recommendTo})
		return nil
	}

	return entities
}

func FindArticleRecommendationByArticleIdAndRecommendedToWithUser(articleId, recommendTo hide.ID) []ArticleRecommendation {
	query := `SELECT * FROM (
		SELECT DISTINCT ON (article_recommendation.id) article_recommendation.*,
		"user".id "organization_member.user.id",
		"user".firstname "organization_member.user.firstname",
		"user".lastname "organization_member.user.lastname",
		"user".picture "organization_member.user.picture",
		"organization_member".username "organization_member.username"
		FROM article_recommendation
		JOIN "user" ON article_recommendation.user_id = "user".id
		JOIN organization_member ON "user".id = organization_member.user_id
		WHERE article_id = $1
		AND recommended_to = $2
		) AS results ORDER BY "organization_member.user.lastname", "organization_member.user.firstname", "organization_member.username" ASC`
	var entities []ArticleRecommendation

	if err := connection.Select(&entities, query, articleId, recommendTo); err != nil {
		logbuch.Error("Article recommendation by article id and recommend to not found", logbuch.Fields{"err": err, "article_id": articleId, "recommend_to": recommendTo})
		return nil
	}

	return entities
}

func DeleteArticleRecommendationById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "article_recommendation" WHERE id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article recommendation by id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveArticleRecommendation(tx *sqlx.Tx, entity *ArticleRecommendation) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article_recommendation" (article_id, user_id, recommended_to)
			VALUES (:article_id, :user_id, :recommended_to)
			RETURNING id`,
		`UPDATE "article_recommendation" SET article_id = :article_id,
			user_id = :user_id,
			recommended_to = :recommended_to
			WHERE id = :id`)
}
