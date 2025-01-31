package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Cerecero/snippetbox/config"
	"github.com/Cerecero/snippetbox/internal/models"
)

func home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFount(app, w)
			return
		}

		snippets, err := app.Snippets.Lastest()
		if err != nil {
			serverError(app, w, err)
			return
		}
		render(app, w, http.StatusOK, "home.tmpl", &templateData{
			Snippets: snippets,
		})
		// files := []string{
		// 	"./ui/html/base.tmpl",
		// 	"./ui/html/pages/home.tmpl",
		// 	"./ui/html/partials/nav.tmpl",
		// }
		// ts, err := template.ParseFiles(files...)
		// if err != nil {
		// 	serverError(app, w, err)
		// 	return
		// }
		//
		// data := &templateData{
		// 	Snippets: snippets,
		// }
		// err = ts.ExecuteTemplate(w, "base", data)
		// if err != nil {
		// 	serverError(app, w, err)
		// }
	}
}
func snippetView(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			notFount(app, w)
			return
		}
		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				notFount(app, w)
			} else {
				serverError(app, w, err)
			}
			return
		}

		render(app, w, http.StatusOK, "view.tmpl", &templateData{
			Snippet: snippet,
		})
		// files := []string{
		// 	"./ui/html/base.tmpl",
		// 	"./ui/html/partials/nav.tmpl",
		// 	"./ui/html/pages/view.tmpl",
		// }
		//
		// ts, err := template.ParseFiles(files...)
		// if err != nil {
		// 	serverError(app, w, err)
		// 	return
		// }
		//
		// data := templateData{
		// 	Snippet: snippet,
		// }
		// err = ts.ExecuteTemplate(w, "base", data)
		// if err != nil {
		// 	serverError(app, w, err)
		// }
	}
}
func snippetCreate(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			clientError(app, w, http.StatusMethodNotAllowed)
			return
		}
		title := "0 snail"
		content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
		expires := 7

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			serverError(app, w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	}
}
