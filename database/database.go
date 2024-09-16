package database

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

// create a New Pgxpool with the singleton pattern
func NewPG(ctx context.Context, dbUrl string) (*Postgres, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, dbUrl)
		if err != nil {
			log.Fatalf("Unable to create connection pool %v", err)
			return
		}
		pgInstance = &Postgres{db}
	})
	return pgInstance, nil
}

func (pg *Postgres) Getpool() *pgxpool.Pool {
	return pg.db
} // get the pool

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
} //self explained by the name

func (pg *Postgres) Close() {
	pg.db.Close()
}
