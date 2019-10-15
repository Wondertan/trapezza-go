package session

import (
	"context"
)

type ID string

// TODO Context handling
type Session struct {
	state *State

	eventReqs chan *eventReq
	stateReqs chan chan *State
	subReqs   chan chan Event
	subs      []chan Event

	ctx    context.Context
	cancel context.CancelFunc
}

func (ses *Session) ID() ID {
	return ses.state.Id
}

func (ses *Session) Event(e Event) (chan error, error) {
	resp := make(chan error, 1)
	select {
	case ses.eventReqs <- &eventReq{event: e, resp: resp}:
		return resp, nil
	case <-ses.ctx.Done():
		return nil, ses.ctx.Err()
	}
}

// TODO Subscription closing
func (ses *Session) SubscribeEvents() (<-chan Event, error) {
	sub := make(chan Event, 1)
	select {
	case ses.subReqs <- sub:
		return sub, nil
	case <-ses.ctx.Done():
		return nil, ses.ctx.Err()
	}
}

func (ses *Session) State() (chan *State, error) {
	resp := make(chan *State, 1)
	select {
	case ses.stateReqs <- resp:
		return resp, nil
	case <-ses.ctx.Done():
		return nil, ses.ctx.Err()
	}
}

func (ses *Session) Stop() {
	// TODO Be more graceful
	ses.cancel()
}

func newSession(ctx context.Context, id ID) *Session {
	ctx, cancel := context.WithCancel(ctx)
	ses := &Session{
		state: &State{
			Id: id,
		},
		eventReqs: make(chan *eventReq),
		stateReqs: make(chan chan *State),
		subReqs:   make(chan chan Event),
		subs:      make([]chan Event, 0),
		ctx:       ctx,
		cancel:    cancel,
	}

	go ses.handle()
	return ses
}

func (ses *Session) handle() {
	// TODO Cleaning on defer?
	for {
		select {
		case req := <-ses.stateReqs:
			req <- ses.state
		case req := <-ses.eventReqs:
			err := req.event.Handle(ses.state)
			if err != nil {
				req.resp <- err
				continue
			}
			req.resp <- nil

			for _, sub := range ses.subs {
				select {
				case sub <- req.event:
				default:
					// subscriber is not listening
					// TODO Log
				}
			}
		case sub := <-ses.subReqs:
			ses.subs = append(ses.subs, sub)
		case <-ses.ctx.Done():
			return
		}
	}
}

type eventReq struct {
	event Event
	resp  chan error
}
