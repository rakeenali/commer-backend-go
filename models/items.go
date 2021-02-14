package models

import (
	"commerce/helpers"

	"gorm.io/gorm"
)

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
	List(*[]Items, int, int) (int, error)
	ByID(uint) (*Items, error)
	Create(item *Items, tags []Tags) (*Items, error)
	Update(uint, *Items, []Tags) error
	Delete(id uint) error
}

func newItemsModel(db *gorm.DB) ItemsModel {
	return &itemsModel{
		db: db,
	}
}

type itemsModel struct {
	db *gorm.DB
}

func (im *itemsModel) List(items *[]Items, pageSize int, page int) (int, error) {
	var count int64

	im.db.Model(&Items{}).Count(&count)
	var err error
	offset := (page - 1) * pageSize

	if offset >= int(count) {
		page = 1
		err = im.db.Offset(0).Limit(pageSize).Preload("Tags").Find(items).Error
	} else {
		err = im.db.Offset(offset).Limit(pageSize).Preload("Tags").Find(items).Error
	}
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (im *itemsModel) ByID(id uint) (*Items, error) {
	var item Items
	err := im.db.Where("id = ?",
		id).Preload("Tags").First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, helpers.ErrNotFound
		}
		panic(err)
	}

	return &item, nil
}

func (im *itemsModel) Create(item *Items, tags []Tags) (*Items, error) {
	err := im.db.Create(item).Error
	if err != nil {
		return nil, err
	}

	err = im.db.Model(&item).Association("Tags").Append(&tags)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (im *itemsModel) Update(id uint, updatedItem *Items, tags []Tags) error {
	var item Items
	item.ID = id

	err := im.db.Model(&item).Updates(&updatedItem).Error
	if err != nil {
		return err
	}

	err = im.db.Model(&item).Association("Tags").Replace(&tags)
	if err != nil {
		return err
	}

	return nil
}

func (im *itemsModel) Delete(id uint) error {
	var item Items
	item.ID = id

	err := im.db.Delete(&item).Error
	if err != nil {
		return err
	}

	return nil
}
