package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type SupportTicket struct {
	db.BaseEntity

	OrganizationId hide.ID `db:"organization_id" json:"organization_id"`
	UserId         hide.ID `db:"user_id" json:"user_id"`
	Type           string  `json:"type"`
	Subject        string  `json:"subject"`
	Message        string  `json:"message"`
	Status         string  `json:"status"`
}

func FindSupportTicketByOrganizationId(orgaId hide.ID) []SupportTicket {
	query := `SELECT * FROM "support_ticket" WHERE organization_id = $1`
	var entities []SupportTicket

	if err := connection.Select(&entities, query, orgaId); err != nil {
		logbuch.Error("Error reading support tickets by organization id", logbuch.Fields{"err": err, "orga_id": orgaId})
		return nil
	}

	return entities
}

func SaveSupportTicket(tx *sqlx.Tx, entity *SupportTicket) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "support_ticket" (organization_id,
			user_id,
			type,
			subject,
			message,
			status)
			VALUES (:organization_id,
			:user_id,
			:type,
			:subject,
			:message,
			:status) RETURNING id`,
		`UPDATE "support_ticket" SET organization_id = :organization_id,
			user_id = :user_id,
			type = :type,
			subject = :subject,
			message = :message,
			status = :status
			WHERE id = :id`)
}
