package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSession(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	state := &fakeState{
		data: "test",
	}

	ses := NewSession(ctx, state)

	t.Run("Error", func(t *testing.T) {
		resp, err := ses.EmitEvent(&fakeEvent{err: fakeError})
		require.Nil(t, err, "unexpected error")

		err = <-resp
		assert.Equal(t, fakeError, err, "wrong error")
	})

	t.Run("Success", func(t *testing.T) {
		in := "success"

		resp, err := ses.EmitEvent(&fakeEvent{data: in})
		require.Nil(t, err, "unexpected error")

		err = <-resp
		require.Nil(t, err, "unexpected error")

		ch, err := ses.State()
		require.Nil(t, err, "unexpected error")

		out := <-ch
		assert.Equal(t, state, out, "states are not equal")
	})

	t.Run("Subscriptions", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)

		in := &fakeEvent{data: "subscriptions"}

		sub1, err := ses.SubscribeUpdates(ctx)
		require.Nil(t, err, "unexpected error")

		sub2, err := ses.SubscribeUpdates(context.Background())
		require.Nil(t, err, "unexpected error")

		resp, err := ses.EmitEvent(in)
		require.Nil(t, err, "unexpected error")

		err = <-resp
		require.Nil(t, err, "unexpected error")

		out := <-sub1
		assert.Equal(t, in, out.Event, "events should be equal")

		out = <-sub2
		assert.Equal(t, in, out.Event, "events should be equal")

		cancel()                          // close one
		time.Sleep(time.Millisecond * 10) // wait till cancel

		resp, err = ses.EmitEvent(in)
		require.Nil(t, err, "unexpected error")

		err = <-resp
		require.Nil(t, err, "unexpected error")

		out, ok := <-sub1
		assert.False(t, ok)
		assert.Nil(t, out, "should be nil")

		out = <-sub2
		assert.Equal(t, in, out.Event, "events should be equal")
	})
}

var fakeError = errors.New("fake")

type fakeState struct {
	data string
}

func (s *fakeState) ID() string {
	return "fake"
}

func (s *fakeState) Handle(e Event) error {
	event := e.(*fakeEvent)
	if event.err != nil {
		return event.err
	} else {
		s.data = event.data
	}

	return nil
}

type fakeEvent struct {
	data string
	err  error
}

func (e *fakeEvent) Type() EventType {
	return "fake"
}
