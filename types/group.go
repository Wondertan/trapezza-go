package types

import (
	"time"
)

type Group struct {
	Payer   string
	Clients []*Client

	Total float64
}

func (g *Group) AddItems(client *Client, items []*GroupItem) error {
	for _, item := range items {
		err := g.AddItem(client, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Group) AddItem(client *Client, item *GroupItem) error {
	if len(client.Items) >= ItemsLimit {
		return ErrItemsLimit
	}
	if client.CheckedOut {
		return ErrCheckedOut
	}

	item.addGroup(g)
	client.addItem(item)
	return nil
}

func (g *Group) RemoveItem(client *Client, itemId string) error {
	item, err := client.item(itemId)
	if err != nil {
		return err
	}

	item.removeGroup(g)
	return client.removeItem(itemId)
}

func (g *Group) Item(client *Client, itemId string) (*GroupItem, error) {
	return client.item(itemId)
}

func (g *Group) Join(client *Client) error {
	if len(g.Clients) >= ClientsLimit {
		return ErrClientsLimit
	}

	g.Clients = append(g.Clients, client)
	return nil
}

func (g *Group) Leave(client *Client) bool {
	for i, c := range g.Clients {
		if c.Id == client.Id {
			g.Clients = append(g.Clients[:i], g.Clients[i+1:]...)
		}
	}

	return len(g.Clients) == 0
}

// TODO Persist after checkout
func (g *Group) CheckoutClient(client *Client) error {
	if g.Payer == client.Id {
		return ErrPayerCheckout
	}

	client.CheckedOut = true
	return nil
}

func (g *Group) CheckoutPayer(client *Client) error {
	client.CheckedOut = true

	for _, c := range g.Clients {
		if !c.CheckedOut {
			return g.ChangePayer(c) // TODO Handle that previous payer still owns his items
		}
	}

	return nil
}

func (g *Group) AddCall(client *Client, call *Call) {
	client.addCall(call)
}

func (g *Group) ChangePayer(client *Client) error {
	if g.Payer == client.Id {
		return ErrAlreadyPayer
	}

	g.Payer = client.Id
	return nil
}

func (g *Group) client(id string) (*Client, error) {
	for _, c := range g.Clients {
		if c.Id == id {
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
}

type Client struct {
	Id         string
	Items      []*GroupItem
	Calls      []*Call
	CheckedOut bool
}

func (c *Client) addItem(item *GroupItem) {
	c.Items = append(c.Items, item)
}

func (c *Client) removeItem(id string) error {
	for i, item := range c.Items {
		if item.Id == id {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return nil
		}
	}

	return ErrWrongItem
}

func (c *Client) item(id string) (*GroupItem, error) {
	for _, item := range c.Items {
		if item.Id == id {
			return item, nil
		}
	}

	return nil, ErrWrongItem
}

func (c *Client) addCall(call *Call) {
	c.Calls = append(c.Calls, call)
}

type GroupItem struct {
	*Item

	Groups []*Group
	split  float64
}

func (i *GroupItem) addGroup(g *Group) {
	i.clear()
	i.Groups = append(i.Groups, g)
	i.calc()
}

func (i *GroupItem) removeGroup(rm *Group) {
	i.clear()
	for j, g := range i.Groups {
		if g.Payer == rm.Payer {
			i.Groups = append(i.Groups[:j], i.Groups[j+1:]...)
		}
	}
	i.calc()
}

func (i *GroupItem) clear() {
	// remove prices from all groups
	for _, g := range i.Groups {
		g.Total -= i.split
	}
}

func (i *GroupItem) calc() {
	// calculate new price
	i.split = i.Price / float64(len(i.Groups))

	// apply new price
	for _, g := range i.Groups {
		g.Total += i.split
	}
}
