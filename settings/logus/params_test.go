package logus

import "testing"

func TestSlogging(t *testing.T) {
	Debug("123")

	Debug("123", TestParam(456))
}
