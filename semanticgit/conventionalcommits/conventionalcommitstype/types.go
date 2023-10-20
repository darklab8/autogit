package conventionalcommitstype

type Footer struct {
	Token   string
	Content string
}

type ParsedCommit struct {
	Type        string
	Exclamation bool
	Scope       string
	Subject     string
	Body        string
	Footers     []Footer
	Hash        string
	Issue       []string
}
