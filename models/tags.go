package models

import "gorm.io/gorm"

// Tags database model implementation
type Tags struct {
	GormModelPK

	Name  string  `gorm:"not null;uniqueIndex" json:"name"`
	Items []Items `gorm:"many2many:items_tags;" json:"items"`

	GormModelSfd
}

// TagsModel interface that implements the tags model
type TagsModel interface {
	Create(*Tags) error
}

func newTagsModel(db *gorm.DB) TagsModel {
	return &tagsModel{
		db: db,
	}
}

type tagsModel struct {
	db *gorm.DB
}

func (tm *tagsModel) Create(t *Tags) error {
	return tm.db.Create(t).Error
}
