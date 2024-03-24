package main

import (
	"expvar"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// static files
	mux.HandleFunc("GET /favicon.ico", app.serveFavicon)
	mux.HandleFunc("GET /static/", app.serveStaticFiles)

	mux.Handle("GET /debug/vars", expvar.Handler())

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/home", app.showMainPageHandler) //TODO: implement pagination
	mux.HandleFunc("/login", app.loginUserHandler)
	mux.HandleFunc("/logout", logoutHandler)

	mux.HandleFunc("GET /admin/dashboard", app.showAdminDashboardHandler)

	// v1 := http.NewServeMux()
	// v1.Handle("/v1/", http.StripPrefix("/v1", mux))
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
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

	stack := createStack(
		app.logging,
		app.rateLimit,
		app.recoverPanic,
		app.metrics,
	)

	return stack(mux)
}
