package models

import (
	"commerce/helpers"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// first will query using the provided gorm.DB and will either set first
// row returned on dest or will return error
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return helpers.ErrNotFound
	}
	return err
}

// GormModelPK embed struct for primary key model
type GormModelPK struct {
	ID uint `gorm:"primarykey" json:"id"`
}

// GormModelSfd embed struct for implementing softdeletes where needed
type GormModelSfd struct {
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`
}
