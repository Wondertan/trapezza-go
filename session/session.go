package session

import (
	"context"
)

type EventType string

type Event interface {
	Type() EventType
}

type State interface {
	Handle(Event) error
}

type Update struct {
	State
	Event
}

type Session struct {
	state State

	updateSub       chan *updateSubReq
	updateSubCancel chan uint64
	updateSubs      map[uint64]chan *Update
	eventReqs       chan *eventReq
	stateReqs       chan chan State

	n uint64

	ctx    context.Context
	cancel context.CancelFunc
}

func NewSession(ctx context.Context, state State) *Session {
	ctx, cancel := context.WithCancel(ctx)
	ses := &Session{
		state:           state,
		updateSub:       make(chan *updateSubReq),
		updateSubCancel: make(chan uint64),
		updateSubs:      make(map[uint64]chan *Update),
		eventReqs:       make(chan *eventReq),
		stateReqs:       make(chan chan State),
		ctx:             ctx,
		cancel:          cancel,
	}

	go ses.handle()
	return ses
}

func (ses *Session) EmitEvent(e Event) (<-chan error, error) {
	select {
	case <-ses.ctx.Done():
		return nil, ses.ctx.Err()
	default:
	}

	resp := make(chan error, 1)
	select {
	case ses.eventReqs <- &eventReq{event: e, resp: resp}:
		return resp, nil
	case <-ses.ctx.Done():
		return nil, ses.ctx.Err()
	}
}

func (ses *Session) SubscribeUpdates(ctx context.Context) (<-chan *Update, error) {
	select {
	case <-ses.ctx.Done():
		return nil, ses.ctx.Err()
	default:
	}

	ch := make(chan *Update, 1)
	select {
	case ses.updateSub <- &updateSubReq{ctx: ctx, sub: ch}:
		return ch, nil
	case <-ctx.Done():
		return nil, ses.ctx.Err()
	case <-ses.ctx.Done():
		return nil, ses.ctx.Err()
	}
}

func (ses *Session) State() (<-chan State, error) {
	select {
	case <-ses.ctx.Done():
		return nil, ses.ctx.Err()
	default:
	}

	resp := make(chan State, 1)
	select {
	case ses.stateReqs <- resp:
		return resp, nil
	case <-ses.ctx.Done():
		return nil, ses.ctx.Err()
	}
}

func (ses *Session) Stop() {
	ses.cancel()
}

func (ses *Session) handleClose(ctx context.Context, id uint64) {
	select {
	case <-ctx.Done():
		ses.updateSubCancel <- id
	case <-ses.ctx.Done():
	}
}

func (ses *Session) handle() {
	for {
		select {
		case req := <-ses.updateSub:
			go ses.handleClose(req.ctx, ses.n)
			ses.updateSubs[ses.n] = req.sub
			ses.n++
		case id := <-ses.updateSubCancel:
			close(ses.updateSubs[id])
			delete(ses.updateSubs, id)
		case req := <-ses.eventReqs:
			err := ses.state.Handle(req.event)
			if err != nil {
				req.resp <- err
				continue
			}
			req.resp <- nil

			for _, sub := range ses.updateSubs {
				select {
				case sub <- &Update{
					State: ses.state,
					Event: req.event,
				}:
				default:
					// subscriber is not listening
					// TODO Log
				}
			}
		case req := <-ses.stateReqs:
			req <- ses.state
			close(req)
		case <-ses.ctx.Done():
			return
		}
	}
}

type updateSubReq struct {
	ctx context.Context
	sub chan *Update
}

type eventReq struct {
	event Event
	resp  chan error
}
