package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/revanfall/bookings/pkg/config"
	"github.com/revanfall/bookings/pkg/handlers"
	"github.com/revanfall/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

var portNum = ":3000"
var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.AboutHandler)

	fmt.Printf("Starting server at %v", portNum)

	srv := &http.Server{Addr: portNum, Handler: routes(&app)}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
