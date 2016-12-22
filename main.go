package main

import (
	"log"
	"net/http"
	"time"

	"github.com/EvanLib/me_users/controllers"
	"github.com/EvanLib/me_users/models"
	"github.com/EvanLib/me_users/views"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var indexView *views.View

func main() {
	db, err := gorm.Open("mysql", "root:NOTMYPASSWORD@/me_schedule?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	ug := models.NewUserGorm(db)
	eg := models.NewEventGorm(db)
	sg := models.NewSchedulsGorm(db)

	modelsInteractor := controllers.ModelsInteractor{
		Users:     ug,
		Events:    eg,
		Schedules: sg,
	}

	if err != nil {
		panic(err)
	}

	indexView = views.NewView("bootstrap", "views/index.html")
	//Create controllers
	ug.DestructiveReset()
	usersC := controllers.NewUsers(modelsInteractor)
	eventsC := controllers.NewEvents(modelsInteractor)
	scheduleC := controllers.NewSchedules(modelsInteractor)

	//Create a mux
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	// User Handles
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.LoginGet).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	//Events Handle
	r.HandleFunc("/events/new", eventsC.New).Methods("GET")
	r.HandleFunc("/events/new", eventsC.Create).Methods("POST")
	r.HandleFunc("/events", eventsC.Index).Methods("GET")
	//Schedules Handles
	r.HandleFunc("/schedules/new", scheduleC.Create).Methods("POST")

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static_files/"))))

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	indexView.Render(w, nil)
}
