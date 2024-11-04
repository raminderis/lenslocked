package controllers

import (
	"html/template"
	"net/http"
	"pro/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQHandler(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "is there a free version?",
			Answer:   "yes free for 30 days.",
		},
		{
			Question: "is there a abc version?",
			Answer:   "yes abc for 30 days.",
		},
		{
			Question: "is there a gfd version?",
			Answer:   "yes gfd for 30 days.",
		},
		{
			Question: "what your email?",
			Answer:   "my email is <a href='mailto:raminderis@live.co'>chmicanga@live.com</a>",
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
