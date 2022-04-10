package web

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Web) routes() http.Handler {
	standartMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.createSecretForm))
	mux.Get("/secret/create", dynamicMiddleware.ThenFunc(app.createSecretForm))
	mux.Post("/secret/create", dynamicMiddleware.ThenFunc(app.createSecret))
	mux.Get("/secret/:id", dynamicMiddleware.ThenFunc(app.showSecretForm))
	mux.Post("/secret/:id", dynamicMiddleware.ThenFunc(app.showSecret))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standartMiddleware.Then(mux)
}
