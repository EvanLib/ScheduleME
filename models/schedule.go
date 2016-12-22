package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Schedule struct {
	gorm.Model
	Title       string
	Description string
}
type ScheduleService interface {
	ByID(id uint) *Schedule
	GetAll() []Schedule
	Create(Schedule *Schedule) error
	Update(Schedule *Schedule) error
	Delete(Schedule *Schedule) error
}

type Schedulsgorm struct {
	*gorm.DB //Create once and interact with it
}

func (sg *Schedulsgorm) DestructiveReset() {
	sg.DropTable(&Schedule{})
	sg.AutoMigrate(&Schedule{})
}

func NewSchedulsGorm(db *gorm.DB) *Schedulsgorm {
	return &Schedulsgorm{db}
}

//CRUD Read Functions
func (sg *Schedulsgorm) ByID(id uint) *Schedule {
	return sg.byQuery(sg.DB.Where("id = ?", id))
}

func (sg *Schedulsgorm) GetAll() []Schedule {
	allSchedules := []Schedule{}
	sg.DB.Find(&allSchedules)
	return allSchedules
}

func (ug *Schedulsgorm) byQuery(query *gorm.DB) *Schedule {
	ret := &Schedule{}
	err := query.First(ret).Error
	switch err {
	case nil:
		return ret
	case gorm.ErrRecordNotFound:
		return nil
	default:
		panic(err)
	}
}

//CRUD Functions
func (sg *Schedulsgorm) Create(Schedule *Schedule) error {
	return sg.DB.Create(Schedule).Error
}

func (sg *Schedulsgorm) Update(Schedule *Schedule) error {
	return sg.DB.Update(Schedule).Error
}

func (sg *Schedulsgorm) Delete(Schedule *Schedule) error {
	return sg.DB.Delete(Schedule).Error
}
