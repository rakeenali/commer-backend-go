package models

import "gorm.io/gorm"

// Items database model implementation
type Items struct {
	GormModelPK

	Name  string `gorm:"not null" json:"name"`
	Sku   string `gorm:"not null" json:"sku"`
	Image string `gorm:"not null" json:"image"`
	Price uint64 `gorm:"not null" json:"price"`
	Tags  []Tags `gorm:"many2many:items_tags;" json:"tags"`

	GormModelSfd
}

// ItemsModel interface that implements the items model
type ItemsModel interface {
}

func newItemsModel(db *gorm.DB) ItemsModel {
	return &itemsModel{
		db: db,
	}
}

type itemsModel struct {
	db *gorm.DB
}
