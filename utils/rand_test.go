package utils

import (
	"testing"
)

func TestRandString(t *testing.T) {
	l := 5
	res := RandString(l)

	if len(res) != l {
		t.Fatal("wrong length")
	}
}
