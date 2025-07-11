package renderer

import (
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

const templateBasePath = "templates"

type TemplateCache map[string]*template.Template

func NewTemplateCache() (TemplateCache, error) {
	cache := TemplateCache{}

	layouts, err := filepath.Glob(templateBasePath + "/layouts/*.html")
	if err != nil {
		return nil, err
	}

	pages, err := filepath.Glob(templateBasePath + "/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		files := append(layouts, page)
		name := filepath.Base(page)
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = tmpl
	}

	return cache, nil
}

type TemplateRenderer struct {
	templates TemplateCache
}

func NewTemplateRenderer() *TemplateRenderer {
	templates, err := NewTemplateCache()
	if err != nil {
		log.Fatal("failed to create TemplateRenderer: ", err)
	}
	return &TemplateRenderer{templates: templates}
}

func (r *TemplateRenderer) Render(w http.ResponseWriter, name string, data any) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return errors.New("template not found: " + name)
	}
	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		return err
	}
	return nil
}
