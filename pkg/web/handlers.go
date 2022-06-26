package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pavel1337/secretbox/pkg/forms"
	"github.com/pavel1337/secretbox/pkg/storage"
)

func (app *Web) showSecretForm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	s, err := app.storage.Get(id)
	if err == storage.ErrNoRecord {
		app.render(w, r, "secret404.page.tmpl", &templateData{})
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "showForm.page.tmpl", &templateData{
		Secret: s, Form: forms.New(nil),
	})
}

func (app *Web) showSecret(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get(":id")

	s, err := app.storage.GetAndDelete(id)

	if err == storage.ErrNoRecord {
		app.render(w, r, "secret404.page.tmpl", &templateData{})
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.MaxLength("passphrase", 1024)
	form.CheckPassword("passphrase", s.Passphrase)

	if !form.Valid() {
		form.Del("passphrase")
		app.render(w, r, "showForm.page.tmpl", &templateData{Secret: s, Form: form})
		return
	}

	plaintext, err := app.crypter.Decrypt(s.Content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	s.Content = plaintext

	app.render(w, r, "show.page.tmpl", &templateData{
		Secret: s,
	})
}

func (app *Web) createSecretForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *Web) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (app *Web) createSecret(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("content", "expires")
	form.MaxLength("content", 1024)
	form.MaxLength("passphrase", 1024)
	form.PermittedValues("expires", "10", "60", "1440")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	stringContent := form.Get("content")

	ciphertext, err := app.crypter.Encrypt([]byte(stringContent))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var s storage.Secret
	s.Content = ciphertext

	if strings.TrimSpace(form.Get("passphrase")) != "" {
		s.Passphrase = form.Get("passphrase")
	} else {
		s.Passphrase = ""
	}

	ttlMinutes, _ := strconv.Atoi(form.Get("expires"))

	id, err := app.storage.Insert(s, ttlMinutes)
	if err != nil {
		app.serverError(w, err)
		return
	}

	link := fmt.Sprintf("%s/secret/%s", app.url, id)
	err = app.addFlash(w, r, link)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "faq.page.tmpl", &templateData{Form: form})
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}
