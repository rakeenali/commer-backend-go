package models

import (
	"gorm.io/gorm"
)

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
	Delete(id uint) error
	ByUserID(uint) (*UserRole, error)
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

func (m *userRoleMode) ByUserID(id uint) (*UserRole, error) {
	var ur UserRole
	err := first(
		m.db.Where("user_id = ?", id),
		&ur,
	)
	if err != nil {
		return nil, err
	}

	return &ur, nil
}

func (m *userRoleMode) Delete(id uint) error {
	var ur UserRole
	ur.ID = id
	return m.db.Unscoped().Delete(&ur).Error
}
