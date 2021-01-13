package models

import "gorm.io/gorm"

// Accounts gorm model
type Accounts struct {
	// gorm.Model
	GormModelPK
	FirstName string `gorm:"not null;" json:"first_name"`
	LastName  string `gorm:"not null;" json:"last_name"`
	UserID    string `gorm:"not null;uniqueIndex" json:"user_id"`
	GormModelSfd
}

// AccountsModel interface that implements the account model
type AccountsModel interface {
	Create(account *Accounts) error
}

// returns pointer to accounts model implementation
func newAccountsModel(db *gorm.DB) AccountsModel {
	return &accountsModel{
		db: db,
	}
}

type accountsModel struct {
	db *gorm.DB
}

func (am *accountsModel) Create(account *Accounts) error {
	return am.db.Create(&account).Error
}
