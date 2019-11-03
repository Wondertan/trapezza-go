package trapezza

import (
	"context"
	"errors"
	"sync"

	"github.com/Wondertan/trapezza-go/utils"
)

const IDLength = 5

var (
	ErrNotFound   = errors.New("session: not found")
	ErrTableTaken = errors.New("session: table already has started session")
)

type Manager struct {
	ids      map[string]*Session            // to allow access through IDs without traversing sessions map
	sessions map[string]map[string]*Session // map[Restaurant][Table]

	ctx context.Context
	l   sync.RWMutex
}

func NewManager(ctx context.Context) *Manager {
	return &Manager{
		ids:      make(map[string]*Session),
		sessions: make(map[string]map[string]*Session),
		ctx:      ctx,
	}
}

func (man *Manager) NewSession(rest, table string) (string, error) {
	// TODO Check that ID is unique in restaurant
	id := utils.RandString(IDLength)

	man.l.Lock()
	defer man.l.Unlock()

	sess, ok := man.sessions[rest]
	if !ok {
		sess = make(map[string]*Session)
		man.sessions[rest] = sess
	}

	_, ok = sess[table]
	if ok {
		return "", ErrTableTaken
	}

	ses := newSession(man.ctx, id)
	sess[table] = ses
	man.ids[id] = ses
	return id, nil
}

func (man *Manager) EndSession(rest, table string) error {
	man.l.Lock()
	defer man.l.Unlock()

	sess, ok := man.sessions[rest]
	if !ok {
		return ErrNotFound
	}

	ses, ok := sess[table]
	if !ok {
		return ErrNotFound
	}

	delete(sess, table)
	delete(man.ids, ses.ID())
	ses.stop()
	return nil
}

func (man *Manager) Session(rest, table string) (*Session, error) {
	man.l.RLock()
	defer man.l.RUnlock()

	sess, ok := man.sessions[rest]
	if !ok {
		return nil, ErrNotFound
	}

	ses, ok := sess[table]
	if !ok {
		return nil, ErrNotFound
	}

	return ses, nil
}

func (man *Manager) SessionByID(id string) (*Session, error) {
	man.l.RLock()
	defer man.l.RUnlock()

	ses, ok := man.ids[id]
	if !ok {
		return nil, ErrNotFound
	}

	return ses, nil
}
