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

	// health check
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	//templates
	mux.HandleFunc("GET /", app.showMainPageHandler)            //TODO: implement pagination
	mux.HandleFunc("GET /admin", app.showAdminDashboardHandler)
	mux.HandleFunc("GET /v1/comments/new", app.getCreateCommentFormHandler)
	mux.HandleFunc("GET /v1/posts/{id}/edit", app.getPostForEditHandler)
	mux.HandleFunc("GET /v1/login", app.getUserLoginHandler)

	//users APIs
	mux.HandleFunc("POST /v1/users", app.createUserHandler)
	mux.HandleFunc("GET /v1/users/{id}", app.showUserHandler)
	mux.HandleFunc("PATCH /v1/users/{id}", app.updateUserHandler)
	mux.HandleFunc("POST /v1/users/login", app.loginUserHandler)

	//posts APIs
	mux.HandleFunc("POST /v1/posts", app.createPostHandler)
	mux.HandleFunc("GET /v1/posts/{id}", app.showPostHandler)
	mux.HandleFunc("PATCH /v1/posts/{id}", app.updatePostHandler)
	//TODO: listPostsByTagsHandler

	//comments API
	mux.HandleFunc("POST /v1/comments", app.createCommentHandler)
	//TODO: updateCommentHandler
	//TODO: deleteCommentHandler

	mux.Handle("GET /debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.rateLimit(mux)))
}
