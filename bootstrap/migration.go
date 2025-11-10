package bootstrap

import (
	"log"

	"github.com/gabrielfmcoelho/platform-core/domain"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		// domains like: &domain.User{},
		&domain.User{},
		&domain.UserLog{},
		&domain.UserServiceLog{},
		&domain.Service{},
		&domain.ContactIntent{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
}
