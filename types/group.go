package types

import (
	"time"
)

type GroupOrder struct {
	Payer  string
	Orders []*ClientOrder

	Total float64
}

func (g *GroupOrder) AddItems(order *ClientOrder, items []*OrderItem) error {
	for _, item := range items {
		err := g.AddItem(order, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GroupOrder) AddItem(order *ClientOrder, item *OrderItem) error {
	if len(order.Items) >= ItemsLimit {
		return ErrItemsLimit
	}
	if order.CheckedOut {
		return ErrCheckedOut
	}

	item.addGroup(g)
	order.addItem(item)
	return nil
}

func (g *GroupOrder) RemoveItem(order *ClientOrder, itemId string) error {
	item, err := order.item(itemId)
	if err != nil {
		return err
	}

	item.removeGroup(g)
	return order.removeItem(itemId)
}

func (g *GroupOrder) Item(order *ClientOrder, itemId string) (*OrderItem, error) {
	return order.item(itemId)
}

func (g *GroupOrder) Join(order *ClientOrder) error {
	if len(g.Orders) >= ClientsLimit {
		return ErrClientsLimit
	}

	g.Orders = append(g.Orders, order)
	return nil
}

func (g *GroupOrder) Leave(order *ClientOrder) bool {
	for i, c := range g.Orders {
		if c.Client == order.Client {
			g.Orders = append(g.Orders[:i], g.Orders[i+1:]...)
		}
	}

	return len(g.Orders) == 0
}

// TODO Persist after checkout
func (g *GroupOrder) CheckoutClient(order *ClientOrder) error {
	if g.Payer == order.Client {
		return ErrPayerCheckout
	}

	order.CheckedOut = true
	return nil
}

func (g *GroupOrder) CheckoutPayer(order *ClientOrder) error {
	order.CheckedOut = true

	for _, c := range g.Orders {
		if !c.CheckedOut {
			return g.ChangePayer(c) // TODO Handle that previous payer still owns his items
		}
	}

	return nil
}

func (g *GroupOrder) AddCall(order *ClientOrder, time time.Time, message string) {
	order.addCall(&Call{
		Time:    time,
		Message: message,
	})
}

func (g *GroupOrder) AnswerCall(order *ClientOrder, waiter string) error {
	call := order.lastCall()
	if call.Waiter != "" {
		return ErrAnswered
	}

	call.Waiter = waiter
	return nil
}

func (g *GroupOrder) ChangePayer(order *ClientOrder) error {
	if g.Payer == order.Client {
		return ErrAlreadyPayer
	}

	g.Payer = order.Client
	return nil
}

func (g *GroupOrder) order(id string) (*ClientOrder, error) {
	for _, c := range g.Orders {
		if c.Client == id {
			if c.CheckedOut {
				return nil, ErrCheckedOut
			}

			return c, nil
		}
	}

	return nil, ErrWrongClient
}

type Call struct {
	Time    time.Time
	Message string
	Waiter  string
}

type ClientOrder struct {
	Client     string
	Items      []*OrderItem
	Calls      []*Call
	CheckedOut bool
}

func (c *ClientOrder) addItem(item *OrderItem) {
	c.Items = append(c.Items, item)
}

func (c *ClientOrder) removeItem(id string) error {
	for i, item := range c.Items {
		if item.Id == id {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return nil
		}
	}

	return ErrWrongItem
}

func (c *ClientOrder) item(id string) (*OrderItem, error) {
	for _, item := range c.Items {
		if item.Id == id {
			return item, nil
		}
	}

	return nil, ErrWrongItem
}

func (c *ClientOrder) addCall(call *Call) {
	c.Calls = append(c.Calls, call)
}

func (c *ClientOrder) lastCall() *Call {
	return c.Calls[len(c.Calls)-1]
}

type OrderItem struct {
	*Item

	Groups []*GroupOrder
	split  float64
}

func (i *OrderItem) addGroup(g *GroupOrder) {
	i.clear()
	i.Groups = append(i.Groups, g)
	i.calc()
}

func (i *OrderItem) removeGroup(rm *GroupOrder) {
	i.clear()
	for j, g := range i.Groups {
		if g.Payer == rm.Payer {
			i.Groups = append(i.Groups[:j], i.Groups[j+1:]...)
		}
	}
	i.calc()
}

func (i *OrderItem) clear() {
	// remove prices from all groups
	for _, g := range i.Groups {
		g.Total -= i.split
	}
}

func (i *OrderItem) calc() {
	// calculate new price
	i.split = i.Price / float64(len(i.Groups))

	// apply new price
	for _, g := range i.Groups {
		g.Total += i.split
	}
}
