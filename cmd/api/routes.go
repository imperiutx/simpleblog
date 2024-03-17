package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// static files
	mux.HandleFunc("GET /favicon.ico", app.serveFavicon)
	mux.HandleFunc("GET /static/", app.serveStaticFiles)
	fs := http.FileServer(http.Dir("./templates"))
	mux.Handle("GET /v1/static/", http.StripPrefix("/static/", fs))

	// health check
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	//templates
	mux.HandleFunc("GET /", app.showMainPageHandler) //TODO: implement pagination
	mux.HandleFunc("GET /v1/posts/new", app.getCreatePostFormHandler)
	mux.HandleFunc("GET /v1/comments/new", app.getCreateCommentFormHandler)
	mux.HandleFunc("GET /v1/posts/{id}/edit", app.getPostForEditHandler)

	//users APIs
	mux.HandleFunc("POST /v1/users", app.createUserHandler)
	mux.HandleFunc("GET /v1/users/{id}", app.showUserHandler)
	mux.HandleFunc("PATCH /v1/users/{id}", app.updateUserHandler)

	//posts APIs
	mux.HandleFunc("POST /v1/posts", app.createPostHandler)
	mux.HandleFunc("GET /v1/posts/{id}", app.showPostHandler)
	mux.HandleFunc("PATCH /v1/posts/{id}", app.updatePostHandler)
	//TODO: listPostsByTagsHandler

	//comments API
	mux.HandleFunc("POST /v1/comments", app.createCommentHandler)
	//TODO: updateCommentHandler
	//TODO: deleteCommentHandler

	return app.recoverPanic(app.rateLimit(mux))
}
