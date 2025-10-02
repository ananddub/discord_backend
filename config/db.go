package config

import (
	"context"
	"discord/gen/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.Database.PostgreSQL.URL)
	if err != nil {
		return err
	}

	poolCfg.MinConns = int32(cfg.Database.PostgreSQL.MinConnections)
	poolCfg.MaxConns = int32(cfg.Database.PostgreSQL.MaxConnections)

	DB, err = pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return err
	}
	return DB.Ping(context.Background())
}

type Queries struct {
	WriteDb *pgxpool.Pool
	ReadDb  *pgxpool.Pool
	Repo    *repo.Queries
}

func RepoQuieries() (*Queries, error) {
	return &Queries{
		WriteDb: GetWriteDb(),
		ReadDb:  GetReadDb(),
		Repo:    repo.New(DB),
	}, nil
}

func GetReadDb() *pgxpool.Pool {
	return DB
}

func GetWriteDb() *pgxpool.Pool {
	return DB
}
