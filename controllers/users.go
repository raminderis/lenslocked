package controllers

import (
	"fmt"
	"net/http"

	"github.com/raminderis/lenslocked/models"
)

type Users struct {
	Templates struct {
		New     Template
		Signin  Template
		General Template
	}
	Message        string
	UserService    *models.UserService
	SessionService *models.SessionService
}

type Account struct {
	Email    string
	Password string
	Message  string
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	recievedAcct := Account{
		Email: r.FormValue("email"),
	}
	u.Templates.New.Execute(w, r, recievedAcct)
}

func (u Users) Signin(w http.ResponseWriter, r *http.Request) {
	recievedAcct := Account{
		Email: r.FormValue("email"),
	}
	u.Templates.Signin.Execute(w, r, recievedAcct)
}

func (u Users) General(w http.ResponseWriter, r *http.Request) {
	recievedAcct := Account{
		Message: u.Message,
	}
	u.Templates.General.Execute(w, r, recievedAcct)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := r.PostForm
	newAcct := Account{
		Email:    s.Get("email"),
		Password: s.Get("password"),
	}
	// fmt.Fprintf(w, "Attempting to create user with email %v and password notprintedHere\n", newAcct.Email)
	user, err := u.UserService.Create(newAcct.Email, newAcct.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
	// fmt.Fprintf(w, "User created: %v", user.Email)
	// recievedAcct := Account{
	// 	Message: "User " + user.Email + " is created",
	// }
	// u.Templates.General.Execute(w, r, recievedAcct)
}

func (u Users) SigninProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := r.PostForm
	newAcct := Account{
		Email:    s.Get("email"),
		Password: s.Get("password"),
	}
	// fmt.Fprintf(w, "Attempting to Login User with email %v and password notprintedHere\n", newAcct.Email)
	user, err := u.UserService.Authenticate(newAcct.Email, newAcct.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	// fmt.Fprintf(w, "User authenticated: %+v", user.Email)
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)

	// recievedAcct := Account{
	// 	Message: "User " + user.Email + " is authenticated",
	// }
	// u.Templates.General.Execute(w, r, recievedAcct)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "signin", http.StatusFound)
		return
	}
	fmt.Println("somewhere here")
	recievedAcct := Account{
		Message: "Current User is: " + user.Email,
	}
	u.Templates.General.Execute(w, r, recievedAcct)
}
