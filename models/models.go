package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ModelConfig every hook of model must satisfy this type
type ModelConfig func(*Models) error

// WithMysql will initialize connection with mysql
func WithMysql(connectionInfo string) ModelConfig {
	return func(m *Models) error {
		db, err := gorm.Open(mysql.Open(connectionInfo), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		m.db = db
		return nil
	}
}

// WithUserModel will initialize user model
func WithUserModel() ModelConfig {
	return func(m *Models) error {
		um := newUserModel(m.db)
		m.User = um
		return nil
	}
}

// WithAccountsModel will initialize user model
func WithAccountsModel() ModelConfig {
	return func(m *Models) error {
		am := newAccountsModel(m.db)
		m.Accounts = am
		return nil
	}
}

// NewModels will initialize models
func NewModels(fns ...ModelConfig) *Models {
	var m Models
	for _, f := range fns {
		err := f(&m)
		if err != nil {
			panic(err)
		}
	}

	return &m
}

// Models struct for db
type Models struct {
	User     userModel
	Accounts AccountsModel
	db       *gorm.DB
}

// ApplyMigration will drop all the tables provided and will then create new migrations
func (m *Models) ApplyMigration() {
	destroyTables(m)
	m.db.AutoMigrate(&User{}, &Accounts{})
}

func destroyTables(m *Models) {
	exist := m.db.Migrator().HasTable(&User{})
	if exist {
		err := m.db.Migrator().DropTable(&User{})
		if err != nil {
			panic(err)
		}
	}
	exist = m.db.Migrator().HasTable(&Accounts{})
	if exist {
		err := m.db.Migrator().DropTable(&Accounts{})
		if err != nil {
			panic(err)
		}
	}
}
