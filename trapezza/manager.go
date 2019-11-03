package trapezza

import (
	"context"
	"errors"
	"sync"

	"github.com/Wondertan/trapezza-go/utils"
)

const IDLength = 5

var (
	ErrNotFound = errors.New("trapezza: not found")
)

type Manager struct {
	sessions map[string]*Session // to allow access through IDs without traversing sessions map

	ctx context.Context
	l   sync.RWMutex
}

func NewManager(ctx context.Context) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		ctx:      ctx,
	}
}

func (man *Manager) NewSession() (*Session, error) {
	// TODO Check that ID is unique
	id := utils.RandString(IDLength)

	man.l.Lock()
	defer man.l.Unlock()

	ses := newSession(man.ctx, id)
	man.sessions[id] = ses
	return ses, nil
}

func (man *Manager) EndSession(id string) error {
	man.l.Lock()
	defer man.l.Unlock()

	ses, ok := man.sessions[id]
	if !ok {
		return ErrNotFound
	}

	delete(man.sessions, ses.ID())
	ses.stop()
	return nil
}

func (man *Manager) Session(id string) (*Session, error) {
	man.l.RLock()
	defer man.l.RUnlock()

	ses, ok := man.sessions[id]
	if !ok {
		return nil, ErrNotFound
	}

	return ses, nil
}
