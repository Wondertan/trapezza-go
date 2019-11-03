package trapezza

import (
	"errors"

	"github.com/Wondertan/trapezza-go/session"
	"github.com/Wondertan/trapezza-go/types"
)

var ErrWrongEvent = errors.New("session: wrong event")

type EventType = session.EventType

const (
	ChangeWaiter   EventType = "CHANGE_WAITER"
	ChangePayer    EventType = "CHANGE_PAYER"
	NewGroup       EventType = "NEW_GROUP"
	AddItems       EventType = "ADD_ITEMS"
	RemoveItem     EventType = "REMOVE_ITEM"
	SplitItem      EventType = "SPLIT_ITEM"
	CheckoutClient EventType = "CHECKOUT_CLIENT"
	CheckoutPayer  EventType = "CHECKOUT_PAYER"
	WaiterCall     EventType = "WAITER_CALL"
	JoinGroup      EventType = "JOIN_GROUP"
)

type Event interface {
	session.Event

	Trapezza() string // TODO Do we need to return Trapezza id here?

	setID(string)
}

type ChangeWaiterEvent struct {
	Waiter string

	session string
}

func (e *ChangeWaiterEvent) Type() EventType {
	return ChangeWaiter
}

func (e *ChangeWaiterEvent) Trapezza() string {
	return e.session
}

func (e *ChangeWaiterEvent) setID(id string) {
	e.session = id
}

type NewGroupEvent struct {
	Payer string

	session string
}

func (e *NewGroupEvent) Type() EventType {
	return NewGroup
}

func (e *NewGroupEvent) Trapezza() string {
	return e.session
}

func (e *NewGroupEvent) setID(id string) {
	e.session = id
}

type AddItemsEvent struct {
	Client string
	Items  []*types.Item

	session string
}

func (e *AddItemsEvent) Type() EventType {
	return AddItems
}

func (e *AddItemsEvent) Trapezza() string {
	return e.session
}

func (e *AddItemsEvent) setID(id string) {
	e.session = id
}

type RemoveItemEvent struct {
	Client string
	Item   string

	session string
}

func (e *RemoveItemEvent) Type() EventType {
	return RemoveItem
}

func (e *RemoveItemEvent) Trapezza() string {
	return e.session
}

func (e *RemoveItemEvent) setID(id string) {
	e.session = id
}

type SplitItemEvent struct {
	Who  string
	With string
	Item string

	session string
}

func (e *SplitItemEvent) Type() EventType {
	return SplitItem
}

func (e *SplitItemEvent) Trapezza() string {
	return e.session
}

func (e *SplitItemEvent) setID(id string) {
	e.session = id
}

type ChangePayerEvent struct {
	Payer string

	session string
}

func (e *ChangePayerEvent) Type() EventType {
	return ChangePayer
}

func (e *ChangePayerEvent) Trapezza() string {
	return e.session
}

func (e *ChangePayerEvent) setID(id string) {
	e.session = id
}

type JoinGroupEvent struct {
	Payer  string
	Client string

	session string
}

func (e *JoinGroupEvent) Type() EventType {
	return JoinGroup
}

func (e *JoinGroupEvent) Trapezza() string {
	return e.session
}

func (e *JoinGroupEvent) setID(id string) {
	e.session = id
}

type CheckoutClientEvent struct {
	Client string

	session string
}

func (e *CheckoutClientEvent) Type() EventType {
	return CheckoutClient
}

func (e *CheckoutClientEvent) Trapezza() string {
	return e.session
}

func (e *CheckoutClientEvent) setID(id string) {
	e.session = id
}

type CheckoutPayerEvent struct {
	Payer string

	session string
}

func (e *CheckoutPayerEvent) Type() EventType {
	return CheckoutPayer
}

func (e *CheckoutPayerEvent) Trapezza() string {
	return e.session
}

func (e *CheckoutPayerEvent) setID(id string) {
	e.session = id
}

type WaiterCallEvent struct {
	Client  string
	Message string

	session string
}

func (e *WaiterCallEvent) Type() EventType {
	return WaiterCall
}

func (e *WaiterCallEvent) Trapezza() string {
	return e.session
}

func (e *WaiterCallEvent) setID(id string) {
	e.session = id
}
