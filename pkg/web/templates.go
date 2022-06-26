package web

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/pavel1337/secretbox/pkg/forms"
	"github.com/pavel1337/secretbox/pkg/storage"
)

type templateData struct {
	CurrentYear int
	Flash       string
	Form        *forms.Form
	Secret      *storage.Secret
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func stringify(bb []byte) string {
	return string(bb)
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"stringify": stringify,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
