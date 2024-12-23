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

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	//err := r.ParseForm()
	//if err != nil {
	//	fmt.Fprint(w, "Error", http.StatusBadRequest)
	//	return
	//}
	// email := r.PostForm.Get("email")
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Fprint(w, "Creating: ", email, " ", password)
}
