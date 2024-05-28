package main

import (
	"fmt"
	"net/http"

	"github.com/Rich-Wilkyness/bookings/internal/helpers"
	"github.com/justinas/nosurf"
)

// we are not using this code, it was an example on how to create our own middleware
func WriteToConsole(next http.Handler) http.Handler { // next is standard syntax when writing middleware
	//  name "next" is conventionally used to represent the next middleware function in the chain.
	// Middleware functions in frameworks like Express.js (Node.js) or Gorilla Mux (Go) are often composed in a stack,
	// where each middleware function in the stack has the opportunity to perform some action before calling the next middleware function in the chain.

	// When a request comes in, it is passed through a chain of middleware functions before reaching the final handler (e.g., a route handler).
	// Each middleware function in the chain has access to the request and response objects, and it can perform operations like logging, authentication, error handling, etc.
	// After performing its operation, a middleware function can either pass the request to the next middleware function in the chain or decide to respond to the request early (e.g., if authentication fails).

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // this is called an anonymous function (similar to lambda functions)
		// an anonymous function is a function without a name that can be defined inline and passed as an argument to another function.
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

// adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",                  // this means it hits the entire site
		Secure:   app.InProduction,     // in production this will be true (https) the s for secure. our global config file will change this to true when in production
		SameSite: http.SameSiteLaxMode, // standard
	})
	return csrfHandler
}

// loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler { // this will just load the session. this is important for state management
	return session.LoadAndSave(next) // hover LoadAndSave. helps remember state essentially
}

// Auth checks if a user is authenticated
func Auth(next http.Handler) http.Handler {
	// not sure how we have access to the w or r here. guess it is passed from next.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
