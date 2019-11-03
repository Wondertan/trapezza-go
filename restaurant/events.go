package restaurant

type EventType string

const (
	NewSession EventType = "NEW_SESSION"
	EndSession EventType = "END_SESSION"
)

type Event interface {
	Type() EventType
	Restaurant() string
}

type NewSessionEvent struct {
	ID    string
	Table string

	restaurant string
}

func (e *NewSessionEvent) Type() EventType {
	return NewSession
}

func (e *NewSessionEvent) Restaurant() string {
	return e.restaurant
}

type EndSessionEvent struct {
	Table string

	restaurant string
}

func (e *EndSessionEvent) Type() EventType {
	return EndSession
}

func (e *EndSessionEvent) Restaurant() string {
	return e.restaurant
}
