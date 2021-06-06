package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func QueryOffset(c *gin.Context) int {
	offset := c.Request.PostForm.Get("start")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}
	return offsetInt
}

func QueryLimit(c *gin.Context) int {
	limit :=  c.Request.PostForm.Get("length")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 25
	}
	return limitInt
}

func QueryOrder(c *gin.Context, columnOrder string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db
		if columnOrder != "" {
			query = query.Order("id "+columnOrder)
		}
		return query
	}
}

func SearchScope(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db
		search := c.PostFormMap("search")
		if search["value"] != "" {
			query = query.Where("name LIKE ? ", "%"+search["value"]+"%")
		}
		return query
	}
}