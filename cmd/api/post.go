package main

import (
	"net/http"
	db "simpleblog/db/sqlc"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

type createPostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type postResponse struct {
	Username  string    `json:"username"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func newPostResponse(post db.Post) postResponse {
	return postResponse{
		Username:  post.Username.String,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	}
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var req createPostRequest

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

		post := db.CreatePostParams{
			Username: pgtype.Text{
				String: "rootadmin",
				Valid:  true,
			},
			Title:   req.Title,
			Content: req.Content,
		}

		newPost, err := app.store.CreatePost(r.Context(), post)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		rsp := newPostResponse(newPost)

		if err = app.writeJSON(
			w,
			http.StatusAccepted,
			envelope{"post": rsp}, nil); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}
