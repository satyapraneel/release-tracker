package cmd

import "gorm.io/gorm"

type Application struct {
	Db   *gorm.DB
	Name string
}

type App struct {
	*Application
}
