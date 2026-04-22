package main

import (
	"log"
	jwt "pet-care/internal/JWT"
	"pet-care/service"
	"pet-care/setting/db"
	"pet-care/setting/env"
	Store "pet-care/store"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config{
		Addr: env.GetString("ADDR", "8080"),
		DBconfig: dbconfig{
			Addr:        env.GetString("postgres://petcare_user:petcare123@localhost:5432/petcare_db?sslmode=disable", "ADDR"),
			MaxIdlecons: env.GetInt("Idle", 30),
			MaxOpencons: env.GetInt("cons", 30),
			MaxIdletime: env.GetString("time", "15m"),
		},
	}

	db, err := db.New(
		cfg.DBconfig.Addr,
		cfg.DBconfig.MaxIdlecons,
		cfg.DBconfig.MaxOpencons,
		cfg.DBconfig.MaxIdletime,
	)

	if err != nil {
		log.Panic(err)
	}

	tokenUtil := jwt.NewTokenUtil(env.GetString("rahasiii123", "rahasiii123"))

	store := Store.NewStorage(db)
	service := service.Services{
		StoreDB:   store,
		TokenUtil: tokenUtil,
	}

	app := Application{
		Config:    cfg,
		Store:     store,
		Service:   service,
		TokenUtil: tokenUtil,
	}

	mux := app.Mount()
	log.Fatal(app.Run(mux))

}
