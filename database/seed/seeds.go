package seed

import (
	"github.com/release-trackers/gin/repositories"
	"gorm.io/gorm"
)

type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

func All() []Seed {
	return []Seed{
		{
			Name: "CreateUser Admin User",
			Run: func(db *gorm.DB) error {
				CreateUser(db,
					"admin",
					"admin@admin.com",
					repositories.PasswordHash("password"),
				)
				return nil
			},
		},
	}
}
