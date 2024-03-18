package main

import (
	"html/template"
	"net/http"
	db "simpleblog/db/sqlc"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

type createCommentRequest struct {
	PostID   int64  `json:"post_id" validate:"required"`
	Username string `json:"username" validate:"required"`
	Content  string `json:"content" validate:"required"`
}

func (app *application) getCreateCommentFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/create_comment.html"))
	if err := r.ParseForm(); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
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
				Int64: req.PostID,
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
		pid, err := strconv.ParseInt(r.PostFormValue("postid"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		cmt := db.CreateCommentParams{
			PostID: pgtype.Int8{
				Int64: pid,
				Valid: true,
			},
			Username: pgtype.Text{
				String: r.PostFormValue("username"),
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

}
