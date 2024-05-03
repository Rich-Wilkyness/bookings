package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Rich-Wilkyness/bookings/internal/config"
	"github.com/Rich-Wilkyness/bookings/internal/forms"
	"github.com/Rich-Wilkyness/bookings/internal/models"
	"github.com/Rich-Wilkyness/bookings/internal/render"
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

	render.RenderTemplateAdvanced(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip") // we are accessing the session we made on our home handler
	stringMap["remote_ip"] = remoteIP                             // we then store our session information in the stringmap to pass it to our frontend

	render.RenderTemplateAdvanced(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap, // we are sending type StringMap the data in stringMap
	})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation // this is an empty reservation, by creating this with the data, we can populate our form with information that they previously input
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplateAdvanced(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil), // creates a new form without anything in it when the make-reservation page is loaded.
		Data: data,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	form := forms.New(r.PostForm)

	// form.Has("first_name", r) // our Has function will add any errors to our struct for a specific field
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	// this is so we can repopulate the form on errors, so they do not have to redo the form
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplateAdvanced(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	// we want to redirect so our user does not submit twice (this is crucial for financial posts)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther) // we can use other status' for redirect, but this one works well enough (303)
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplateAdvanced(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplateAdvanced(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplateAdvanced(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start") // we get these values from the form, this is the syntax to get them from the request "r" value. these are returned as strings.
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end))) // have to wrap it in a print
}

type jsonResponse struct {
	OK      bool   `json:"ok"` // first) our attributes need to be capital so they are accessible. Second) we have to declare what we want the json to look like
	Message string `json:"message"`
}

func (m *Repository) PostAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{ // the response to the request being made
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ") // generates the json from our passed in object, indent is just formatting
	if err != nil {
		log.Println(err)
	}
	// create a header that tells the web browser what kind of response this is
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplateAdvanced(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) // the .(models.Reservation) comes from us Registering the type on our main. This is asserting what type we are expecting to get from our Session
	if !ok {
		// what if someone tries to access this page by typing in the url? We don't want this, so we are going to redirect them to another page
		log.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	// we need to remove our reservation from our session so it does not persist. we are currently done with it at this point.
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplateAdvanced(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
