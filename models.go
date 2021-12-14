package main

import (
	"time"
	//	"github.com/jinzhu/gorm"
	//	"github.com/shopspring/decimal"
)

/////////////////
//models
//////////////////

type GormCustomModel struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	//	UpdatedAt time.Time  `json:"-"`
	UpdatedAt time.Time
	DeletedAt *time.Time `json:"-";sql:"index"`
}
type Config struct {
	GormCustomModel
	Key   string
	Value string
}

//    "Date": "2017-01-08T00:00:00Z", DATE FORMAT
type Position struct {
	GormCustomModel

	PositionName                string
	PositionCode                string
	PositionInmediateSuperior   string
	PositionInmediateSuperiorID uint
	AdvisingAuthority           bool
	Location                    string
	Abbreviation                string
	Company                     string
	DedicationRegime            string
	CreationDate                string
	SpacesToSupervisor          uint
	PositionObjective           string
	PositionPurpose             string
	PositionValues              string
	Clients                     string
	Products                    string
	Notes                       string
	AttachmentsName             string
}

type Contact struct {
	GormCustomModel
	Name string
}
type ContactRight struct {
	GormCustomModel
	ContactID uint
	Itemtype  string // tree node, sheet node, sheet
	Optype    string //additional,exception
	See       bool
	Modify    bool
	Del       bool
}

type Project struct {
	GormCustomModel
	ProjectName string `gorm:"unique_index:project_name"`
	Data        string
	//	UserID      int `gorm:"unique_index:project_name"`
}

type Sheet struct {
	GormCustomModel
	ProjectID                         uint `gorm:"unique_index:sheet_name"`
	PositionName                      string
	PositionAdvisors                  int
	DirectSupervisedPeople            uint
	DirectAndIndirectSupervisedPeople uint
	SheetName                         string `gorm:"unique_index:sheet_name"`
	Data                              string
	Attrs                             string
	//	UserID                            int `gorm:"unique_index:sheet_name"`
}
