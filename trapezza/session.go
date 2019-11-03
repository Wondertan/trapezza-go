package trapezza

import (
	"context"

	"github.com/Wondertan/trapezza-go/session"
	"github.com/Wondertan/trapezza-go/types"
)

type Update struct {
	State *types.Trapezza
	Event Event
}

type Session struct {
	id  string
	ses *session.Session

	ctx context.Context
}

func newSession(ctx context.Context, id string) *Session {
	return &Session{
		id:  id,
		ses: session.NewSession(ctx, &state{types.NewTrapezza(id)}),
		ctx: ctx,
	}
}

func (ses *Session) ID() string {
	return ses.id
}

func (ses *Session) SubscribeUpdates(ctx context.Context) (<-chan *Update, error) {
	out := make(chan *Update)

	sub, err := ses.ses.SubscribeUpdates(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		for s := range sub {
			select {
			case out <- &Update{
				State: s.State.(*state).trapezza,
				Event: s.Event.(Event),
			}:
			case <-ctx.Done():
				return
			case <-ses.ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

func (ses *Session) Trapezza(ctx context.Context) (*types.Trapezza, error) {
	ch, err := ses.ses.State()
	if err != nil {
		return nil, err
	}

	select {
	case s := <-ch:
		return s.(*state).trapezza, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-ses.ctx.Done():
		return nil, ses.ctx.Err()
	}
}

func (ses *Session) EmitEvent(ctx context.Context, event Event) error {
	event.setID(ses.id)
	resp, err := ses.ses.EmitEvent(event)
	if err != nil {
		return err
	}

	select {
	case err := <-resp:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-ses.ctx.Done():
		return ses.ctx.Err()
	}
}

func (ses *Session) stop() {
	ses.ses.Stop()
}

type state struct {
	trapezza *types.Trapezza
}

func (s *state) Handle(e session.Event) error {
	switch e.Type() {
	case ChangeWaiter:
		event := e.(*ChangeWaiterEvent)
		return s.trapezza.ChangeWaiter(event.Waiter)
	case ChangePayer:
		event := e.(*ChangePayerEvent)
		return s.trapezza.ChangePayer(event.Payer)
	case NewGroupOrder:
		event := e.(*NewGroupOrderEvent)
		return s.trapezza.NewGroup(event.Payer)
	case JoinGroupOrder:
		event := e.(*JoinGroupOrderEvent)
		return s.trapezza.NewGroup(event.Payer)
	case AddItems:
		event := e.(*AddItemsEvent)
		return s.trapezza.AddItems(event.Client, event.Items)
	case RemoveItem:
		event := e.(*RemoveItemEvent)
		return s.trapezza.RemoveItem(event.Client, event.Item)
	case SplitItem:
		event := e.(*SplitItemEvent)
		return s.trapezza.SplitItem(event.Who, event.With, event.Item)
	case CheckoutClient:
		event := e.(*CheckoutClientEvent)
		return s.trapezza.CheckoutClient(event.Client)
	case CheckoutPayer:
		event := e.(*CheckoutPayerEvent)
		return s.trapezza.CheckoutPayer(event.Payer)
	case WaiterCall:
		event := e.(*WaiterCallEvent)
		return s.trapezza.WaiterCall(event.Client, event.Message)
	default:
		return ErrWrongEvent
	}
}
