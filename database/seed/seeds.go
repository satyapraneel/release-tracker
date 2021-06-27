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
					"ecom",
					"roopa@gmail.com",
					"1",
					"2",
					"3",
					"5",
					"LMGID,LOYAL",
					"1",
				)
				return nil
			},
		},
		{
			Name: "CreateProject ReactNative",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"ReactNative",
					"react",
					"roopa@gmail.com",
					"1",
					"2",
					"3",
					"5",
					"LMGID,LOYAL",
					"1",
				)
				return nil
			},
		},
		{
			Name: "CreateProject React",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"React",
					"react",
					"roopa@gmail.com",
					"1",
					"2",
					"3",
					"5",
					"LMGID,LOYAL",
					"1",
				)
				return nil
			},
		},
		{
			Name: "CreateProject Hybris",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"Hybris",
					"hybris",
					"roopa@gmail.com",
					"1",
					"2",
					"3",
					"5",
					"LMGID,LOYAL",
					"1",
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

		{
			Name: "Create Reviewer",
			Run: func(db *gorm.DB) error {
				CreateReviwers(db,
					"Roopa J",
					"roopa.j@landmarkgroup.in",
					"roopa1118",
				)
				return nil
			},
		},
		{
			Name: "Create DLs User",
			Run: func(db *gorm.DB) error {
				CreateDls(db,
					"satyapraneel@gmail.com",
					"qa",
				)
				return nil
			},
		},
		{
			Name: "Create Reviewer 2",
			Run: func(db *gorm.DB) error {
				CreateReviwers(db,
					"Satya P H",
					"satyapraneelh@yahoo.com",
					"satyapraneelh",
				)
				return nil
			},
		},
		{
			Name: "Map DLs User",
			Run: func(db *gorm.DB) error {
				CreateDlsProject(db,
					1,
					1,
				)
				return nil
			},
		},
	}
}
