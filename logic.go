package main

import (
	"fmt"
	//	"time"

	//	"github.com/jinzhu/gorm"
	//	_ "github.com/jinzhu/gorm/dialects/sqlite"
	//	"gorm.io/driver/postgres"
	"errors"

	"gorm.io/gorm"
)

////////////////
// business logic
////////

////////

//project
func projectGetAll(db *gorm.DB, searchKey string) ([]Project, string) {
	var records []Project
	var err string
	result := db.Where("project_name LIKE ? ", "%"+searchKey+"%").Find(&records) // find all
	if result.Error != nil {
		err = result.Error.Error() + " - projectGetAll(), project"
	}
	return records, err
}
func projectGet(db *gorm.DB, id string) ([]Project, string) {
	var data = []Project{}
	var record Project
	var err string
	result := db.First(&record, id)
	if result.Error != nil {
		err = result.Error.Error() + " - projectGet(), project"
	}
	if record.ID != 0 {
		data = append(data, record)
	}
	return data, err
}
func projectCreate(db *gorm.DB, recordVal Project) (err error) {
	// check for name exists
	var record Project
	db.Where("project_name=?", recordVal.ProjectName).First(&record)
	if record.ID != 0 {
		fmt.Println("Error: Project Exists")
		//return result.Error
		return errors.New("UNIQUE constraint failed")
	}

	result := db.Create(&Project{
		ProjectName: recordVal.ProjectName,
	})
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
	}
	return result.Error
}
func projectUpdate(db *gorm.DB, recordVal Project, id string, isNameUpdate string) (err error) {
	// check for name exists
	if isNameUpdate != "" {
		var recordExists Project
		db.Where("project_name=?", recordVal.ProjectName).First(&recordExists)
		if recordExists.ID != 0 {
			fmt.Println("Error: Project Exists")
			return errors.New("UNIQUE constraint failed")
		}

	}

	var record Project
	result := db.First(&record, id)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
		return result.Error
	}
	record.ProjectName = recordVal.ProjectName
	record.Data = recordVal.Data

	result = db.Save(&record)
	return result.Error
}
func projectDelete(db *gorm.DB, id string) (err error) {
	//delete all project sheets
	resultSheet := db.Where("project_id = ?", id).Unscoped().Delete(&Sheet{})
	if resultSheet.Error != nil {
		fmt.Println("Error: ", resultSheet.Error.Error())
		return resultSheet.Error
	}
	//	result = db.Delete(&record)

	var record Project
	result := db.First(&record, id)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
		return result.Error
	}
	//	result = db.Delete(&record)
	result = db.Unscoped().Delete(&record)

	return result.Error
}

//contactRight

