package session

import (
	"fmt"
)

type EventType string

const (
	Waiter EventType = "WAITER"
	Client EventType = "CLIENT"
	Item   EventType = "ITEM"
	Table  EventType = "TABLE"
)

type Event interface {
	Handle(*State) error
}

type WaiterEvent struct {
	Session ID
	Type EventType
	Waiter string
}

func (e *WaiterEvent) Handle(s *State) error {
	s.Waiter = e.Waiter
	return nil
}

type ClientEvent struct {
	Session ID
	Type EventType
	Client string
}

func (e *ClientEvent) Handle(s *State) error {
	for _, order := range s.Orders {
		if order != nil && order.Client == e.Client {
			return fmt.Errorf("session: client already joined")
		}
	}

	s.Orders = append(s.Orders, &Order{
		Client: e.Client,
	})

	return nil
}

type ItemEvent struct {
	Session ID
	Type EventType
	Client string
	Item   string
}

func (e *ItemEvent) Handle(s *State) error {
	for _, order := range s.Orders {
		if order != nil && order.Client == e.Client {
			order.Items = append(order.Items, e.Item)
			return nil
		}
	}

	return fmt.Errorf("session: client is not joined")
}

type TableEvent struct {
	Session ID
	Type EventType
	Table string
}

func (e *TableEvent) Handle(s *State) error {
	s.Table = e.Table
	return nil
}
