package models

import (
	"gorm.io/gorm"
)

//Project struct declaration
type Project struct {
	gorm.Model
	Name                 string `gorm:"type:varchar(100);index:,unique"`
	BitbucketUrl         string `json:"bitbucket_url"`
	ReviewerList         string `json:"reviewer_list"`
	BetaReleaseDate      string `json:"beta_release_date"`
	RegressionSignorDate string `json:"regression_signor_date"`
	CodeFreezeDate       string `json:"code_freeze_date"`
	DevCompletionDate    string `json:"dev_completion_date"`
}

type ProjectResult struct {
	Total    int64     `json:"recordsTotal"`
	Filtered int64     `json:"recordsFiltered"`
	Data     []Project `json:"data"`
}
