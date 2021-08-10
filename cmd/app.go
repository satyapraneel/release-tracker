package cmd

import (
	"github.com/jinzhu/gorm"
)

type Application struct {
	Db   *gorm.DB
	Name string
}

type App struct {
	*Application
}
