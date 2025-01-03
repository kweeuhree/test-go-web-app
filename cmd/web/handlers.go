package main

import (
	"html/template"
	"net/http"
)

func (app *applicaton) Home(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *applicaton) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	// parse the template from disk
	parsedTemplate, err := template.ParseFiles("./templates/" + t)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	// execute the template passing it data if any
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
