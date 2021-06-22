package models

import (
	"gorm.io/gorm"
)

//Project struct declaration
type Project struct {
	gorm.Model
	Name                 string `gorm:"type:varchar(100);index:,unique"`
	RepoName             string `json:"repo_name"`
	ReviewerList         string `json:"reviewer_list"`
	BetaReleaseDate      string `json:"beta_release_date"`
	RegressionSignorDate string `json:"regression_signor_date"`
	CodeFreezeDate       string `json:"code_freeze_date"`
	DevCompletionDate    string `json:"dev_completion_date"`
	Status               string `gorm:"type:varchar(1)" sql:"DEFAULT:1"`
}

type ProjectResult struct {
	Total    int64     `json:"recordsTotal"`
	Filtered int64     `json:"recordsFiltered"`
	Data     []Project `json:"data"`
}
