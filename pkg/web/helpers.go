package web

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *Web) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	return td
}

func (app *Web) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))

	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

func (app *Web) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Web) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Web) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
