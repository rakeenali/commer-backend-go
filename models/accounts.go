package models

import (
	"commerce/helpers"

	"gorm.io/gorm"
)

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
	Update(uint, *Accounts) (*Accounts, error)
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

func (am *accountsModel) Update(id uint, account *Accounts) (*Accounts, error) {
	var find Accounts
	find.ID = id

	err := am.db.First(&find).Error
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	if account.FirstName != "" {
		find.FirstName = account.FirstName
	}
	if account.LastName != "" {
		find.LastName = account.LastName
	}

	err = am.db.Save(find).Error
	if err != nil {
		return nil, helpers.ErrAccountUpdate
	}

	return &find, nil
}
