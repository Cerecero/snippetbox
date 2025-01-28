package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Cerecero/snippetbox/config"
)

func home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/pages/home.tmpl",
			"./ui/html/partials/nav.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
}
func snippetView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}
func snippetCreate(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte("Create a new snippet..."))
	}
}
