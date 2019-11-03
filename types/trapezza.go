package types

import (
	"fmt"
	"time"
)

const (
	ItemsLimit        = 10
	GroupsLimit       = 10
	ClientsLimit      = 10
	WaiterCallTimeout = 30 * time.Second
)

var (
	ErrJoinedPayer   = fmt.Errorf("trapezza: payer already joined")
	ErrAlreadyPayer  = fmt.Errorf("trapezza: clint already payer")
	ErrWrongClient   = fmt.Errorf("trapezza: wrong client or not joined")
	ErrWrongItem     = fmt.Errorf("trapezza: wrong item")
	ErrGroupsLimit   = fmt.Errorf("trapezza: groups limit exceded")
	ErrClientsLimit  = fmt.Errorf("trapezza: clients per group limit exceded")
	ErrItemsLimit    = fmt.Errorf("trapezza: items per client limit exceded")
	ErrCallTimeout   = fmt.Errorf("trapezza: waiter call timeout")
	ErrPayerCheckout = fmt.Errorf("trapezza: payer can't checkput himself")
	ErrCheckedOut    = fmt.Errorf("trapezza: client already checked out")
)

type Item struct {
	Id    string
	Price float64
}

type Trapezza struct {
	Id       string
	Waiter   string
	Started  time.Time
	LastCall time.Time

	groups Groups
}

func NewTrapezza(id string) *Trapezza {
	return &Trapezza{
		Id:       id,
		Started:  time.Now(),
		LastCall: time.Now(),
		groups:   make(Groups, 0, GroupsLimit),
	}
}

// TODO Needed for GraphQL. Find better way.
func (t *Trapezza) Groups() []*GroupOrder {
	return t.groups
}

func (t *Trapezza) ChangeWaiter(waiter string) error {
	t.Waiter = waiter
	return nil
}

func (t *Trapezza) NewGroup(payer string) error {
	return t.groups.New(payer)
}

func (t *Trapezza) JoinGroup(client, payer string) error {
	c, cg, err := t.groups.OrderGroup(client)
	if err != nil {
		return err
	}

	_, pg, err := t.groups.OrderGroup(payer)
	if err != nil {
		return err
	}

	err = pg.Join(c)
	if err != nil {
		return err
	}

	if cg.Leave(c) {
		return t.groups.Remove(payer)
	}

	return nil
}

func (t *Trapezza) AddItems(client string, items []*Item) error {
	c, g, err := t.groups.OrderGroup(client)
	if err != nil {
		return err
	}

	gitems := make([]*OrderItem, len(items))
	for i, item := range items {
		gitems[i] = &OrderItem{Item: item}
	}

	return g.AddItems(c, gitems)
}

func (t *Trapezza) RemoveItem(client string, item string) error {
	c, g, err := t.groups.OrderGroup(client)
	if err != nil {
		return err
	}

	return g.RemoveItem(c, item)
}

func (t *Trapezza) SplitItem(who, with, item string) error {
	whoC, whoG, err := t.groups.OrderGroup(who)
	if err != nil {
		return err
	}

	gitem, err := whoG.Item(whoC, item)
	if err != nil {
		return err
	}

	withC, withG, err := t.groups.OrderGroup(with)
	if err != nil {
		return err
	}

	return withG.AddItem(withC, gitem)
}

func (t *Trapezza) ChangePayer(payer string) error {
	p, g, err := t.groups.OrderGroup(payer)
	if err != nil {
		return err
	}

	return g.ChangePayer(p)
}

func (t *Trapezza) CheckoutPayer(payer string) error {
	c, g, err := t.groups.OrderGroup(payer)
	if err != nil {
		return err
	}

	return g.CheckoutPayer(c)
}

func (t *Trapezza) CheckoutClient(client string) error {
	c, g, err := t.groups.OrderGroup(client)
	if err != nil {
		return err
	}

	return g.CheckoutClient(c)
}

func (t *Trapezza) WaiterCall(client, message string) error {
	call := time.Now()

	if t.LastCall.Add(WaiterCallTimeout).After(call) {
		return ErrCallTimeout
	}

	c, g, err := t.groups.OrderGroup(client)
	if err != nil {
		return err
	}

	g.AddCall(c, &Call{
		Time:    call,
		Message: message,
	})

	t.LastCall = call
	return nil
}

type Groups []*GroupOrder

func (gs *Groups) OrderGroup(client string) (*ClientOrder, *GroupOrder, error) {
	for _, g := range *gs {
		c, err := g.order(client)
		if err == nil {
			return c, g, nil
		}
	}

	return nil, nil, ErrWrongClient
}

func (gs *Groups) New(payer string) error {
	if len(*gs) >= GroupsLimit {
		return ErrGroupsLimit
	}

	_, _, err := gs.OrderGroup(payer)
	if err == nil {
		return ErrJoinedPayer
	}

	*gs = append(
		*gs,
		&GroupOrder{
			Payer: payer,
			Orders: []*ClientOrder{
				{
					Client: payer,
				},
			},
		},
	)

	return nil
}

func (gs *Groups) Remove(payer string) error {
	for i, g := range *gs {
		if g.Payer == payer {
			*gs = append((*gs)[:i], (*gs)[i+1:]...)
			return nil
		}
	}

	return ErrWrongClient
}
