package controllers

import (
	"github.com/EvanLib/me_users/models"
)

type ModelsInteractor struct {
	Users     models.UserService
	Events    models.EventService
	Schedules models.ScheduleService
}
