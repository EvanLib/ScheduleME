package main

import (
	"log"
	"net/http"
	"time"

	"github.com/EvanLib/me_users/controllers"
	"github.com/EvanLib/me_users/models"
	"github.com/EvanLib/me_users/views"
	"github.com/gorilla/mux"
)

var indexView *views.View

func main() {

	ug, err := models.NewUserGorm("root:lol626465@/me_schedule?charset=utf8&parseTime=True&loc=Local")
	eg, err := models.NewEventGorm("root:lol626465@/me_schedule?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	indexView = views.NewView("bootstrap", "views/index.html")
	//Create controllers
	ug.DestructiveReset()
	usersC := controllers.NewUsers(ug)
	eventsC := controllers.NewEvents(eg)

	//Create a mux
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.LoginGet).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	r.HandleFunc("/events/new", eventsC.New).Methods("GET")
	r.HandleFunc("/events/new", eventsC.Create).Methods("POST")
	r.HandleFunc("/events", eventsC.Index).Methods("GET")
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