func contactRightGetAll(db *gorm.DB, contactID int, itemtype string, optype string) ([]ContactRight, string) {
	var records []ContactRight
	var err string
	result := db.Where("contact_id = ? and itemtype=? and optype=? ", contactID, itemtype, optype).Find(&records) // find all
	if result.Error != nil {
		err = result.Error.Error() + " - contactRightGetAll(), contactRight"
	}
	return records, err
}
func contactGetAll(db *gorm.DB) ([]Contact, string) {
	var records []Contact
	var err string
	result := db.Where("").Find(&records) // find all
	if result.Error != nil {
		err = result.Error.Error() + " - contactRightGetAll(), contactRight"
	}
	return records, err
}
func contactRightGet(db *gorm.DB, id string) ([]ContactRight, string) {
	var data = []ContactRight{}
	var record ContactRight
	var err string
	result := db.First(&record, id)
	if result.Error != nil {
		err = result.Error.Error() + " - contactRightGet(), contactRight"
	}
	if record.ID != 0 {
		data = append(data, record)
	}
	return data, err
}
func contactRightCreate(db *gorm.DB, recordVal ContactRight) (err error) {
	result := db.Create(&ContactRight{
		ContactID: recordVal.ContactID,
		Itemtype:  recordVal.Itemtype,
		Optype:    recordVal.Optype,
		See:       recordVal.See,
		Modify:    recordVal.Modify,
		Del:       recordVal.Del,
	})
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
	}
	return result.Error
}
func contactCreate(db *gorm.DB, recordVal Contact) (err error) {
	result := db.Create(&Contact{
		Name: recordVal.Name,
	})
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
	}
	return result.Error
}
func contactRightUpdate(db *gorm.DB, recordVal ContactRight, id string) (err error) {
	var record ContactRight
	result := db.First(&record, id)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
		return result.Error
	}
	record.ContactID = recordVal.ContactID
	record.Itemtype = recordVal.Itemtype
	record.Optype = recordVal.Optype
	record.See = recordVal.See
	record.Modify = recordVal.Modify
	record.Del = recordVal.Del
	result = db.Save(&record)
	return result.Error
}
func contactUpdate(db *gorm.DB, recordVal Contact, id string) (err error) {
	var record Contact
	result := db.First(&record, id)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
		return result.Error
	}
	record.Name = recordVal.Name
	result = db.Save(&record)
	return result.Error
}
func contactRightDelete(db *gorm.DB, id string) (err error) {
	var record ContactRight
	result := db.First(&record, id)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
		return result.Error
	}
	result = db.Delete(&record)
	return result.Error
}
func contactDelete(db *gorm.DB, id string) (err error) {
	var record Contact
	result := db.First(&record, id)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
		return result.Error
	}
	result = db.Delete(&record)
	return result.Error
}

//position
func positionGetAll(db *gorm.DB, searchKey string) ([]Position, string) {
	var records []Position
	var err string
	result := db.Where("position_name LIKE ? OR position_code LIKE ?", "%"+searchKey+"%", "%"+searchKey+"%").Find(&records) // find all

	if result.Error != nil {
		err = result.Error.Error() + " - positionGetAll(), position"
	}
	return records, err
}
func positionGet(db *gorm.DB, id string) ([]Position, string) {
	var data = []Position{}
	var record Position
	var err string
	result := db.First(&record, id)
	if result.Error != nil {
		err = result.Error.Error() + " - positionGet(), position"
	}
	if record.ID != 0 {
		data = append(data, record)
	}
	return data, err
}
func positionCreate(db *gorm.DB, recordVal Position) (err error) {
	result := db.Create(
		&Position{

			PositionName:                recordVal.PositionName,
			PositionCode:                recordVal.PositionCode,
			PositionInmediateSuperior:   recordVal.PositionInmediateSuperior,
			PositionInmediateSuperiorID: recordVal.PositionInmediateSuperiorID,
			AdvisingAuthority:           recordVal.AdvisingAuthority,
			Location:                    recordVal.Location,
			Abbreviation:                recordVal.Abbreviation,
			Company:                     recordVal.Company,
			DedicationRegime:            recordVal.DedicationRegime,
			CreationDate:                recordVal.CreationDate,
			SpacesToSupervisor:          recordVal.SpacesToSupervisor,
			PositionObjective:           recordVal.PositionObjective,
			PositionPurpose:             recordVal.PositionPurpose,
			PositionValues:              recordVal.PositionValues,
			Clients:                     recordVal.Clients,
			Products:                    recordVal.Products,
			Notes:                       recordVal.Notes,
			AttachmentsName:             recordVal.AttachmentsName})

	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
	}
	return result.Error
}
func positionUpdate(db *gorm.DB, recordVal Position, id string) (err error) {
	var record Position
	result := db.First(&record, id)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
		return result.Error
	}

	record.PositionName = recordVal.PositionName
	record.PositionCode = recordVal.PositionCode
	record.PositionInmediateSuperior = recordVal.PositionInmediateSuperior
	record.PositionInmediateSuperiorID = recordVal.PositionInmediateSuperiorID
	record.AdvisingAuthority = recordVal.AdvisingAuthority
	record.Location = recordVal.Location
	record.Abbreviation = recordVal.Abbreviation
	record.Company = recordVal.Company
	record.DedicationRegime = recordVal.DedicationRegime
	record.CreationDate = recordVal.CreationDate
	record.SpacesToSupervisor = recordVal.SpacesToSupervisor
	record.PositionObjective = recordVal.PositionObjective
	record.PositionPurpose = recordVal.PositionPurpose
	record.PositionValues = recordVal.PositionValues
	record.Clients = recordVal.Clients
	record.Products = recordVal.Products
	record.Notes = recordVal.Notes
	record.AttachmentsName = recordVal.AttachmentsName

	result = db.Save(&record)
	return result.Error
}
func positionDelete(db *gorm.DB, id string) (err error) {
	var record Position
	result := db.First(&record, id)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
		return result.Error
	}
	result = db.Delete(&record)
	return result.Error
}

