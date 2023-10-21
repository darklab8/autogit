package settings

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test_data/autogit.partial.yml
var ConfigPartial string

func TestReadPartialConfig(t *testing.T) {
	result := configRead([]byte(ConfigPartial))

	assert.Equal(t, 7, result.Validation.Rules.Header.Subject.MinWords)
	assert.Equal(t, true, result.Validation.Rules.Header.Scope.Lowercase)
}
