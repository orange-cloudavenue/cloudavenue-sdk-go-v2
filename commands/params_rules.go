/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commands

// // ConditionExpr est une interface pour tout noeud logique
// type ConditionExpr interface {
// 	Eval(params map[string]interface{}) bool
// }

// // Condition est une feuille (test sur un champ)
// type Condition struct {
// 	Field string
// 	Value interface{}
// }

// func (c Condition) Eval(params map[string]interface{}) bool {
// 	val, ok := params[c.Field]
// 	return ok && val == c.Value
// }

// // AndExpr est un noeud ET logique
// type AndExpr struct {
// 	Exprs []ConditionExpr
// }

// func (a AndExpr) Eval(params map[string]interface{}) bool {
// 	for _, expr := range a.Exprs {
// 		if !expr.Eval(params) {
// 			return false
// 		}
// 	}
// 	return true
// }

// // OrExpr est un noeud OU logique
// type OrExpr struct {
// 	Exprs []ConditionExpr
// }

// func (o OrExpr) Eval(params map[string]interface{}) bool {
// 	for _, expr := range o.Exprs {
// 		if expr.Eval(params) {
// 			return true
// 		}
// 	}
// 	return false
// }

// // RuleValues defines the constraints for a field.
// type RuleValues struct {
// 	Editable    bool
// 	Min         *int
// 	Max         *int
// 	Equal       *int
// 	Enum        []interface{}
// 	Pattern     string
// 	Description string
// }

// type ParamsRules []ConditionalRule

// // ConditionalRule represents a rule that applies only if all conditions are met.
// // Now supports Console(s) restriction.
// type ConditionalRule struct {
// 	When     ConditionExpr
// 	Target   string
// 	Rule     RuleValues
// 	Consoles []consoles.Console // If empty, applies to all consoles. If not, applies only to listed consoles.
// }

// NewRules creates a new rules.
func NewRules(rules []ConditionalRule) ParamsRules {
	return ParamsRules(rules)
}

// // MarkdownDoc generates a markdown documentation for all rules of a resource.
// func (r ParamsRules) MarkdownDoc(resource string) string {
// 	// Group rules by unique combinations of When conditions and Consoles
// 	grouped := map[string][]ConditionalRule{}
// 	for _, rule := range r {
// 		conds := []string{}
// 		for _, c := range rule.When {
// 			conds = append(conds, fmt.Sprintf("%s=%v", c.Field, c.Value))
// 		}
// 		if len(rule.Consoles) > 0 {
// 			conds = append(conds, fmt.Sprintf("Consoles=%v", rule.Consoles))
// 		}
// 		grouped[strings.Join(conds, ", ")] = append(grouped[strings.Join(conds, ", ")], rule)
// 	}
// 	var sb strings.Builder
// 	sb.WriteString(fmt.Sprintf("## %s Rules\n\n", strings.Title(resource)))
// 	for conds, group := range grouped {
// 		if conds != "" {
// 			sb.WriteString(fmt.Sprintf("### When %s\n\n", conds))
// 		}
// 		sb.WriteString("| Field | Editable | Min | Max | Equal | Enum | Pattern | Description |\n")
// 		sb.WriteString("|-------|----------|-----|-----|-------|------|---------|-------------|\n")
// 		for _, rule := range group {
// 			sb.WriteString(fmt.Sprintf("| %s | %v | %v | %v | %v | %v | %s | %s |\n",
// 				rule.Target,
// 				rule.Rule.Editable,
// 				valOrEmpty(rule.Rule.Min),
// 				valOrEmpty(rule.Rule.Max),
// 				valOrEmpty(rule.Rule.Equal),
// 				enumToString(rule.Rule.Enum),
// 				rule.Rule.Pattern,
// 				rule.Rule.Description,
// 			))
// 		}
// 		sb.WriteString("\n")
// 	}
// 	return sb.String()
// }

// func valOrEmpty(v *int) string {
// 	if v == nil {
// 		return ""
// 	}
// 	return fmt.Sprintf("%v", *v)
// }

// func enumToString(e []interface{}) string {
// 	if len(e) == 0 {
// 		return ""
// 	}
// 	parts := make([]string, len(e))
// 	for i, v := range e {
// 		parts[i] = fmt.Sprintf("%v", v)
// 	}
// 	return strings.Join(parts, ", ")
// }
