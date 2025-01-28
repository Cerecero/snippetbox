package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Cerecero/snippetbox/config"
)

// type application struct {
// 	errorLog *log.Logger
// 	infoLog *log.Logger
// }

func main() {
	// Logger
	logFile, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	defer logFile.Close()
	infoLog := log.New(mw, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(mw, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Command Line flag addr, default address :4000
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	
	app := &config.Application{
		ErrorLog: errorLog,
		InfoLog: infoLog,
	}
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", home(app))
	mux.HandleFunc("/snippet/view", snippetView(app))
	mux.HandleFunc("/snippet/create", snippetCreate(app))
	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
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
