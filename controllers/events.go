package controllers

import (
	"fmt"
	"net/http"

	"github.com/EvanLib/me_users/models"
	"github.com/EvanLib/me_users/views"
)

type Events struct {
	NewView      *views.View
	EventsView   *views.View
	UserService  models.UserService
	EventService models.EventService
}

type EventForm struct {
	Title       string `scheme:"title"`
	Description string `scheme:"description"`
}

func NewEvents(mi ModelsInteractor) *Events {
	return &Events{
		NewView:      views.NewView("bootstrap", "views/events/new.html"),
		EventsView:   views.NewView("bootstrap", "views/events/index.html"),
		UserService:  mi.Users,
		EventService: mi.Events,
	}
}

func (e *Events) New(w http.ResponseWriter, r *http.Request) {
	e.NewView.Render(w, nil)
}

func (e Events) Index(w http.ResponseWriter, r *http.Request) {

	events := e.EventService.GetAll()
	e.EventsView.Render(w, events)
}

//API Post
func (e *Events) Create(w http.ResponseWriter, r *http.Request) {
	form := EventForm{}
	if err := parseForm(r, &form); err != nil {
		fmt.Println(err)
		panic(err)
	}

	event := &models.Event{
		Title:       form.Title,
		Description: form.Description,
	}

	e.EventService.Create(event)
	http.Redirect(w, r, "/events", http.StatusFound)
}
