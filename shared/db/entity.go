package db

import (
	"github.com/emvi/hide"
	"time"
)

// Entity is an identifiable database entity.
type Entity interface {
	// GetId must return the ID of the entity.
	GetId() hide.ID

	// SetId must set the ID of the entity.
	SetId(hide.ID)
}

// BaseEntity is the base for all database entities.
type BaseEntity struct {
	ID      hide.ID   `json:"id"`
	DefTime time.Time `db:"def_time" json:"def_time"`
	ModTime time.Time `db:"mod_time" json:"mod_time"`
}

// GetId returns the ID of this entity.
func (entity *BaseEntity) GetId() hide.ID {
	return entity.ID
}

// SetId sets the ID to this entity.
func (entity *BaseEntity) SetId(id hide.ID) {
	entity.ID = id
}
