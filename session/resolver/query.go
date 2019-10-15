package resolver

import (
	"context"

	"github.com/Wondertan/trapezza-go/session"
)

type query Resolver

func (r *query) Session(ctx context.Context, session session.ID) (*session.State, error) {
	return r.Manager.State(ctx, session)
}
