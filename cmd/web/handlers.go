package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pavel1337/secretbox/pkg/forms"
	"github.com/pavel1337/secretbox/pkg/models"
)

func (app *application) showSecretForm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	s, err := app.secrets.Get(id, app.encryptionKey)

	if err == models.ErrNoRecord {
		app.render(w, r, "secret404.page.tmpl", &templateData{})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "showForm.page.tmpl", &templateData{
		Secret: s, Form: forms.New(nil),
	})
}

func (app *application) showSecret(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get(":id")

	s, err := app.secrets.Get(id, app.encryptionKey)

	if err == models.ErrNoRecord {
		app.render(w, r, "secret404.page.tmpl", &templateData{})
		return
	} else if err != nil {
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

	err = app.secrets.Delete(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Secret: s,
	})
}

func (app *application) createSecretForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSecret(w http.ResponseWriter, r *http.Request) {
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

	var s models.Secret

	s.Content = form.Get("content")

	if strings.TrimSpace(form.Get("passphrase")) != "" {
		s.Passphrase = form.Get("passphrase")
	} else {
		s.Passphrase = ""
	}

	i, _ := strconv.Atoi(form.Get("expires"))
	key, err := app.secrets.Insert(s, app.encryptionKey, i)
	if err != nil {
		app.serverError(w, err)
		return
	}
	link := fmt.Sprintf("%s/secret/%s", app.config.Url, key)
	app.session.Put(r, "flash", link)

	app.render(w, r, "faq.page.tmpl", &templateData{Form: form})
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}
