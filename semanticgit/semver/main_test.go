package semver

import (
	"autogit/settings/testutils"
	_ "autogit/settings/testutils/autouse"
	"testing"
)

func TestSemver(t *testing.T) {
	ver1, _ := Parse("v0.0.1")
	testutils.Equal(t, 0, ver1.Major)
	testutils.Equal(t, 0, ver1.Minor)
	testutils.Equal(t, 1, ver1.Patch)
}

func TestSemver2(t *testing.T) {
	ver1, _ := Parse("v1.0.1")
	testutils.Equal(t, 1, ver1.Major)
	testutils.Equal(t, 0, ver1.Minor)
	testutils.Equal(t, 1, ver1.Patch)
}

func TestSemver3(t *testing.T) {
	ver1, _ := Parse("v1.0.1-a.1")
	testutils.Equal(t, 1, ver1.Major)
	testutils.Equal(t, 0, ver1.Minor)
	testutils.Equal(t, 1, ver1.Patch)
	testutils.Equal(t, "a.1", ver1.Prerelease)
}

func TestSemver4(t *testing.T) {
	ver1, _ := Parse("v1.0.1-a.1-b.2+324234")
	testutils.Equal(t, 1, ver1.Major)
	testutils.Equal(t, 0, ver1.Minor)
	testutils.Equal(t, 1, ver1.Patch)
	testutils.Equal(t, "a.1-b.2", ver1.Prerelease)
	testutils.Equal(t, "324234", ver1.Build)
}
