package semver

type NotParsedSemver struct{}

func (m NotParsedSemver) Error() string {
	return "not parsed semver at all"
}
