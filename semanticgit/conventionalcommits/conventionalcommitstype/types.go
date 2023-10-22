package conventionalcommitstype

type Footer struct {
	Token   string
	Content string
}

type Scope string

type ParsedCommit struct {
	Type        string
	Exclamation bool
	Scope       Scope
	Subject     string
	Body        string
	Footers     []Footer
	Hash        string
	Issue       []string
}
