package models

import (
	"errors"

	"gorm.io/gorm"
)

// UserBalance model
type UserBalance struct {
	GormModelPK

	Balance uint64 `gorm:"not null;uniqueIndex" json:"balance"`
	UserID  uint   `gorm:"not null;uniqueIndex" json:"user_id"`

	GormModelSfd
}

// UserBalanceModel interface for models
type UserBalanceModel interface {
	Credit(*User, uint64) (*UserBalance, error)
	Debit(*User, uint64) (*UserBalance, error)
}

func newUserBalanceModel(db *gorm.DB) UserBalanceModel {
	return &userBalanceModel{
		db: db,
	}
}

type userBalanceModel struct {
	db *gorm.DB
}

func (ub *userBalanceModel) Credit(user *User, balance uint64) (*UserBalance, error) {
	var userBalance UserBalance
	if user.Balance.ID == 0 {
		userBalance.Balance = balance
		userBalance.UserID = user.ID
		err := ub.db.Create(&userBalance).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := first(ub.db.Where("user_id = ?", user.ID), &userBalance)
		if err != nil {
			return nil, err
		}
		userBalance.Balance = userBalance.Balance + balance
		err = ub.db.Save(&userBalance).Error
		if err != nil {
			return nil, err
		}
	}

	return &userBalance, nil
}

func (ub *userBalanceModel) Debit(user *User, balance uint64) (*UserBalance, error) {
	var userBalance UserBalance
	if user.Balance.ID == 0 {
		return nil, errors.New("Cannot debit a balance from a user that does not exist")
	}

	err := first(ub.db.Where("user_id = ?", user.ID), &userBalance)
	if err != nil {
		return nil, err
	}
	userBalance.Balance = userBalance.Balance - balance
	err = ub.db.Save(&userBalance).Error
	if err != nil {
		return nil, err
	}

	return &userBalance, nil
}
