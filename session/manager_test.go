package session

import (
	"context"
	"testing"
)

func TestManager_Session(t *testing.T) {
	ctx := context.Background()
	man := NewManager(ctx)

	id := man.NewSession()
	if len(id) != IDLength {
		t.Fatal("id length is not right")
	}

	ses, err := man.Session(id)
	if err != nil {
		t.Fatal(err)
	}
	if ses == nil {
		t.Fatal("nil session")
	}

	err = man.EndSession(id)
	if err != nil {
		t.Fatal(err)
	}

	ses, err = man.Session(id)
	if err != ErrNotFound {
		t.Fatal("wrong error")
	}
	if ses != nil {
		t.Fatal("should be nil")
	}
}
