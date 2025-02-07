package main

import (
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler{
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		app.notFount(w)
	})
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		app.notFount(w)
	})

	filesServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", filesServer))
	
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create",  dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create",  dynamic.ThenFunc(app.snippetCreatePost))

	standar := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standar.Then(router)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()

	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
