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
			Name: "CreateProject  Ecommerce",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"Ecommerce",
					"http://bitbucket",
					"roopa@gmail.com",
					"1",
					"2",
					"3",
					"5",
				)
				return nil
			},
		},
		{
			Name: "CreateProject ReactNative",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"ReactNative",
					"http://bitbucket",
					"roopa@gmail.com",
					"1",
					"2",
					"3",
					"5",
				)
				return nil
			},
		},
		{
			Name: "CreateProject React",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"React",
					"http://bitbucket",
					"roopa@gmail.com",
					"1",
					"2",
					"3",
					"5",
				)
				return nil
			},
		},
		{
			Name: "CreateProject Hybris",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"Hybris",
					"http://bitbucket",
					"roopa@gmail.com",
					"1",
					"2",
					"3",
					"5",
				)
				return nil
			},
		},
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
