package restaurant

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Wondertan/trapezza-go/trapezza"
)

func TestManager(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rest := "test"
	table := "test"

	man := NewManager(ctx, &fakeSessionManager{})

	ses, err := man.NewTrapezzaSession(rest, table)
	require.Nil(t, err, "unexpected error")

	err = man.EndTrapezzaSession(rest, table)
	require.Nil(t, err, "unexpected error")

	ses, err = man.TrapezzaSession(rest, table)
	assert.Equal(t, err, trapezza.ErrNotFound)
	assert.Nil(t, ses)

	t.Run("Subscriptions", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		sub, err := man.SubscribeEvents(ctx, rest)
		require.Nil(t, err, "unexpected error")

		_, err = man.NewTrapezzaSession(rest, table)
		require.Nil(t, err, "unexpected error")

		time.Sleep(time.Millisecond * 10) // wait till published

		event := <-sub
		assert.Equal(t, event.Restaurant(), rest)

		cancel()
		time.Sleep(time.Millisecond * 10) // wait till closed

		err = man.EndTrapezzaSession(rest, table)
		require.Nil(t, err, "unexpected error")

		event, ok := <-sub
		assert.False(t, ok)
		assert.Nil(t, event)
	})
}

type fakeSessionManager struct{}

func (m *fakeSessionManager) NewSession() (*trapezza.Session, error) {
	return &trapezza.Session{}, nil
}

func (m *fakeSessionManager) EndSession(string) error {
	return nil
}

func (m *fakeSessionManager) Session(string) (*trapezza.Session, error) {
	return &trapezza.Session{}, nil
}
