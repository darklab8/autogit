package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Equal[T interface{}](t *testing.T, tag_name T, actual T) {
	assert.Equal(t, tag_name, actual)
}
