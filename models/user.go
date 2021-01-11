package models

import "gorm.io/gorm"

// User db model
type User struct {
	GormModelPK
	Username string `gorm:"not null;uniqueIndex" json:"username"`
	Password string `gorm:"not null;" json:"password"`
	GormModelSfd
}

// userModel interface that must be implemented by services trying to use user
type userModel interface {
	Create(*User) error

	ByUsername(username string) (*User, error)
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

func (ug *userGorm) ByUsername(username string) (*User, error) {
	var user User
	err := first(ug.db.Where("username = ?", username), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
