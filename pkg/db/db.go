package db

import (
	"github.com/Iowel/app-base-server/configs"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	Pool *pgxpool.Pool
}

func NewDB(conf configs.Config) (*Db, error) {
	pool, err := pgxpool.New(context.Background(), conf.DB.Dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Connected to PostgreSQL")

	return &Db{
		Pool: pool,
	}, nil

}
