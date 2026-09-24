package handler

import (
	"html/template"
	"net/http"
)

func MainPage(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(
			w,
			"Failed to load page",
			http.StatusInternalServerError,
		)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(
			w,
			"Failed to render page",
			http.StatusInternalServerError,
		)
	}
}