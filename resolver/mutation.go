package resolver

import (
	"context"

	"github.com/Wondertan/trapezza-go/trapezza"
	"github.com/Wondertan/trapezza-go/types"
)

type mutation Resolver

func (m *mutation) New(_ context.Context, rest string, table string) (string, error) {
	return m.trapezza.NewSession(rest, table)
}

func (m *mutation) ChangeWaiter(ctx context.Context, session, waiter string) (bool, error) {
	err := m.emitEvent(ctx, session, &trapezza.ChangeWaiterEvent{Waiter: waiter})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *mutation) NewGroup(ctx context.Context, session, payer string) (bool, error) {
	err := m.emitEvent(ctx, session, &trapezza.NewGroupEvent{Payer: payer})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *mutation) JoinGroup(ctx context.Context, session, client, payer string) (bool, error) {
	err := m.emitEvent(ctx, session, &trapezza.JoinGroupEvent{Client: client, Payer: payer})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *mutation) AddItems(ctx context.Context, session, client string, ids []string) (bool, error) {
	items := make([]*types.Item, len(ids))
	for i, id := range ids {
		items[i] = &types.Item{
			Id:    id,
			Price: 10, // TODO Get price from DB
		}
	}

	err := m.emitEvent(ctx, session, &trapezza.AddItemsEvent{Client: client, Items: items})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *mutation) RemoveItem(ctx context.Context, session, client, item string) (bool, error) {
	err := m.emitEvent(ctx, session, &trapezza.RemoveItemEvent{Client: client, Item: item})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *mutation) SplitItem(ctx context.Context, session, who, with, item string) (bool, error) {
	err := m.emitEvent(ctx, session, &trapezza.SplitItemEvent{Who: who, With: with, Item: item})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *mutation) ChangePayer(ctx context.Context, session, payer string) (bool, error) {
	err := m.emitEvent(ctx, session, &trapezza.ChangePayerEvent{Payer: payer})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *mutation) CheckoutPayer(ctx context.Context, session, payer string) (bool, error) {
	err := m.emitEvent(ctx, session, &trapezza.CheckoutPayerEvent{Payer: payer})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *mutation) CheckoutClient(ctx context.Context, session, client string) (bool, error) {
	err := m.emitEvent(ctx, session, &trapezza.CheckoutClientEvent{Client: client})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *mutation) WaiterCall(ctx context.Context, session, client, message string) (bool, error) {
	err := m.emitEvent(ctx, session, &trapezza.WaiterCallEvent{Client: client, Message: message})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *mutation) emitEvent(ctx context.Context, session string, event trapezza.Event) error {
	s, err := m.trapezza.SessionByID(session)
	if err != nil {
		return err
	}

	return s.EmitEvent(ctx, event)
}
