package mypg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	*pgxpool.Pool
}

func NewPostgresDB(db *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{
		Pool: db,
	}
}

func (p *PostgresDB) CloseDB() {
	p.Close()
}

func (p *PostgresDB) GetPool() *pgxpool.Pool {
	return p.Pool
}

func (p *PostgresDB) Ping(ctx context.Context) error {
	return p.Pool.Ping(ctx)
}

func (p *PostgresDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return p.Pool.Exec(ctx, sql, args...)

}

func (p *PostgresDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	v, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return v, err
	}
	fmt.Println("query:", sql, args)
	go ReactiveToPublish(sql, v, args...)
	return v, err

}

func (p *PostgresDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	v := p.Pool.QueryRow(ctx, sql, args...)
	fmt.Println("query:", sql, args)
	go ReactiveToPublish(sql, v, args...)
	return v
}
