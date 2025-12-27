package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
}

type Account struct {
	Email    string
	Password string
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	recievedAcct := Account{
		Email: r.FormValue("email"),
	}
	u.Templates.New.Execute(w, recievedAcct)
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
	fmt.Fprintf(w, "Temprory Response email is %v and password is %v", newAcct.Email, newAcct.Password)
}
