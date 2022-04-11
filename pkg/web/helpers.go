package web

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func (app *Web) addDefaultData(td *templateData, w http.ResponseWriter, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.Flash = app.flashOrEmpty(w, r)
	return td
}

func (app *Web) flashOrEmpty(w http.ResponseWriter, r *http.Request) string {
	// ignoring error as Get() always returns a session, even if empty
	session, _ := app.session.Get(r, "session")
	flashes := session.Flashes()
	if len(flashes) != 1 {
		return ""
	}
	err := session.Save(r, w)
	if err != nil {
		return ""
	}
	return flashes[0].(string)
}

func (app *Web) addFlash(w http.ResponseWriter, r *http.Request, link string) error {
	// ignoring error as Get() always returns a session, even if empty
	session, _ := app.session.Get(r, "session")
	session.AddFlash(link)
	return session.Save(r, w)
}

func (app *Web) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, w, r))

	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

func (app *Web) serverError(w http.ResponseWriter, err error) {
	// trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Web) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Web) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
