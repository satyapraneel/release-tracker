package models

type ReleaseProject struct {
	ID  		uint
	ReleaseId  	uint
	ProjectId	uint
	Release  	Release `gorm:"ForeignKey:ReleaseId;References:id"`
	Project 	Project `gorm:"ForeignKey:ProjectId;References:id"`
}
