package controllers

import (
	"fmt"
	"net/http"

	"github.com/EvanLib/me_users/models"
	"github.com/EvanLib/me_users/views"
	"github.com/EvanLib/rand"
	_ "github.com/go-sql-driver/mysql"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	models.UserService
}
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:     views.NewView("bootstrap", "views/users/new.html"),
		LoginView:   views.NewView("bootstrap", "views/users/login.html"),
		UserService: us,
	}
}

//Get /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

//Post /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	form := SignupForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := &models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.UserService.Create(user); err != nil {
		panic(err)
	}
	if err := u.signIn(w, user); err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		fmt.Fprintln(w, "Err retrivieing cookie: ", err)
	}
	fmt.Fprintln(w, "RememberToken is:", cookie.Value)
	user := u.UserService.ByRemember(cookie.Value)
	fmt.Fprintln(w, "User is : ", user)
}

//Authentication
func (u *Users) LoginGet(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Render(w, nil)
}
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := u.UserService.Authenticate(form.Email, form.Password)
	if user == nil {
		fmt.Fprintln(w, "Invalid Login")
	}

	if err := u.signIn(w, user); err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)

}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	rememberToken, err := rand.RemeberToken()
	if err != nil {
		return err
	}

	user.Remember = rememberToken
	if err := u.UserService.Update(user); err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  "remember_token",
		Value: rememberToken,
	}
	http.SetCookie(w, cookie)
	return nil
}