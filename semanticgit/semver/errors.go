package semver

type NotParsed struct{}

func (m NotParsed) Error() string {
	return "not parsed at all"
}
