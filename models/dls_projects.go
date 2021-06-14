package models

type DlsProjects struct {
	ID        uint
	DlsId     uint
	ProjectId uint
	Dls       DLS     `gorm:"ForeignKey:DlsId;References:id"`
	Project   Project `gorm:"ForeignKey:ProjectId;References:id"`
}
