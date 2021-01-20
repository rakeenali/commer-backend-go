package models

import "gorm.io/gorm"

// UserBalance model
type UserBalance struct {
	GormModelPK

	Balance float32 `gorm:"not null;uniqueIndex" json:"balance"`
	UserID  uint    `gorm:"not null;uniqueIndex" json:"user_id"`

	GormModelSfd
}

// UserBalanceModel interface for models
type UserBalanceModel interface {
	Credit(*User, float32) (*UserBalance, error)
}

func newUserBalanceModel(db *gorm.DB) UserBalanceModel {
	return &userBalanceModel{
		db: db,
	}
}

type userBalanceModel struct {
	db *gorm.DB
}

func (ub *userBalanceModel) Credit(user *User, balance float32) (*UserBalance, error) {
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
