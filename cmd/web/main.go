package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Rich-Wilkyness/bookings/internal/config"
	"github.com/Rich-Wilkyness/bookings/internal/handlers"
	"github.com/Rich-Wilkyness/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

// Important: Windows Users
// In the next few lectures, lecture, I run multiple Go files at the same time on a Macintosh like this:
// "go run *.go"
// On Windows, though, this will not work unless you have customized your IDE or terminal. Instead, use this command:
// "go run ."

const portNumber = ":8080"

var app config.AppConfig // We want our global template cache here. not sure why, he just said it's better
// this is global here to allow us to use the global variables in our middleware and other places.
// if it was inside main we couldn't call it in our middleware file

var session *scs.SessionManager

func main() {

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour // lifetime is part of the session package. time.Hour is built into go. We are assigning the lifetime of the session/cookie to 24 hours
	// there are more advanced ways of storing a session like redis and a variety of databases. The above is just a default cookie on the browser / individual pc
	session.Cookie.Persist = true                  // when true, this means that after someone closes the browser/webpage the session ends
	session.Cookie.SameSite = http.SameSiteLaxMode // how strick is the site the cookie applies to. laxmode is default
	session.Cookie.Secure = app.InProduction       // for production we want this to be true, when true this makes the site https. and the cookies are encrypted

	app.Session = session

	tc, err := render.CreateTemplateCacheAdvanced()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc // add the template sets to the global cache
	app.UseCache = false

	// setup handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// this gives our render package access to the AppConfig (cache of template sets)
	// we do this after the cache has been created
	render.NewTemplates(&app) // needs a pointer, so we reference the app with &

	// We don't need our routes here anymore, in a bigger system this over complicates our main file.
	// instead we made a seperate routes file in the main package.
	// http is a built in library to access the web
	// http.HandleFunc("/", handlers.Repo.Home) // first variable is the pathname,
	// http.HandleFunc("/about", handlers.Repo.About)

	// n, err := fmt.Fprintf(w, "Hello, world!") // n is the number of bytes printed
	// if err != nil { // this means if there is an error (error is not nothing)
	// 	fmt.Println(err)
	// }
	// fmt.Println(fmt.Sprintf("Number of Bytes written: %d", n))

	fmt.Printf("Starting application on port %s", portNumber) // %s is a placeholder for a string

	// start a web server that is listening
	// the _ = means that if there is an error with this, ignore it (this is ok here, because we'll know if the server does not start)
	// we made a new serve that handles our routes
	// not sure the benefits yet
	// _ = http.ListenAndServe(portNumber, nil) // it expects a port and a handler (we use nil since we have our Handler function above)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app), // we point to our app
	}
	err = srv.ListenAndServe()
	log.Fatal(err)

}
