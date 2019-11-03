package restaurant

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	man := NewManager(ctx, &fakeSessionManager{})

	t.Run("Subscriptions", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		rest := "test"
		table := "test"

		sub, err := man.SubscribeEvents(ctx, rest)
		require.Nil(t, err, "unexpected error")

		_, err = man.NewSession(rest, table)
		require.Nil(t, err, "unexpected error")

		time.Sleep(time.Millisecond * 10) // wait till published

		event := <-sub
		assert.Equal(t, event.Restaurant(), rest)

		cancel()
		time.Sleep(time.Millisecond * 10) // wait till closed

		err = man.EndSession(rest, table)
		require.Nil(t, err, "unexpected error")

		event, ok := <-sub
		assert.False(t, ok)
		assert.Nil(t, event)
	})
}

type fakeSessionManager struct{}

func (m *fakeSessionManager) NewSession(rest, table string) (string, error) {
	return "test", nil
}

func (m *fakeSessionManager) EndSession(rest, table string) error {
	return nil
}
