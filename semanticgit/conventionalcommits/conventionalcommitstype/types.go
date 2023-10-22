package conventionalcommitstype

type FooterToken string
type FooterContent string

type Footer struct {
	Token   FooterToken
	Content FooterContent
}

type Scope string
type Type string
type Subject string
type Body string
type Hash string
type Issue string

type ParsedCommit struct {
	Type        Type
	Exclamation bool
	Scope       Scope
	Subject     Subject
	Body        Body
	Footers     []Footer
	Hash        Hash
	Issue       []Issue
}
