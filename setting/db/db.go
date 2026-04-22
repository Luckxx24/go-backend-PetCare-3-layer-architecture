package db

import (
	"context"
	"database/sql"
	"time"
)

func New(Addr string, MaxIdlecons, MaxOpencons int, MaxIdletime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", Addr)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(MaxOpencons)
	db.SetMaxIdleConns(MaxIdlecons)

	IdleTIme, err := time.ParseDuration(MaxIdletime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(IdleTIme)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	erro := db.PingContext(ctx)

	if erro != nil {
		return nil, erro
	}

	return db, nil
}
