package web

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Web) routes() http.Handler {
	standartMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()
	mux.Get("/ping", http.HandlerFunc(app.ping))
	mux.Get("/", http.HandlerFunc((app.createSecretForm)))
	mux.Get("/secret/create", http.HandlerFunc((app.createSecretForm)))
	mux.Post("/secret/create", http.HandlerFunc((app.createSecret)))
	mux.Post("/secret/generate", http.HandlerFunc((app.generatePassword)))
	mux.Get("/secret/:id", http.HandlerFunc((app.showSecretForm)))
	mux.Post("/secret/:id", http.HandlerFunc((app.showSecret)))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standartMiddleware.Then(mux)
}
