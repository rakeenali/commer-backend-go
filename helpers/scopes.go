package helpers

import (
	"gorm.io/gorm"
)

// Paginate can be use a scope for gorm  to paginate models
func Paginate(count *int, pageSize *int, page *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		switch {
		case *pageSize > 13:
			*pageSize = 13
		case *pageSize <= 0:
			*pageSize = 10
		}

		offset := (*page - 1) * *pageSize
		if offset >= *count {
			*page = 1
			return db.Offset(0).Limit(*pageSize)
		}
		return db.Offset(offset).Limit(*pageSize)

	}
}
