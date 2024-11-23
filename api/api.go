package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/devtiagoneves/api-restful/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(db *pkg.Application) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Route("/api", func(r chi.Router) {
		r.Post("/users", createUser(db))
		r.Get("/users", allUsers(db))
		r.Get("/users/{id}", findByIdUser(db))
		r.Delete("/users/{id}", deleteUser(db))
		r.Put("/users/{id}", updateUser(db))
	})

	return r
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type InputCreateBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type InputUpdateBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Biography string `json:"biography"`
}

func createUser(db *pkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body InputCreateBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			pkg.SendJSON(
				w,
				pkg.Response{Error: "Please provide FirstName LastName and bio for the user"},
				http.StatusBadRequest,
			)
			slog.Error("NewDecoder", "error", err)
			return
		}

		newUser, err := db.Insert(body.FirstName, body.LastName, body.Biography)
		if err != nil {
			pkg.SendJSON(
				w,
				pkg.Response{Error: "There was an error while saving the user to the database"},
				http.StatusInternalServerError,
			)
			return
		}

		pkg.SendJSON(
			w,
			pkg.Response{Data: newUser},
			http.StatusCreated,
		)

	}
}

func allUsers(db *pkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := db.FindAll()
		pkg.SendJSON(w, pkg.Response{
			Data: users,
		}, http.StatusOK)
	}
}

func findByIdUser(db *pkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		user := db.FindById(id)

		pkg.SendJSON(w, pkg.Response{
			Data: user,
		}, http.StatusOK)

	}
}

func updateUser(db *pkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body InputUpdateBody
		id := chi.URLParam(r, "id")

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			pkg.SendJSON(
				w,
				pkg.Response{Error: "Please provide FirstName and bio for the user"},
				http.StatusUnprocessableEntity,
			)
			slog.Error("NewDecoder", "error", err)
			return
		}

		err := db.Update(id, body.FirstName, body.LastName, body.Biography)
		if err != nil {
			pkg.SendJSON(
				w,
				pkg.Response{Error: err.Error()},
				http.StatusInternalServerError,
			)
			return
		}

		pkg.SendJSON(w, pkg.Response{}, http.StatusOK)

	}
}

func deleteUser(db *pkg.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		err := db.Delete(id)
		if err != nil {
			pkg.SendJSON(
				w,
				pkg.Response{Error: err.Error()},
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
