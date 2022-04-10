package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/pavel1337/secretbox/pkg/crypt"
	"github.com/pavel1337/secretbox/pkg/storage"
)

type Web struct {
	url           string
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	storage       storage.Store
	crypter       crypt.Crypter
	templateCache map[string]*template.Template
}

func New(sessionSecretKey string, store storage.Store, crypter crypt.Crypter) *Web {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session := sessions.New([]byte(sessionSecretKey))
	session.Lifetime = 12 * time.Hour

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	return &Web{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		storage:       store,
		crypter:       crypter,
		templateCache: templateCache,
	}
}

func (w *Web) Start(addr string) error {
	srv := &http.Server{
		Addr:     addr,
		ErrorLog: w.errorLog,
		Handler:  w.routes(),
	}

	w.infoLog.Printf("Starting server on %s\n", addr)
	return srv.ListenAndServe()
}
