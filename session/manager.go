package session

import (
	"context"
	"errors"
	"sync"
)

const IDLength = 5

var ErrNotFound = errors.New("session: not found")

type Manager struct {
	ctx context.Context

	l        sync.RWMutex
	sessions map[ID]*Session
}

func NewManager(ctx context.Context) *Manager {
	return &Manager{
		ctx:      ctx,
		sessions: make(map[ID]*Session),
	}
}

func (man *Manager) NewSession() ID {
	// TODO Check that ID is unique
	id := RandID(IDLength)

	man.l.Lock()
	man.sessions[id] = newSession(man.ctx, id)
	man.l.Unlock()

	return id
}

func (man *Manager) EndSession(id ID) error {
	ses, err := man.Session(id)
	if err != nil {
		return err
	}

	man.l.Lock()
	delete(man.sessions, id)
	man.l.Unlock()

	ses.Stop()
	return nil
}

func (man *Manager) Session(id ID) (*Session, error) {
	man.l.RLock()
	defer man.l.RUnlock()

	ses, ok := man.sessions[id]
	if !ok {
		return nil, ErrNotFound
	}

	return ses, nil
}

func (man *Manager) State(ctx context.Context, id ID) (*State, error) {
	ses, err := man.Session(id)
	if err != nil {
		return nil, err
	}

	ch, err := ses.State()
	if err != nil {
		return nil, err
	}

	select {
	case state := <-ch:
		return state, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (man *Manager) AddItem(ctx context.Context, id ID, client, item string) error {
	ses, err := man.Session(id)
	if err != nil {
		return err
	}

	return EmitEvent(ctx, ses, &itemEvent{client: client, item: item})
}

func (man *Manager) AddClient(ctx context.Context, id ID, client string) error {
	ses, err := man.Session(id)
	if err != nil {
		return err
	}

	return EmitEvent(ctx, ses, &clientEvent{client: client})
}

func (man *Manager) SetWaiter(ctx context.Context, id ID, waiter string) error {
	ses, err := man.Session(id)
	if err != nil {
		return err
	}

	return EmitEvent(ctx, ses, &waiterEvent{waiter: waiter})
}

func (man *Manager) SetTable(ctx context.Context, id ID, table string) error {
	ses, err := man.Session(id)
	if err != nil {
		return err
	}

	return EmitEvent(ctx, ses, &tableEvent{table: table})
}

func EmitEvent(ctx context.Context, ses *Session, event Event) error {
	resp, err := ses.Event(event)
	if err != nil {
		return err
	}

	select {
	case err := <-resp:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
