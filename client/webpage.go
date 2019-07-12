package client

import (
	"html/template"
	"net/http"
)

func DisplayWebPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "error loading template :" + err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}