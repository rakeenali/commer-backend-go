package models

import "gorm.io/gorm"

// User db model
type User struct {
	gorm.Model

	Username string `gorm:"not null;uniqueIndex"`
	Password string `gorm:"not null;"`
}

// userModel interface that must be implemented by services trying to use user
type userModel interface {
	Create(*User) error
}

// newUserModel will initialize the user model
func newUserModel(db *gorm.DB) userModel {
	return &userGorm{
		db: db,
	}
}

// UserGorm struct
type userGorm struct {
	db *gorm.DB
}

var _ userModel = &userGorm{}

func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}
