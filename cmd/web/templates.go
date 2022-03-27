package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/swissarmybox/snippetbox/pkg/forms"
	"github.com/swissarmybox/snippetbox/pkg/models"
	"github.com/swissarmybox/snippetbox/ui"
)

type templateData struct {
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	CSRFToken       string
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	User            *models.User
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCacheNonEmbed(dir string) (map[string]*template.Template, error) {
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

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/*.page.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFS(ui.Files, "html/*.layout.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFS(ui.Files, "html/*.partial.tmpl")
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
