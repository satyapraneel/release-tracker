package models

type ReleaseTickets struct {
	ID  		uint
	Key  		string
	Type		string
	Summary 	string
	Project 	string
	Status      string
	ReleaseId  	uint
	Release  	Release `gorm:"ForeignKey:ReleaseId;References:id"`
}
