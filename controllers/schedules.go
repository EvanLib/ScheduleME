package controllers

import (
	"net/http"

	"github.com/EvanLib/me_users/models"
)

type Schedules struct {
	SchedulesService models.ScheduleService
	EventsService    models.EventService
}

type ScheduleForm struct {
	Title       string `scheme:"title"`
	Description string `scheme:"description"`
}

func NewSchedules(mi ModelsInteractor) *Schedules {
	return &Schedules{
		SchedulesService: mi.Schedules,
		EventsService:    mi.Events,
	}
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
