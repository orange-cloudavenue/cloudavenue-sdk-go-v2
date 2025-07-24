package commands

type Validator interface {
	GetKey() string
	GetDescription() string
	GetMarkdownDescription() string
}
