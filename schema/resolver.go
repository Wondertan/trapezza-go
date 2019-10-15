//go:generate go run github.com/99designs/gqlgen

package schema

import (
	"context"

	"github.com/Wondertan/trapezza-go/session"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	*session.Manager
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) NewSession(ctx context.Context, waiter string, table string) (string, error) {
	id := r.Manager.NewSession()

	err := r.Manager.SetWaiter(ctx, id, waiter)
	if err != nil {
		return "", err
	}

	err = r.Manager.SetTable(ctx, id, table)
	if err != nil {
		return "", err
	}

	return string(id), nil
}

func (r *mutationResolver) EndSession(_ context.Context, ses string) (bool, error) {
	return true, r.Manager.EndSession(session.ID(ses))
}

func (r *mutationResolver) AddClient(ctx context.Context, ses string, client string) (bool, error) {
	return true, r.Manager.AddClient(ctx, session.ID(ses), client)
}

func (r *mutationResolver) AddItem(ctx context.Context, ses string, client string, item string) (bool, error) {
	return true, r.Manager.AddItem(ctx, session.ID(ses), client, item)
}

func (r *mutationResolver) SetWaiter(ctx context.Context, ses string, waiter string) (bool, error) {
	return true, r.Manager.SetWaiter(ctx, session.ID(ses), waiter)
}

func (r *mutationResolver) SetTable(ctx context.Context, ses string, table string) (bool, error) {
	return true, r.Manager.SetTable(ctx, session.ID(ses), table)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Session(ctx context.Context, id string) (*Session, error) {
	state, err := r.Manager.State(ctx, session.ID(id))
	if err != nil {
		return nil, err
	}

	var orders []*Order
	for client, items := range state.Orders {
		out := make([]*string, len(items))
		for i, item := range items {
			out[i] = &item
		}

		orders = append(orders, &Order{Client: &client, Items: out})
	}

	return &Session{
		ID:     &id,
		Waiter: &state.Waiter,
		Orders: orders,
		Table:  &state.Table,
	}, nil
}
