package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/revanfall/bookings/internal/config"
	"github.com/revanfall/bookings/internal/forms"
	"github.com/revanfall/bookings/internal/models"
	"github.com/revanfall/bookings/internal/render"
	"log"
	"net/http"
)

// Repo is the repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{App: a}
}

// NewHandlers sets the repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{}, r)
}

// Generals room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{}, r)
}

// Majors room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "majors.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil), Data: data},
		r)
}

// PostReservation handles the posting of the reservation form
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

	//form.Has("first_name", r)
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data}, r)
	} else {
		m.App.Session.Put(r.Context(), "reservation", reservation)
		http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
	}
}

// Majors room page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send json response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      false,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Bye again"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	}, r)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		log.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation for session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	}, r)
}
