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

type NewTrapezzaSessionEvent struct {
	Trapezza string
	Table    string

	restaurant string
}

func (e *NewTrapezzaSessionEvent) Type() EventType {
	return NewSession
}

func (e *NewTrapezzaSessionEvent) Restaurant() string {
	return e.restaurant
}

type EndTrapezzaSessionEvent struct {
	Trapezza string
	Table    string

	restaurant string
}

func (e *EndTrapezzaSessionEvent) Type() EventType {
	return EndSession
}

func (e *EndTrapezzaSessionEvent) Restaurant() string {
	return e.restaurant
}
