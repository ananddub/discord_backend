package database

import (
	"context"

	"discord/gen/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, *repo.Queries, error) {
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, nil, err
	}

	queries := repo.New(db)
	return db, queries, nil
}
