package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	//templates
	// mux.HandleFunc("GET /", app.showMainPageHandler)

	// health check
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	//users APIs
	mux.HandleFunc("POST /v1/users", app.createUserHandler)
	// mux.HandleFunc("GET /v1/users/{id}", app.showUserHandler)


	//movies APIs
	// mux.HandleFunc("POST /v1/movies", app.createMovieHandler)
	// mux.HandleFunc("GET /v1/movies/{id}", app.showMovieHandler)
	// mux.HandleFunc("GET /v1/movies/", app.listMoviesHandler)
	// mux.HandleFunc("PATCH /v1/movies/{id}", app.updateMovieHandler)
	// for API endpoints which perform partial updates on a resource,
	// it’s appropriate to the use the HTTP method PATCH
	// rather than PUT (which is intended for replacing a resource in full).

	return app.recoverPanic(
		app.rateLimit(
			mux),
	)
}
