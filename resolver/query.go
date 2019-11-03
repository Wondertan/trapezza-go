package resolver

import (
	"context"

	"github.com/Wondertan/trapezza-go/types"
)

type query Resolver

func (r *query) TrapezzaSession(ctx context.Context, rest, table string) (*types.Trapezza, error) {
	s, err := r.restaurant.TrapezzaSession(rest, table)
	if err != nil {
		return nil, err
	}

	return s.Trapezza(ctx)
}

func (r *query) TrapezzaSessionByID(ctx context.Context, id string) (*types.Trapezza, error) {
	s, err := r.trapezza.Session(id)
	if err != nil {
		return nil, err
	}

	return s.Trapezza(ctx)
}
