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

	//datastar
	mux.HandleFunc("/datastar", app.showDataStarPostHandler) //TODO: implement pagination

	//posts APIs
	mux.HandleFunc("GET /v1/posts", app.listDataStarPostHandler)
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
	mux.HandleFunc("GET /v1/contacts/new", app.showCreateContactHandler)
	mux.HandleFunc("POST /v1/contacts", app.createContactHandler)
	mux.HandleFunc("GET /v1/contacts", app.listContactsHandler)
	mux.HandleFunc("DELETE /v1/contacts/{id}", app.deleteContactHandler)
	mux.HandleFunc("/v1/contacts/{id}/edit", app.updateContactHandler)

	// udm apps
	mux.HandleFunc("GET /udm/v1", app.showAppPageHandler)
	mux.HandleFunc("GET /udm/v1/info", app.showInfoPageHandler)
	mux.HandleFunc("POST /udm/v1/note", app.postNoteHandler)
	mux.HandleFunc("GET /udm/v2/goals", app.showGoalPageHandler)
	mux.HandleFunc("POST /udm/v2/goals", app.postGoalHandler)
	mux.HandleFunc("DELETE /udm/v2/goals/{id}", app.deleteGoalHandler)
	mux.HandleFunc("GET /udm/v3/places", app.showPlacePageHandler)
	mux.HandleFunc("POST /udm/v3/places", app.postPlaceHandler)
	mux.HandleFunc("DELETE /udm/v3/places/{id}", app.deletePlaceHandler)

	stack := createStack(
		app.logging,
		app.rateLimit,
		app.enableCORS,
		app.recoverPanic,
		app.metrics,
	)

	return stack(mux)
}
