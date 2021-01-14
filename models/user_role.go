package models

import "gorm.io/gorm"

// UserRole model database representation
type UserRole struct {
	GormModelPK

	Type   string `gorm:"not null;index" json:"type"`
	UserID uint   `gorm:"not null;uniqueIndex" json:"user_id"`

	GormModelSfd
}

// UserRoleModel interface that implements the user_role model
type UserRoleModel interface {
	Create(*UserRole) error
}

func newUserRoleModel(db *gorm.DB) UserRoleModel {
	return &userRoleMode{
		db: db,
	}
}

type userRoleMode struct {
	db *gorm.DB
}

func (m *userRoleMode) Create(userRole *UserRole) error {
	return m.db.Create(userRole).Error
}
