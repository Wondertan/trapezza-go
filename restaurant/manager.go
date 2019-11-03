package restaurant

import (
	"context"
	"errors"
	"sync"

	"github.com/Wondertan/trapezza-go/trapezza"
)

var (
	ErrTableTaken = errors.New("session: table already has started session")
)

type sessionManager interface {
	NewSession() (*trapezza.Session, error)
	EndSession(string) error
	Session(string) (*trapezza.Session, error)
}

type Manager struct {
	sessionManager sessionManager

	sessions map[string]map[string]*trapezza.Session // map[Restaurant][Table]id

	subs          map[string]map[int]chan Event
	subReqs       chan *subReq
	subCancelReqs chan *subCancelReq

	l   sync.RWMutex
	ctx context.Context
}

func NewManager(ctx context.Context, sessions sessionManager) *Manager {
	man := &Manager{
		sessionManager: sessions,
		sessions:       make(map[string]map[string]*trapezza.Session),
		subs:           make(map[string]map[int]chan Event),
		subReqs:        make(chan *subReq),
		subCancelReqs:  make(chan *subCancelReq),
		ctx:            ctx,
	}

	go man.handleSubs()
	return man
}

func (man *Manager) NewTrapezzaSession(rest, table string) (*trapezza.Session, error) {
	man.l.Lock()
	defer man.l.Unlock()

	sess, ok := man.sessions[rest]
	if !ok {
		sess = make(map[string]*trapezza.Session)
		man.sessions[rest] = sess
	}

	_, ok = sess[table]
	if ok {
		return nil, ErrTableTaken
	}

	ses, err := man.sessionManager.NewSession()
	if err != nil {
		return nil, err
	}

	go man.publish(&NewTrapezzaSessionEvent{ses.ID(), table, rest})

	sess[table] = ses
	return ses, nil
}

func (man *Manager) EndTrapezzaSession(rest, table string) error {
	man.l.Lock()
	defer man.l.Unlock()

	sess, ok := man.sessions[rest]
	if !ok {
		return trapezza.ErrNotFound
	}

	ses, ok := sess[table]
	if !ok {
		return trapezza.ErrNotFound
	}

	err := man.sessionManager.EndSession(ses.ID())
	if err != nil {
		return err
	}

	go man.publish(&EndTrapezzaSessionEvent{ses.ID(), table, rest})

	delete(sess, table)
	return nil
}

func (man *Manager) TrapezzaSession(rest, table string) (*trapezza.Session, error) {
	man.l.RLock()
	defer man.l.RUnlock()

	sess, ok := man.sessions[rest]
	if !ok {
		return nil, trapezza.ErrNotFound
	}

	ses, ok := sess[table]
	if !ok {
		return nil, trapezza.ErrNotFound
	}

	return ses, nil
}

func (man *Manager) SubscribeEvents(ctx context.Context, restaurant string) (<-chan Event, error) {
	select {
	case <-man.ctx.Done():
		return nil, man.ctx.Err()
	default:
	}

	sub := make(chan Event, 1)
	select {
	case man.subReqs <- &subReq{ctx, sub, restaurant}:
		return sub, nil
	case <-man.ctx.Done():
		return nil, ctx.Err()
	case <-man.ctx.Done():
		return nil, ctx.Err()
	}
}

func (man *Manager) publish(event Event) {
	for _, sub := range man.subs[event.Restaurant()] {
		select {
		case sub <- event:
		default:
			// subscription is not listening
			// TODO Log
		}
	}
}

func (man *Manager) handleClose(ctx context.Context, id int, restaurant string) {
	select {
	case <-ctx.Done():
		man.subCancelReqs <- &subCancelReq{id, restaurant}
	case <-man.ctx.Done():
	}
}

func (man *Manager) handleSubs() {
	for {
		select {
		case req := <-man.subReqs:
			subs := man.subs[req.restaurant]
			if subs == nil {
				subs = make(map[int]chan Event)
				man.subs[req.restaurant] = subs
			}

			id := len(subs)
			subs[id] = req.sub
			go man.handleClose(req.ctx, id, req.restaurant)
		case req := <-man.subCancelReqs:
			close(man.subs[req.restaurant][req.id])
			delete(man.subs[req.restaurant], req.id)
		}
	}
}

type subReq struct {
	ctx context.Context
	sub chan Event

	restaurant string
}

type subCancelReq struct {
	id         int
	restaurant string
}
