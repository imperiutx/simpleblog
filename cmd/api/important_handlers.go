package main

import (
	"html/template"
	"net/http"
	db "simpleblog/db/sqlc"

	"simpleblog/util"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.Environment,
			"version":     app.config.Version,
		},
	}

	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showAdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/admin_dashboard.html"))

	comments, err := app.store.ListAllComments(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	users, err := app.store.ListAllUsers(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := tmpl.Execute(w,
		envelope{
			"Users":    users,
			"Comments": comments}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) showMainPageHandler(w http.ResponseWriter, r *http.Request) {

	username := "Guest"
	cookie, err := r.Cookie("username")
	if err == nil {
		username = cookie.Value
	}

	data := struct {
		Username string
	}{
		Username: username,
	}
	
	var input struct {
		Title string
		Tags  []string
		db.Filters
	}

	v := util.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Tags = app.readCSV(qs, "tags", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 10)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafeValues = []string{
		"id", "title", "-id", "-title",
	}

	if db.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	posts, metadata, err := app.store.ListPostsDynamic(
		r.Context(),
		input.Title,
		input.Tags,
		input.Filters,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	author, err := app.store.GetUserById(r.Context(), 1)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/index.html"))

	if err := tmpl.Execute(w,
		envelope{
			"Posts":    posts,
			"Metadata": metadata,
			"Author":   author,
			"Data":     data}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
