package models

import (
	"commerce/helpers"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Tags database model implementation
type Tags struct {
	GormModelPK

	Name  string  `gorm:"not null;uniqueIndex" json:"name"`
	Items []Items `gorm:"many2many:items_tags;" json:"items"`

	GormModelSfd
}

// TagsModel interface that implements the tags model
type TagsModel interface {
	List() (*[]Tags, error)
	ByID(id uint) (*Tags, error)
	WithItems(id uint) (*Tags, error)
	ByName(name string) (*Tags, error)

	Create(*Tags) error
	Update(uint, *Tags) error
}

func newTagsModel(db *gorm.DB) TagsModel {
	return &tagsModel{
		db: db,
	}
}

type tagsModel struct {
	db *gorm.DB
}

func (tm *tagsModel) List() (*[]Tags, error) {
	var tags []Tags

	err := tm.db.Preload("Items").Preload("Items.Tags").Find(&tags).Error

	return &tags, err
}

func (tm *tagsModel) ByID(id uint) (*Tags, error) {
	var tag Tags
	err := first(tm.db.Where("id = ?", id), &tag)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (tm *tagsModel) ByName(name string) (*Tags, error) {
	var tag Tags
	err := first(tm.db.Where("name = ?", strings.ToLower(name)), &tag)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (tm *tagsModel) Create(t *Tags) error {
	t.Name = strings.ToLower(t.Name)
	err := tm.db.Create(t).Error

	fmt.Println("error", err)
	if err != nil {
		return helpers.ErrTagExist
	}
	return err
}

func (tm *tagsModel) Update(id uint, tag *Tags) error {
	var find Tags

	err := first(tm.db.Where("id", id), &find)
	if err != nil {
		return err
	}

	if tag.Name != "" {
		find.Name = strings.ToLower(tag.Name)
	}

	err = tm.db.Save(&find).Error
	if err != nil {
		return err
	}

	return nil
}

func (tm *tagsModel) WithItems(id uint) (*Tags, error) {
	var tags Tags
	err := tm.db.Where("id = ?",
		id).Preload("Items").Preload("Items.Tags").First(&tags).Error
	if err != nil {
		return nil, err
	}

	return &tags, nil
}
