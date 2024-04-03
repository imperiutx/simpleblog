package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	db "simpleblog/db/sqlc"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

type createPostRequest struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content" validate:"required"`
	Tags    []string `json:"tags"`
}

type postResponse struct {
	Username string    `json:"username"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Tags     []string  `json:"tags"`
	EditedAt time.Time `json:"created_at"`
}

func newPostResponse(post db.Post) postResponse {
	return postResponse{
		Username: post.Username.String,
		Title:    post.Title,
		Content:  post.Content,
		Tags:     post.Tags,
		EditedAt: post.EditedAt,
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
			Tags:    []string{"news"},
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

	case "application/x-www-form-urlencoded":

		post := db.CreatePostParams{
			Username: pgtype.Text{
				String: "rootadmin",
				Valid:  true,
			},
			Title:   r.PostFormValue("title"),
			Content: r.PostFormValue("content"),
			Tags:    []string{"news"},
		}

		_, err := app.store.CreatePost(r.Context(), post)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

	default:
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}
}

func (app *application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	pid, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post, err := app.store.GetPostById(r.Context(), pid)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	postID := pgtype.Int8{
		Int64: pid,
		Valid: true,
	}
	comments, err := app.store.ListCommentsByPostID(r.Context(), postID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	file := fePath + "post.html"
	tmpl := template.Must(template.ParseFiles(file))

	if err := tmpl.Execute(w, envelope{"Post": post, "Comments": comments}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	pid, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	switch r.Method {

	case http.MethodGet:
		post, err := app.store.GetPostForUpdate(r.Context(), pid)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				app.notFoundResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		file := fePath + "edit_post.html"
		tmpl := template.Must(template.ParseFiles(file))

		if err := tmpl.Execute(w, post); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

	case http.MethodPatch:
		switch r.Header.Get("Content-Type") {
		case "application/json":

			var updatePostRequest struct {
				Title   string `json:"title"`
				Content string `json:"content"`
			}

			if err := app.readJSON(w, r, &updatePostRequest); err != nil {
				app.badRequestResponse(w, r, err)
				return
			}

			data := db.UpdatePostParams{
				ID: pid,
				Title: pgtype.Text{
					String: updatePostRequest.Title,
					Valid:  true,
				},
				Content: pgtype.Text{
					String: updatePostRequest.Content,
					Valid:  true,
				},
			}

			_, err = app.store.UpdatePost(r.Context(), data)
			if err != nil {
				if errors.Is(err, db.ErrRecordNotFound) {
					app.notFoundResponse(w, r)
					return
				}
				app.serverErrorResponse(w, r, err)
				return
			}

			if err = app.writeJSON(
				w,
				http.StatusAccepted,
				envelope{"update_post": "success"}, nil); err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

		case "application/x-www-form-urlencoded":
			if err := r.ParseForm(); err != nil {
				app.badRequestResponse(w, r, err)
				return
			}

			data := db.UpdatePostParams{
				ID: pid,
				Title: pgtype.Text{
					String: r.FormValue("title"),
					Valid:  true,
				},
				Content: pgtype.Text{
					String: r.FormValue("content"),
					Valid:  true,
				},
			}

			_, err := app.store.UpdatePost(r.Context(), data)
			if err != nil {
				if errors.Is(err, db.ErrRecordNotFound) {
					app.notFoundResponse(w, r)
					return
				}
				app.serverErrorResponse(w, r, err)
				return
			}

			u := fmt.Sprintf("/v1/posts/%d", pid)
			http.Redirect(w, r, u, http.StatusSeeOther)
			return

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

func (app *application) listDataStarPostHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.store.ListPosts(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for _, post := range posts {
		fmt.Fprintf(w, "data: %v\n\n", post)
		flusher.Flush()
	}
}

func (app *application) showDataStarPostHandler(w http.ResponseWriter, r *http.Request) {
	file := fePath + "home.html"
	tmpl := template.Must(template.ParseFiles(file))

	if err := tmpl.Execute(w, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
