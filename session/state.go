package session

type State struct {
	Id     ID
	Table  string
	Waiter string
	Orders []*Order
}

type Order struct {
	Client string
	Items  []string
}