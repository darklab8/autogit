package testutils

import (
	"autogit/settings/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func EqualTag(t *testing.T, tag_name string, actual types.TagName) {
	assert.Equal(t, types.TagName(tag_name), actual)
}
