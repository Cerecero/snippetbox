package main

import (
	"net/http"
	"path/filepath"

	"github.com/Cerecero/snippetbox/config"
)

func routes(app *config.Application) *http.ServeMux{
	mux := http.NewServeMux()
	filesServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static/", http.StripPrefix("/static", filesServer))

	mux.HandleFunc("/", home(app))
	mux.HandleFunc("/snippet/view", snippetView(app))
	mux.HandleFunc("/snippet/create",  snippetCreate(app))

	return mux
}
// mux.Handle("/static", http.NotFoundHandler())

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
