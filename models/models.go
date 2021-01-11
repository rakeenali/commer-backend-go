package models

import (
	"commerce/helpers"
	"database/sql"
	"time"

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
	User userModel
	db   *gorm.DB
}

// ApplyMigration will drop all the tables provided and will then create new migrations
func (m *Models) ApplyMigration() {
	exist := m.db.Migrator().HasTable(&User{})
	if exist {
		err := m.db.Migrator().DropTable(&User{})
		if err != nil {
			panic(err)
		}
	}

	m.db.AutoMigrate(&User{})
}

// first will query using the provided gorm.DB and will either set first
// row returned on dest or will return error
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return helpers.ErrNotFound
	}
	return err
}

type GormModelPK struct {
	ID uint `gorm:"primarykey" json:"id"`
}
type GormModelSfd struct {
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`
}
