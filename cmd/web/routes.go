package main

import (
	"net/http"

	"github.com/Rich-Wilkyness/bookings/pkg/config"
	"github.com/Rich-Wilkyness/bookings/pkg/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// we downloaded an external library called chi from github
// you can follow the github to see the install
// this library simplifies routing as we will see and provides middleware authentication
func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer) // the documentation on this is on the github
	mux.Use(NoSurf)               // this is our middleware for csrf toxen authentication
	mux.Use(SessionLoad)          // this helps us have state management

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}
