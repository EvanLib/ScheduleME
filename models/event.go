package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	Title       string
	Description string
}
type EventService interface {
	ByID(id uint) *Event
	GetAll() []Event
	Create(event *Event) error
	Update(event *Event) error
	Delete(event *Event) error
}

type EventGorm struct {
	*gorm.DB //Create once and interact with it
}

func (eg *EventGorm) DestructiveReset() {
	eg.DropTable(&Event{})
	eg.AutoMigrate(&Event{})
}

func NewEventGorm(db *gorm.DB) *EventGorm {
	return &EventGorm{db}
}

//CRUD Read Functions
func (eg *EventGorm) ByID(id uint) *Event {
	return eg.byQuery(eg.DB.Where("id = ?", id))
}

func (eg *EventGorm) GetAll() []Event {
	allEvents := []Event{}
	eg.DB.Find(&allEvents)
	return allEvents
}

func (ug *EventGorm) byQuery(query *gorm.DB) *Event {
	ret := &Event{}
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
func (eg *EventGorm) Create(event *Event) error {
	return eg.DB.Create(event).Error
}

func (eg *EventGorm) Update(event *Event) error {
	return eg.DB.Update(event).Error
}

func (eg *EventGorm) Delete(event *Event) error {
	return eg.DB.Delete(event).Error
}
