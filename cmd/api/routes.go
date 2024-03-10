package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	//templates
	mux.HandleFunc("GET /", app.showMainPageHandler)

	// health check
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	//users APIs
	mux.HandleFunc("POST /v1/users", app.createUserHandler)
	mux.HandleFunc("GET /v1/users/{id}", app.showUserHandler)
	mux.HandleFunc("PATCH /v1/users/{id}", app.updateUserHandler)

	//posts APIs
	mux.HandleFunc("POST /v1/posts", app.createPostHandler)
	mux.HandleFunc("GET /v1/posts/{id}", app.showPostHandler)
	// mux.HandleFunc("GET /v1/posts/", app.listPostsHandler)
	mux.HandleFunc("PATCH /v1/posts/{id}", app.updatePostHandler)
	// for API endpoints which perform partial updates on a resource,
	// it’s appropriate to the use the HTTP method PATCH
	// rather than PUT (which is intended for replacing a resource in full).

	//comments API

	return app.recoverPanic(
		app.rateLimit(
			mux),
	)
}
