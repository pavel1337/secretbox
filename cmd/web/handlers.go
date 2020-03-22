package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pavel1337/secretbox/pkg/forms"
	"github.com/pavel1337/secretbox/pkg/models"
)

func (app *application) showSecret(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	s, err := app.secrets.Get(id)

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

func (app *application) showSecretTrue(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	s, err := app.secrets.Get(id)

	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	content, err := DecryptSecret(s.EncContent, app.encryptionKey)
	if err != nil {
		app.serverError(w, err)
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

	s.Content = string(content)
	app.render(w, r, "show_true.page.tmpl", &templateData{
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
	form.PermittedValues("expires", "10", "60", "1440")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	encContent, err := EncryptSecret([]byte(form.Get("content")), app.encryptionKey)
	if err != nil {
		app.serverError(w, err)
		return
	}
	i, _ := strconv.Atoi(form.Get("expires"))
	key, err := app.secrets.Insert(encContent, i)
	if err != nil {
		app.serverError(w, err)
		return
	}
	link := fmt.Sprintf("%s/secret/%s", app.config.Url, key)
	app.session.Put(r, "flash", link)

	app.render(w, r, "faq.page.tmpl", &templateData{Form: form})
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}
