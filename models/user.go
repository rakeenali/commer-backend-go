package models

import (
	"gorm.io/gorm"
)

// User db model
type User struct {
	// gorm.Model
	GormModelPK

	Username string      `gorm:"not null;uniqueIndex" json:"username"`
	Password string      `gorm:"not null;" json:"password"`
	Account  Accounts    `gorm:"foreignKey:UserID;references:ID" json:"account"`
	Role     UserRole    `gorm:"foreignKey:UserID;references:ID" json:"user_role"`
	Balance  UserBalance `gorm:"foreignKey:UserID;references:ID" json:"balance"`

	GormModelSfd
}

// userModel interface that must be implemented by services trying to use user
type userModel interface {
	Create(*User) error

	ByUsername(string) (*User, error)
	ByID(uint) (*User, error)
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
	err := ug.db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (ug *userGorm) ByUsername(username string) (*User, error) {
	var user User

	err := ug.db.Where("username = ?",
		username).Preload("Account").Preload("Role").Preload("Balance").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User

	err := ug.db.Where("id = ?",
		id).Preload("Account").Preload("Role").Preload("Balance").First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// st := ug.db.Session(&gorm.Session{DryRun: true}).Where("username = ?", username).First(&user).Statement
// fmt.Println(st.SQL.String())
