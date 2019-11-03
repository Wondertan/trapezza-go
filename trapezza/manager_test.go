package trapezza

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_Session(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	man := NewManager(ctx)

	rest := "test"
	table := "test"

	id, err := man.NewSession(rest, table)
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, len(id), IDLength, "id length is not right")

	ses, err := man.Session(rest, table)
	require.Nil(t, err, "unexpected error")
	assert.NotNil(t, ses)

	err = man.EndSession(rest, table)
	require.Nil(t, err, "unexpected error")

	ses, err = man.SessionByID(id)
	assert.Equal(t, err, ErrNotFound)
	assert.Nil(t, ses)
}
