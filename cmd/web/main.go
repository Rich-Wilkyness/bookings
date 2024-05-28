package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Rich-Wilkyness/bookings/internal/config"
	"github.com/Rich-Wilkyness/bookings/internal/driver"
	"github.com/Rich-Wilkyness/bookings/internal/handlers"
	"github.com/Rich-Wilkyness/bookings/internal/helpers"
	"github.com/Rich-Wilkyness/bookings/internal/models"
	"github.com/Rich-Wilkyness/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

// We want our global template cache here. not sure why, he just said it's better
// this is global here to allow us to use the global variables in our middleware and other places.
// if it was inside main we couldn't call it in our middleware file
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	// this allows us to store our type struct Reservation in a session.
	// the session library by default can do strings and ints,
	// but not specific objects or types unless we do this (where we are registering the type to be stored)
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour // lifetime is part of the session package. time.Hour is built into go. We are assigning the lifetime of the session/cookie to 24 hours
	// there are more advanced ways of storing a session like redis and a variety of databases. The above is just a default cookie on the browser / individual pc
	session.Cookie.Persist = true                  // when true, this means that after someone closes the browser/webpage the session ends
	session.Cookie.SameSite = http.SameSiteLaxMode // how strick is the site the cookie applies to. laxmode is default
	session.Cookie.Secure = app.InProduction       // for production we want this to be true, when true this makes the site https. and the cookies are encrypted

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=lebrum1203")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	// add the template sets to the global cache
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
