/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package main

import (
	"fmt"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

// --- Structures d'export ---

type ConditionExport struct {
	Field string            `json:"field,omitempty"`
	Value interface{}       `json:"value,omitempty"`
	And   []ConditionExport `json:"and,omitempty"`
	Or    []ConditionExport `json:"or,omitempty"`
}

type RuleExport struct {
	Consoles    []string        `json:"consoles"`
	WhenHuman   string          `json:"whenHuman"`
	When        ConditionExport `json:"when"`
	Target      string          `json:"target"`
	Min         *int            `json:"min,omitempty"`
	Max         *int            `json:"max,omitempty"`
	Enum        []interface{}   `json:"enum,omitempty"`
	Pattern     string          `json:"pattern,omitempty"`
	Description string          `json:"description,omitempty"`
	Unit        string          `json:"unit,omitempty"`
}

// --- Helpers pour exporter les conditions ---

func ExportCondition(expr commands.ConditionExpr) ConditionExport {
	switch v := expr.(type) {
	case commands.Condition:
		return ConditionExport{Field: v.Field, Value: v.Value}
	case commands.AndExpr:
		children := make([]ConditionExport, len(v.Exprs))
		for i, e := range v.Exprs {
			children[i] = ExportCondition(e)
		}
		return ConditionExport{And: children}
	case commands.OrExpr:
		children := make([]ConditionExport, len(v.Exprs))
		for i, e := range v.Exprs {
			children[i] = ExportCondition(e)
		}
		return ConditionExport{Or: children}
	default:
		return ConditionExport{}
	}
}

func ConditionToString(expr commands.ConditionExpr) string {
	switch v := expr.(type) {
	case commands.Condition:
		return fmt.Sprintf("%s = %v", v.Field, v.Value)
	case commands.AndExpr:
		parts := make([]string, len(v.Exprs))
		for i, e := range v.Exprs {
			parts[i] = ConditionToString(e)
		}
		return "(" + strings.Join(parts, " AND ") + ")"
	case commands.OrExpr:
		parts := make([]string, len(v.Exprs))
		for i, e := range v.Exprs {
			parts[i] = ConditionToString(e)
		}
		return "(" + strings.Join(parts, " OR ") + ")"
	default:
		return ""
	}
}

// --- Helper pour extraire les noms ou ID de consoles ---

func getConsoleNames(cs []consoles.ConsoleName) []string {
	var out []string
	for _, c := range cs {
		out = append(out, c.GetSiteName())
	}
	return out
}
