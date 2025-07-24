package commands

import (
	"strings"
	"testing"
)

func TestValidatorOneOf_GetKey(t *testing.T) {
	v := ValidatorOneOf("foo", "bar", "baz")
	expected := "oneof=foo bar baz"
	if v.GetKey() != expected {
		t.Errorf("GetKey() = %v, want %v", v.GetKey(), expected)
	}
}

func TestValidatorOneOf_GetDescription(t *testing.T) {
	v := ValidatorOneOf("foo", "bar")
	expected := "Validates that the value is one of: foo, bar"
	if v.GetDescription() != expected {
		t.Errorf("GetDescription() = %v, want %v", v.GetDescription(), expected)
	}
}

func TestValidatorOneOf_GetMarkdownDescription(t *testing.T) {
	v := ValidatorOneOf("foo", "bar")
	expected := "Validates that the value is one of: `foo`, `bar`"
	if v.GetMarkdownDescription() != expected {
		t.Errorf("GetMarkdownDescription() = %v, want %v", v.GetMarkdownDescription(), expected)
	}
}

// Optionally, test with empty values
func TestValidatorOneOf_EmptyValues(t *testing.T) {
	v := ValidatorOneOf()
	if !strings.HasPrefix(v.GetKey(), "oneof=") {
		t.Errorf("GetKey() with empty values should start with 'oneof='")
	}
	if v.GetDescription() != "Validates that the value is one of: " {
		t.Errorf("GetDescription() with empty values failed")
	}
}
