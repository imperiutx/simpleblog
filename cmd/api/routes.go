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
	mux.HandleFunc("/logout", app.logoutHandler)

	mux.HandleFunc("GET /admin/dashboard", app.showAdminDashboardHandler)
	mux.HandleFunc("GET /contacts/dashboard", app.showContactDashboardHandler)

	// v1 := http.NewServeMux()
	// v1.Handle("/v1/", http.StripPrefix("/v1", mux))
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	//users APIs
	mux.HandleFunc("POST /v1/users", app.createUserHandler)
	mux.HandleFunc("GET /v1/users/{id}", app.showUserHandler)
	mux.HandleFunc("PATCH /v1/users/{id}", app.updateUserHandler)

	//posts APIs
	mux.HandleFunc("POST /v1/posts", app.createPostHandler)
	mux.HandleFunc("GET /v1/posts/{id}", app.showPostHandler)
	mux.HandleFunc("/v1/posts/{id}/edit", app.updatePostHandler)
	//TODO: listPostsByTagsHandler

	//comments API
	mux.HandleFunc("/v1/comments/{id}/new", app.createCommentHandler)
	mux.HandleFunc("GET /v1/comments", app.ListCommentsHandler)
	//TODO: updateCommentHandler
	//TODO: deleteCommentHandler
	mux.HandleFunc("DELETE /v1/comments/{id}", app.deleteCommentHandler)
	
	//contacts API
	mux.HandleFunc("POST /v1/contacts", app.createContactHandler)
	mux.HandleFunc("GET /v1/contacts", app.listContactsHandler)
	mux.HandleFunc("DELETE /v1/contacts/{id}", app.deleteContactHandler)
	mux.HandleFunc("/v1/contacts/{id}/edit", app.updateContactHandler)

	stack := createStack(
		app.logging,
		app.rateLimit,
		app.recoverPanic,
		app.metrics,
	)

	return stack(mux)
}
