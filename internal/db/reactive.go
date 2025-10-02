package pgxrx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	rxgo "github.com/reactivex/rxgo/v2"
)

type HookedPool struct {
	*pgxpool.Pool
	events chan rxgo.Item
}

// Constructor
func NewHookedPool(pool *pgxpool.Pool) *HookedPool {
	return &HookedPool{
		Pool:   pool,
		events: make(chan rxgo.Item, 100),
	}
}

// Observable to subscribe
func (p *HookedPool) Observable() rxgo.Observable {
	return rxgo.FromChannel(p.events)
}

// Exec wrapper with hooks / reactive event
func (p *HookedPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	p.events <- rxgo.Of(fmt.Sprintf("Before Exec: %s", sql))

	tag, err := p.Pool.Exec(ctx, sql, args...)

	p.events <- rxgo.Of(fmt.Sprintf("After Exec: %s", sql))

	return tag, err
}
