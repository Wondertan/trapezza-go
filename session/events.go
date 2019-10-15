package session

import (
	"fmt"
)

type Event interface {
	Handle(*State) error
}

type waiterEvent struct {
	waiter string
}

func (e *waiterEvent) Handle(s *State) error {
	s.Waiter = e.waiter
	return nil
}

type clientEvent struct {
	client string
}

func (e *clientEvent) Handle(s *State) error {
	if s.Orders == nil {
		s.Orders = make(map[string][]string)
	}

	s.Orders[e.client] = []string{}
	return nil
}

type itemEvent struct {
	client string
	item   string
}

func (e *itemEvent) Handle(s *State) error {
	items, ok := s.Orders[e.client]
	if !ok {
		return fmt.Errorf("client %s is not in the session", e.client)
	}

	s.Orders[e.client] = append(items, e.item)
	return nil
}

type tableEvent struct {
	table string
}

func (e *tableEvent) Handle(s *State) error {
	s.Table = e.table
	return nil
}
