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
