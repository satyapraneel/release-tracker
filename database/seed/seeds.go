package seed

import "gorm.io/gorm"

func All() []Seed {
	return []Seed{
		Seed{
			Name: "CreateProject  Ecommerce",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"Ecommerce",
					"http://bitbucket",
					"roopa@gmail.com",
					"T-1",
					"T-2",
					"T-3",
					"T-5",
					)
				return nil
			},
		},
		Seed{
			Name: "CreateProject ReactNative",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"ReactNative",
					"http://bitbucket",
					"roopa@gmail.com",
					"T-1",
					"T-2",
					"T-3",
					"T-5",
				)
				return nil
			},
		},
		Seed{
			Name: "CreateProject React",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"React",
					"http://bitbucket",
					"roopa@gmail.com",
					"T-1",
					"T-2",
					"T-3",
					"T-5",
				)
				return nil
			},
		},
		Seed{
			Name: "CreateProject Hybris",
			Run: func(db *gorm.DB) error {
				CreateProject(db,
					"Hybris",
					"http://bitbucket",
					"roopa@gmail.com",
					"T-1",
					"T-2",
					"T-3",
					"T-5",
				)
				return nil
			},
		},
	}
}
