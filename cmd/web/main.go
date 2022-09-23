package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/revanfall/bookings/internal/config"
	"github.com/revanfall/bookings/internal/handlers"
	"github.com/revanfall/bookings/internal/helpers"
	"github.com/revanfall/bookings/internal/models"
	"github.com/revanfall/bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

var portNum = ":3000"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting server at %v", portNum)

	srv := &http.Server{Addr: portNum, Handler: routes(&app)}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	// store value in a session

	gob.Register(models.Reservation{})

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)

	render.NewTemplates(&app)
	return nil
}
