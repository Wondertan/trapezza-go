package resolver

import (
	"context"

	"github.com/Wondertan/trapezza-go/session"
)

type mutation Resolver

func (r *mutation) NewSession(ctx context.Context, waiter string, table string) (session.ID, error) {
	id := r.Manager.NewSession()

	err := r.Manager.SetWaiter(ctx, id, waiter)
	if err != nil {
		return "", err
	}

	err = r.Manager.SetTable(ctx, id, table)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *mutation) EndSession(_ context.Context, session session.ID) (bool, error) {
	return true, r.Manager.EndSession(session)
}

func (r *mutation) AddClient(ctx context.Context, session session.ID, client string) (bool, error) {
	return true, r.Manager.AddClient(ctx, session, client)
}

func (r *mutation) AddItem(ctx context.Context, session session.ID, client string, item string) (bool, error) {
	return true, r.Manager.AddItem(ctx, session, client, item)
}

func (r *mutation) SetWaiter(ctx context.Context, session session.ID, waiter string) (bool, error) {
	return true, r.Manager.SetWaiter(ctx, session, waiter)
}

func (r *mutation) SetTable(ctx context.Context, session session.ID, table string) (bool, error) {
	return true, r.Manager.SetTable(ctx, session, table)
}