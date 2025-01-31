package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Cerecero/snippetbox/config"
)

func serverError(app *config.Application,w http.ResponseWriter, err error) {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.ErrorLog.Output(2, trace)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func clientError(_ *config.Application,w http.ResponseWriter, status int) {
		http.Error(w, http.StatusText(status), status)
}

func notFount(app *config.Application, w http.ResponseWriter) {
	clientError(app, w, http.StatusNotFound)
}

func render(app *config.Application, w http.ResponseWriter, status int, page string, data *templateData){
	ts, ok := app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the tempolate %s does not exist", page)
		serverError(app, w, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		serverError(app, w, err)
	}
}
