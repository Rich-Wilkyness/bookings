package handlers

import (
	"net/http"

	"github.com/Rich-Wilkyness/bookings/pkg/config"
	"github.com/Rich-Wilkyness/bookings/pkg/models"
	"github.com/Rich-Wilkyness/bookings/pkg/render"
)

// --------------------------------------- Setup Repositories ----------------------------------
// repository used by the handlers
var Repo *Repository

// sets the type of repository
type Repository struct {
	App *config.AppConfig
}

// creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// --------------------------------------- Handlers ------------------------------------------------
// for web browsers to work we need a response and request
// (m *Repository) gives access to the handlers everything that is inside repository
// (m *Repository) is a "reciever" - not sure what that is
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// we can now pass TemplateData to our render func
	// this Data allows us to access it on the frontend

	// grab remote ip address and store it in the session
	// the r comes from the request in our parameters. this contains information from the user being sent to us
	remoteIP := r.RemoteAddr // this is built into go, we can get the ip address. When someone makes a request it is part of the request header

	// m is accessed via the Repository struct reciever. This is the site wide config
	// the first variable = comes from the user request. not sure what Context does
	// second variable = key to save and access the info
	// third variable = information we are saving
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP) // this is how we store a session.

	render.RenderTemplateAdvanced(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip") // we are accessing the session we made on our home handler
	stringMap["remote_ip"] = remoteIP                             // we then store our session information in the stringmap to pass it to our frontend

	render.RenderTemplateAdvanced(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap, // we are sending type StringMap the data in stringMap
	})
}
