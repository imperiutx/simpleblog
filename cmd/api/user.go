package main

import (
	"net/http"
	db "simpleblog/db/sqlc"
	"simpleblog/util"
	"time"

	"github.com/go-playground/validator/v10"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type userResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Header.Get("Content-Type") {

	case "application/json":
		var req createUserRequest

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

		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		user := db.CreateUserParams{
			Username: req.Username,
			Password: hashedPassword,
			Email:    req.Email,
		}

		usr, err := app.store.CreateUser(r.Context(), user)
		if err != nil {
			if db.ErrorCode(err) == db.UniqueViolation {
				app.forbiddenResponse(w, r, err)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		rsp := newUserResponse(usr)

		if err = app.writeJSON(
			w,
			http.StatusAccepted,
			envelope{"user": rsp}, nil); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		//TODO: handle more

	default:
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}

}
