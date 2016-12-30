package main

import (
	"fmt"
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
	db, err := gorm.Open("mysql", "root:somepassword@/me_schedule?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	//LOGGIN
	db.LogMode(false)
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

	//Create controllers
	ug.DestructiveReset()
	eg.DestructiveReset()
	sg.DestructiveReset()

	usersC := controllers.NewUsers(modelsInteractor)
	eventsC := controllers.NewEvents(modelsInteractor)
	scheduleC := controllers.NewSchedules(modelsInteractor)

	//Create a mux
	r := mux.NewRouter()

	// User Handles
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.LoginGet).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	//Events Handle
	r.HandleFunc("/events/new", usersC.AuthenticateCookie(eventsC.New)).Methods("GET")
	r.HandleFunc("/events/new", usersC.AuthenticateCookie(eventsC.Create)).Methods("POST")
	r.HandleFunc("/events", usersC.AuthenticateCookie(eventsC.Index)).Methods("GET")

	//Schedules Handles
	r.HandleFunc("/schedules/new", usersC.AuthenticateCookie(scheduleC.Create)).Methods("POST")
	r.HandleFunc("/schedules/new", usersC.AuthenticateCookie(scheduleC.New)).Methods("GET")
	r.HandleFunc("/schedules/{id:[0-9]+}", usersC.AuthenticateCookie(scheduleC.SingleSchedule)).Methods("GET")

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
