package models

import (
	"gorm.io/gorm"
)

// Orders database model implementation
type Orders struct {
	GormModelPK

	Charge  uint64  `gorm:"not null" json:"charge"`
	Address string  `gorm:"not null" json:"address"`
	UserID  uint    `gorm:"not null" json:"user_id"`
	User    User    `gorm:"references:ID" json:"user"`
	Items   []Items `gorm:"many2many:order_items;" json:"items"`

	GormModelSfd
}

func newOrdersModel(db *gorm.DB) OrdersModel {
	return &ordersModel{
		db: db,
	}
}

// OrdersModel wil implements the order
type OrdersModel interface {
	List(userID uint) (*[]Orders, error)
	Create(*Orders, *[]Items) error
}

type ordersModel struct {
	db *gorm.DB
}

func (om *ordersModel) List(userID uint) (*[]Orders, error) {
	var orders []Orders
	err := om.db.Preload("Items").Preload("Items.Tags").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return &orders, nil
}

func (om *ordersModel) Create(order *Orders, items *[]Items) error {
	err := om.db.Create(order).Error
	if err != nil {
		return err
	}

	err = om.db.Model(order).Association("Items").Append(items)
	if err != nil {
		return err
	}

	return nil
}
