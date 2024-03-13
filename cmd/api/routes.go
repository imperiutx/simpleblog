package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	//templates
	mux.HandleFunc("GET /", app.showMainPageHandler)
	mux.HandleFunc("GET /v1/post/new", app.getCreateFormHandler)
	mux.HandleFunc("GET /v1/post/{id}/edit", app.getPostForEditHandler)

	// health check
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	//users APIs
	mux.HandleFunc("POST /v1/users", app.createUserHandler)
	mux.HandleFunc("GET /v1/users/{id}", app.showUserHandler)
	mux.HandleFunc("PATCH /v1/users/{id}", app.updateUserHandler)

	//posts APIs
	mux.HandleFunc("POST /v1/post", app.createPostHandler)
	mux.HandleFunc("GET /v1/post/{id}", app.showPostHandler)
	// mux.HandleFunc("GET /v1/posts/", app.listPostsHandler)
	mux.HandleFunc("PATCH /v1/post/{id}", app.updatePostHandler)

	//comments API

	return app.recoverPanic(
		app.rateLimit(
			mux),
	)
}
