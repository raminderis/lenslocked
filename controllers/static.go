package controllers

import (
	"net/http"
)

type Questions struct {
	Question string
	Answer   string
}

type User struct {
	Email string
	Phone string
	QA    []Questions
}

func StaticHandler(tpl Template, data User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, data)
	}
}
