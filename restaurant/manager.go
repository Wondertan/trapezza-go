package restaurant

import (
	"context"
)

type sessionManager interface {
	NewSession(rest, table string) (string, error)
	EndSession(rest, table string) error
}

type Manager struct {
	sessions sessionManager

	subs          map[string]map[int]chan Event
	subReqs       chan *subReq
	subCancelReqs chan *subCancelReq

	ctx context.Context
}

func NewManager(ctx context.Context, sessions sessionManager) *Manager {
	man := &Manager{
		sessions:      sessions,
		subs:          make(map[string]map[int]chan Event),
		subReqs:       make(chan *subReq),
		subCancelReqs: make(chan *subCancelReq),
		ctx:           ctx,
	}

	go man.handleSubs()
	return man
}

func (man *Manager) NewSession(rest, table string) (string, error) {
	id, err := man.sessions.NewSession(rest, table)
	if err != nil {
		return "", err
	}

	go man.publish(&NewSessionEvent{id, table, rest})
	return id, nil
}

func (man *Manager) EndSession(rest, table string) error {
	err := man.sessions.EndSession(rest, table)
	if err != nil {
		return err
	}

	go man.publish(&EndSessionEvent{table, rest})
	return nil
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
