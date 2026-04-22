package main

import (
	"log"
	"net/http"
	jwt "pet-care/internal/JWT"
	"pet-care/service"
	Store "pet-care/store"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type dbconfig struct {
	Addr        string
	MaxIdlecons int
	MaxOpencons int
	MaxIdletime string
}

type config struct {
	Addr     string
	DBconfig dbconfig
}

type Application struct {
	Config    config
	Store     Store.Storage
	Service   service.Services
	TokenUtil jwt.TokenUtil
}

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowedOrigins:   []string{"https//", "https//*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Use(middleware.Timeout(60 * time.Second))

	return r
}

func (app *Application) Run(MUX http.Handler) error {

	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      MUX,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server berjalan di %v", srv)
	return srv.ListenAndServe()
}
