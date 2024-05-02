package main

import (
	"net/http"

	"github.com/Rich-Wilkyness/bookings/internal/config"
	"github.com/Rich-Wilkyness/bookings/internal/handlers"

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
	mux.Get("/generals-quaters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.PostAvailabilityJSON)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/make-reservation", handlers.Repo.Reservation)

	// this allows our tmpl templates to access our static directory
	// this directory is where we will store things like images
	fileServer := http.FileServer(http.Dir("./static/"))             // we first find our directory. by using "./" this is our root. and this is what is required by the Dir function
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // then we feed our mux the directory by directing it to our static directory and removing static from the pathname of our files to get our filename
	return mux
}
