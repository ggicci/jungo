package http

import (
	"testing"
)

func TestGenRequestID(t *testing.T) {
	t.Logf("generated: %s", GenRequestID())
}
