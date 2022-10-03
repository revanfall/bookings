package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/revanfall/bookings/driver"
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
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Sql.Close()

	defer close(app.MailChan)
	listenForMail()
	fmt.Println("Starting mail listener...")
	fmt.Printf("Starting server at %v", portNum)

	srv := &http.Server{Addr: portNum, Handler: routes(&app)}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// store value in a session

	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	//in production
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()
	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = *inProduction
	app.UseCache = *useCache

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

	//connect to db
	log.Println("Connecting to db")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		*dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("cannot connect to db ", err)
	}
	log.Println("Connected to db")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	render.NewRenderer(&app)
	return db, nil
}