//sheet
func sheetGetAll(db *gorm.DB, searchKey string) ([]Sheet, string) {
	var records []Sheet
	var err string
	result := db.Where("sheet_name LIKE ?", "%"+searchKey+"%").Find(&records) // find all
	//result := db.Table("positions").Select("*").Where("name LIKE ?", "%"+searchKey+"%").Scan(&records)

	if result.Error != nil {
		err = result.Error.Error() + " - sheetGetAll(), position"
	}
	return records, err
}
func sheetGet(db *gorm.DB, id string) ([]Sheet, string) {
	var data = []Sheet{}
	var record Sheet
	var err string
	result := db.First(&record, id)
	if result.Error != nil {
		err = result.Error.Error() + " - sheetGet(), position"
	}
	if record.ID != 0 {
		data = append(data, record)
	}
	return data, err
}

func sheetCreate(db *gorm.DB, recordVal Sheet) (err error) {

	// check for name exists
	var record Sheet
	db.Where("sheet_name=? and project_id=?", recordVal.SheetName, recordVal.ProjectID).First(&record)
	if record.ID != 0 {
		fmt.Println("Error: Sheet Exists")
		//return result.Error
		return errors.New("UNIQUE constraint failed")
	}

	result := db.Create(&Sheet{
		PositionName:                      recordVal.PositionName,
		PositionAdvisors:                  recordVal.PositionAdvisors,
		DirectSupervisedPeople:            recordVal.DirectSupervisedPeople,
		DirectAndIndirectSupervisedPeople: recordVal.DirectAndIndirectSupervisedPeople,
		SheetName:                         recordVal.SheetName,
		ProjectID:                         recordVal.ProjectID})

	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
	}
	return result.Error
}
func sheetUpdate(db *gorm.DB, recordVal Sheet, id string, isNameUpdate string) (err error) {

	// check for name exists
	if isNameUpdate != "" {
		// check for name exists
		var record Sheet
		db.Where("sheet_name=? and project_id=?", recordVal.SheetName, recordVal.ProjectID).First(&record)
		if record.ID != 0 {
			fmt.Println("Error: Sheet Exists")
			//return result.Error
			return errors.New("UNIQUE constraint failed")
		}

	}

	var record Sheet
	result := db.First(&record, id)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
		return result.Error
	}
	record.PositionName = recordVal.PositionName
	record.PositionAdvisors = recordVal.PositionAdvisors
	record.DirectSupervisedPeople = recordVal.DirectSupervisedPeople
	record.DirectAndIndirectSupervisedPeople = recordVal.DirectAndIndirectSupervisedPeople
	record.SheetName = recordVal.SheetName
	record.ProjectID = recordVal.ProjectID
	record.Data = recordVal.Data
	record.Attrs = recordVal.Attrs
	result = db.Save(&record)
	return result.Error
}
func sheetDelete(db *gorm.DB, id string) (err error) {
	var record Sheet
	result := db.First(&record, id)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error.Error())
		return result.Error
	}
	//	result = db.Delete(&record)
	result = db.Unscoped().Delete(&record)
	return result.Error
}

func exportTreeToExcel(values []string) {

}
