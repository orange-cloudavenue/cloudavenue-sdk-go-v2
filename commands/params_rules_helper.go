/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commands

// ConditionExprBuilder permet de construire dynamiquement des arbres de conditions complexes
type ConditionExprBuilder struct {
	expr ConditionExpr
}

// NewCondition crée un builder pour une feuille (champ == valeur)
func NewCondition(field string, value interface{}) ConditionExprBuilder {
	return ConditionExprBuilder{expr: Condition{Field: field, Value: value}}
}

// And combine plusieurs builders avec un AND logique
func And(exprs ...ConditionExprBuilder) ConditionExprBuilder {
	cExprs := make([]ConditionExpr, len(exprs))
	for i, e := range exprs {
		cExprs[i] = e.expr
	}
	return ConditionExprBuilder{expr: AndExpr{Exprs: cExprs}}
}

// Or combine plusieurs builders avec un OR logique
func Or(exprs ...ConditionExprBuilder) ConditionExprBuilder {
	cExprs := make([]ConditionExpr, len(exprs))
	for i, e := range exprs {
		cExprs[i] = e.expr
	}
	return ConditionExprBuilder{expr: OrExpr{Exprs: cExprs}}
}

// Build retourne l'arbre ConditionExpr final
func (b ConditionExprBuilder) Build() ConditionExpr {
	return b.expr
}
