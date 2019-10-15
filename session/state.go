package session

type State struct {
	Table  string
	Waiter string
	Orders map[string][]string
}
