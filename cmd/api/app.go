package main

import (
	"html/template"
	"net/http"
)

var data = []string{
	"HTMX is a great alternative to React etc.",
	"It offers a different way of loading data into your frontend web UI.",
	"It might be especially interesting for server-side developers who are not so familiar with frontend development.",
	"But - as you will see - it's actually also a very promising alternative to React, Angular etc.",
	"You just have to be open for a diffent mental model.",
	"When using HTMX you typically write way less frontend JavaScript code.",
	"You also don't need to manage any frontend state.",
	"Though you can always add extra JS code if needed.",
	"And you can also combine HTMX with other libraries like AlpineJS or integrate it into React apps etc.",
}

func (app *application) showAppPageHandler(w http.ResponseWriter, r *http.Request) {
	file1 := fePath + "main.tmpl"
	file2 := fePath + "info.tmpl"
	tmpl := template.Must(template.ParseFiles(file1, file2))

	if err := tmpl.ExecuteTemplate(w, "main", envelope{"Data": data}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) showInfoPageHandler(w http.ResponseWriter, r *http.Request) {

	file := fePath + "info.tmpl"
	tmpl := template.Must(template.ParseFiles(file))

	if err := tmpl.Execute(w, envelope{"Data": data}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) postNoteHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fv := r.PostFormValue("note")

	data = append([]string{fv}, data...)

	file := fePath + "info.tmpl"
	tmpl := template.Must(template.ParseFiles(file))

	if err := tmpl.ExecuteTemplate(w, "info", envelope{"Data": data}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
