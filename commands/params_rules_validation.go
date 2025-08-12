/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commands

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

// ConditionExpr is an interface for logical conditions (AND/OR/leaf)
type ConditionExpr interface {
	Eval(val reflect.Value) bool
}

// Condition is a leaf node: field == value
type Condition struct {
	Field string
	Value interface{}
}

func (c Condition) Eval(val reflect.Value) bool {
	fieldVal := getFieldByParamSpecName(val, c.Field)
	if !fieldVal.IsValid() {
		return false
	}
	return reflect.DeepEqual(fieldVal.Interface(), c.Value)
}

// AndExpr is a logical AND node
type AndExpr struct {
	Exprs []ConditionExpr
}

func (a AndExpr) Eval(val reflect.Value) bool {
	for _, expr := range a.Exprs {
		if !expr.Eval(val) {
			return false
		}
	}
	return true
}

// OrExpr is a logical OR node
type OrExpr struct {
	Exprs []ConditionExpr
}

func (o OrExpr) Eval(val reflect.Value) bool {
	for _, expr := range o.Exprs {
		if expr.Eval(val) {
			return true
		}
	}
	return false
}

// RuleValues defines the constraints for a field.
type RuleValues struct {
	Editable    bool
	Min         *int
	Max         *int
	Equal       *int
	Unit        string
	Enum        []interface{}
	Pattern     string
	Description string
}

type ParamsRules []ConditionalRule

// ConditionalRule with ConditionExpr
type ConditionalRule struct {
	When     ConditionExpr
	Target   string
	Rule     RuleValues
	Consoles []consoles.ConsoleName
}

// Validate applies all rules and returns error on first validation failure.
func (rules ParamsRules) validate(client cav.Client, params interface{}) error {
	val := reflect.ValueOf(params)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return errors.New("params must be a struct or pointer to struct")
	}

	for _, rule := range rules {
		if len(rule.Consoles) > 0 {
			// Check if the rule applies to the current console
			if !slices.ContainsFunc(rule.Consoles, func(c consoles.ConsoleName) bool {
				return c.GetSiteID() == client.GetConsole().GetSiteID()
			}) {
				continue // Skip this rule if it doesn't apply to the current console
			}
		}

		if rule.When != nil && !rule.When.Eval(val) {
			continue
		}

		values, err := GetAllValuesAtTarget(params, rule.Target)
		if err != nil {
			return fmt.Errorf("field %s not found in params: %w", rule.Target, err)
		}

		if len(values) == 0 {
			return fmt.Errorf("field %s not found in params", rule.Target)
		}
		for _, v := range values {
			vVal := reflect.ValueOf(v)
			if err := applyRuleValues(vVal, rule.Rule, rule.Target); err != nil {
				return err
			}
		}
	}
	return nil
}

// getFieldByParamSpecName: snake_case matching
func getFieldByParamSpecName(val reflect.Value, name string) reflect.Value {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if toSnakeCase(field.Name) == name {
			return val.Field(i)
		}
	}
	return reflect.Value{}
}

// toSnakeCase convertit CamelCase en snake_case
func toSnakeCase(str string) string {
	var out []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out = append(out, '_')
		}
		out = append(out, r)
	}
	return strings.ToLower(string(out))
}

func stringToLower(s string) string {
	res := []rune(s)
	for i, c := range res {
		if c >= 'A' && c <= 'Z' {
			res[i] = c + ('a' - 'A')
		}
	}
	return string(res)
}

// applyRuleValues applies RuleValues validation logic to fieldVal
// Now supports RuleValues.Enum values that can be regexp.Regexp or *regexp.Regexp
func applyRuleValues(fieldVal reflect.Value, rule RuleValues, fieldName string) error {
	if rule.Min != nil {
		if fieldVal.Kind() == reflect.Int && fieldVal.Int() < int64(*rule.Min) {
			return fmt.Errorf("%s must be >= %d", fieldName, *rule.Min)
		}
	}
	if rule.Max != nil {
		if fieldVal.Kind() == reflect.Int && fieldVal.Int() > int64(*rule.Max) {
			return fmt.Errorf("%s must be <= %d", fieldName, *rule.Max)
		}
	}
	if rule.Equal != nil {
		if fieldVal.Kind() == reflect.Int && fieldVal.Int() != int64(*rule.Equal) {
			return fmt.Errorf("%s must be == %d", fieldName, *rule.Equal)
		}
	}
	if len(rule.Enum) > 0 {
		found := false
		val := fieldVal.Interface()
		for _, e := range rule.Enum {
			switch enumVal := e.(type) {
			case regexp.Regexp:
				strVal, ok := val.(string)
				if ok && enumVal.MatchString(strVal) {
					found = true
					break
				}
			case *regexp.Regexp:
				strVal, ok := val.(string)
				if ok && enumVal != nil && enumVal.MatchString(strVal) {
					found = true
					break
				}
			default:
				if reflect.DeepEqual(val, e) {
					found = true
					break
				}
			}
		}
		if !found {
			return fmt.Errorf("%s must be one of %v", fieldName, rule.Enum)
		}
	}
	if rule.Pattern != "" && fieldVal.Kind() == reflect.String {
		matched, err := regexp.MatchString(rule.Pattern, fieldVal.String())
		if err != nil || !matched {
			return fmt.Errorf("%s must match pattern %s", fieldName, rule.Pattern)
		}
	}
	return nil
}
