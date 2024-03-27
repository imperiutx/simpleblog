package main

import (
	"html/template"
	"net/http"
	db "simpleblog/db/sqlc"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

type createCommentRequest struct {
	Username string `json:"username" validate:"required"`
	Content  string `json:"content" validate:"required"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	pid, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	switch r.Method {

	case http.MethodGet:
		if err := r.ParseForm(); err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		file := fePath + "create_comment.html"
		tmpl := template.Must(template.ParseFiles(file))

		if err := tmpl.Execute(w, envelope{"Post": pid}); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

	case http.MethodPost:
		switch r.Header.Get("Content-Type") {

		case "application/json":
			var req createCommentRequest

			if err := app.readJSON(w, r, &req); err != nil {
				app.badRequestResponse(w, r, err)
				return
			}

			validate := validator.New()

			if err := validate.Struct(req); err != nil {
				errors := err.(validator.ValidationErrors)
				app.badRequestResponse(w, r, errors)
				return
			}

			cmt := db.CreateCommentParams{
				PostID: pgtype.Int8{
					Int64: pid,
					Valid: true,
				},
				Username: pgtype.Text{
					String: req.Username,
					Valid:  true,
				},
				Content: req.Content,
			}

			newCmt, err := app.store.CreateComment(r.Context(), cmt)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			if err = app.writeJSON(
				w,
				http.StatusAccepted,
				envelope{"post": newCmt}, nil); err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

		case "application/x-www-form-urlencoded":
			username := "Guest"
			cookie, err := r.Cookie("username")
			if err == nil {
				username = cookie.Value
			}

			cmt := db.CreateCommentParams{
				PostID: pgtype.Int8{
					Int64: pid,
					Valid: true,
				},
				Username: pgtype.Text{
					String: username,
					Valid:  true,
				},
				Content: r.PostFormValue("content"),
			}

			_, err = app.store.CreateComment(r.Context(), cmt)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

		default:
			http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
			return
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		app.logger.Info("Using method", r.Method, "<<<<<<<")
		return
	}

}
