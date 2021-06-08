package seed

import "gorm.io/gorm"

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
	}
}
