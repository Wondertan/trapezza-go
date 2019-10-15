package session

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestSessionEvent(t *testing.T) {
	ctx := context.Background()
	id := ID("test")

	ses := newSession(ctx, id)
	defer ses.Stop()

	if ses.ID() != id {
		t.Fatal("ids are not equal")
	}

	in := &State{
		Waiter: "test",
		Orders: map[string][]string{
			"testclient": {
				"chicken",
				"potato",
			},
		},
	}

	t.Run("Error", func(t *testing.T) {
		resp, err := ses.Event(&testEvent{
			State: in,
			err:   testError,
		})
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		err = <-resp
		if err != testError {
			t.Fatal("wrong error")
		}

		ch, err := ses.State()
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		out := <-ch
		if reflect.DeepEqual(in, out) {
			t.Fatal("should not be equal")
		}
	})

	t.Run("Success", func(t *testing.T) {
		resp, err := ses.Event(&testEvent{State: in})
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		err = <-resp
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		ch, err := ses.State()
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		out := <-ch
		if !reflect.DeepEqual(in, out) {
			t.Fatal("states are not equal")
		}
	})

	t.Run("Sub", func(t *testing.T) {
		event := &testEvent{State: in}

		sub, err := ses.SubscribeEvents()
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		resp, err := ses.Event(event)
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		err = <-resp
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		e := <-sub
		if !reflect.DeepEqual(e, event) {
			t.Fatal("events should be equal")
		}
	})
}

var testError = errors.New("test")

type testEvent struct {
	*State
	err error
}

func (e *testEvent) Handle(s *State) error {
	if e.err != nil {
		return e.err
	}

	s.Waiter = e.Waiter
	s.Orders = e.Orders
	return nil
}
