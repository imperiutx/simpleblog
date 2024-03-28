package main

import (
	"errors"
	"html/template"
	"net/http"
	db "simpleblog/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func (app *application) showContactDashboardHandler(w http.ResponseWriter, r *http.Request) {
	file := ctPath + "index.html"
	tmpl := template.Must(template.ParseFiles(file))

	if err := tmpl.Execute(w, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	file1 := ctPath + "create_contact.html"
	tmpl1 := template.Must(template.ParseFiles(file1))

	if err := tmpl1.Execute(w, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createContactHandler(w http.ResponseWriter, r *http.Request) {
	conctact := db.CreateContactParams{
		FirstName: pgtype.Text{
			String: r.PostFormValue("first_name"),
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: r.PostFormValue("last_name"),
			Valid:  true,
		},
		Email: r.PostFormValue("email"),
		Phone: pgtype.Text{
			String: r.PostFormValue("phone"),
			Valid:  true,
		},
	}

	_, err := app.store.CreateContact(r.Context(), conctact)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/v1/contacts/dashboard", http.StatusSeeOther)

	file := ctPath + "create_contact.html"
	tmpl := template.Must(template.ParseFiles(file))

	if err := tmpl.Execute(w, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteContactHandler(w http.ResponseWriter, r *http.Request) {
	cid, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	_, err = app.store.GetContactById(r.Context(), cid)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.store.DeleteContact(r.Context(), cid); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/contacts/dashboard", http.StatusSeeOther)
}

func (app *application) updateContactHandler(w http.ResponseWriter, r *http.Request) {
	cid, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		contact, err := app.store.GetContactForUpdate(r.Context(), cid)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				app.notFoundResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		file := ctPath + "edit_contact.html"
		tmpl := template.Must(template.ParseFiles(file))

		if err := tmpl.Execute(w, contact); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	case http.MethodPatch:
		if err := r.ParseForm(); err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		data := db.UpdateContactParams{
			ID: cid,
			FirstName: pgtype.Text{
				String: r.FormValue("first_name"),
				Valid:  true,
			},
			LastName: pgtype.Text{
				String: r.FormValue("last_name"),
				Valid:  true,
			},
			Phone: pgtype.Text{
				String: r.FormValue("phone"),
				Valid:  true,
			},
		}

		_, err := app.store.UpdateContact(r.Context(), data)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				app.notFoundResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		http.Redirect(w, r, "/contacts/dashboard", http.StatusSeeOther)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		app.logger.Info("Using method", r.Method, "<<<<<<<")
		return
	}
}

func (app *application) listContactsHandler(w http.ResponseWriter, r *http.Request) {
	var filters db.Filters
	qs := r.URL.Query()

	filters.Page = app.readInt(qs, "page", 1)
	filters.PageSize = app.readInt(qs, "page_size", 3)

	pgnt := db.ListContactsParams{
		Limit:  filters.PageSize,
		Offset: (filters.Page - 1) * filters.PageSize,
	}

	contacts, err := app.store.ListContacts(r.Context(), pgnt)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	file := ctPath + "lad_contact.html"
	tmpl := template.Must(template.ParseFiles(file))

	if err := tmpl.Execute(w, envelope{"Contacts": contacts}); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
