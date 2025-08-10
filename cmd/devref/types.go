package main

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"

type (
	Functionality struct {
		Title                 string
		Documentation         string
		MarkdownDocumentation string

		Commands         map[string]Func          // E.g Get, List, Create, Update, Delete
		SubFunctionality map[string]Functionality //E.g FirewallRule
	}

	Func struct {
		Namespace string `json:"namespace"`
		Resource  string `json:"resource"`
		Verb      string `json:"verb"`

		// Documentation is the documentation of the command.
		ShortDocumentation string `json:"short_documentation"`
		LongDocumentation  string `json:"long_documentation"`

		// MarkdownDocumentation is the markdown documentation of the command. Used for top-level commands.
		MarkdownDocumentation string `json:"markdown_documentation"`

		Params []FuncParam `json:"params"`

		Deprecated        bool   `json:"deprecated"`
		DeprecatedMessage string `json:"deprecated_message"`

		Model []commands.DocModel `json:"model"`

		Rules []RuleExport `json:"rules,omitempty"`
	}

	FuncParam struct {
		Name                  string `json:"name"`
		Description           string `json:"description"`
		Required              bool   `json:"required"`
		Example               string `json:"example"`
		Type                  string `json:"type"`
		ValidatorsDescription string `json:"validators_description"`
	}
)
