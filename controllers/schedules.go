package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/EvanLib/me_users/models"
	"github.com/EvanLib/me_users/views"
	"github.com/gorilla/mux"
)

type Schedules struct {
	SchedulesService models.ScheduleService
	EventsService    models.EventService
	NewView          *views.View
	SingleView       *views.View
}

type ScheduleForm struct {
	Title       string `scheme:"title"`
	Description string `scheme:"description"`
}

func NewSchedules(mi ModelsInteractor) *Schedules {
	return &Schedules{
		SchedulesService: mi.Schedules,
		NewView:          views.NewView("bootstrap", "views/schedules/new.html"),
		SingleView:       views.NewView("bootstrap", "views/schedules/view.html"),
		EventsService:    mi.Events,
	}
}
func (s *Schedules) SingleSchedule(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		println("Not an id..")
		return
	}
	schedule := s.SchedulesService.ByID(uint(id))
	s.SingleView.Render(w, schedule.Events)
}

func (s *Schedules) New(w http.ResponseWriter, r *http.Request) {
	s.NewView.Render(w, nil)
}

func (s *Schedules) Create(w http.ResponseWriter, r *http.Request) {
	form := &ScheduleForm{}
	if err := parseForm(r, form); err != nil {
		panic(err)
	}

	schedule := &models.Schedule{
		Title:       form.Title,
		Description: form.Description,
	}
	s.SchedulesService.Create(schedule)
}

func (s *Schedules) AddEvent(w http.ResponseWriter, r *http.Request) {
	var schid = r.FormValue("scheduleid")
	var evnid = r.FormValue("eventid")

	scheduleID, err := strconv.ParseUint(schid, 10, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	eventID, err := strconv.ParseUint(evnid, 10, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	event := s.EventsService.ByID(uint(eventID))
	if event == nil {
		fmt.Println("Event not found")
	}

	schedule := s.SchedulesService.ByID(uint(scheduleID))
	if schedule == nil {
		fmt.Println("Schedule not found") //Need to implement logging
		return
	}
	s.SchedulesService.AddEvent(schedule, event)

	//Testin relations
	//sc := s.SchedulesService.ByID(1)
	//event := s.EventsService.ByID(2)
	//s.SchedulesService.AddEvent(sc, event)
	//s.SchedulesService.Update(sc)
}
