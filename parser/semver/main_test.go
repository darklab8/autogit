package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemver(t *testing.T) {
	ver1 := Parse("v0.0.1")
	assert.Equal(t, 0, ver1.Major)
	assert.Equal(t, 0, ver1.Minor)
	assert.Equal(t, 1, ver1.Patch)
}

func TestSemver2(t *testing.T) {
	ver1 := Parse("v1.0.1")
	assert.Equal(t, 1, ver1.Major)
	assert.Equal(t, 0, ver1.Minor)
	assert.Equal(t, 1, ver1.Patch)
}

func TestSemver3(t *testing.T) {
	ver1 := Parse("v1.0.1-a.1")
	assert.Equal(t, 1, ver1.Major)
	assert.Equal(t, 0, ver1.Minor)
	assert.Equal(t, 1, ver1.Patch)
	assert.Equal(t, "a.1", ver1.Prerelease)
}

func TestSemver4(t *testing.T) {
	ver1 := Parse("v1.0.1-a.1-b.2+324234")
	assert.Equal(t, 1, ver1.Major)
	assert.Equal(t, 0, ver1.Minor)
	assert.Equal(t, 1, ver1.Patch)
	assert.Equal(t, "a.1-b.2", ver1.Prerelease)
	assert.Equal(t, "324234", ver1.Build)
}
