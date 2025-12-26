package controllers

import (
	"net/http"

	"github.com/raminderis/lenslocked/views"
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

func StaticHandler(tpl views.Template, data User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, data)
	}
}
