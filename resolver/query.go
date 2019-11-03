package resolver

import (
	"context"

	"github.com/Wondertan/trapezza-go/types"
)

type query Resolver

func (r *query) Session(ctx context.Context, rest, table string) (*types.Trapezza, error) {
	s, err := r.trapezza.Session(rest, table)
	if err != nil {
		return nil, err
	}

	return s.Trapezza(ctx)
}

func (r *query) SessionByID(ctx context.Context, id string) (*types.Trapezza, error) {
	s, err := r.trapezza.SessionByID(id)
	if err != nil {
		return nil, err
	}

	return s.Trapezza(ctx)
}
