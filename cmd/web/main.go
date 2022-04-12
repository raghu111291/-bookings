package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/raghu/bookings/pckgs/config"
	"github.com/raghu/bookings/pckgs/handlers"
	"github.com/raghu/bookings/pckgs/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.SameSite = http.SameSiteDefaultMode
	session.Cookie.Persist = true
	session.Cookie.Secure = app.InProduction

	app.InProduction = false
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

	fmt.Println(fmt.Sprint("Starting application on port ", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
