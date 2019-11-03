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

	ses, err := man.NewSession()
	require.Nil(t, err, "unexpected error")
	assert.Equal(t, len(ses.ID()), IDLength, "id length is not right")

	err = man.EndSession(ses.ID())
	require.Nil(t, err, "unexpected error")

	ses, err = man.Session(ses.ID())
	assert.Equal(t, err, ErrNotFound)
	assert.Nil(t, ses)
}
