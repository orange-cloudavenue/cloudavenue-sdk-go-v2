/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

const (
	DynamicSecurityGroupCriteriaRuleTypeVMName = "VM_NAME"
	DynamicSecurityGroupCriteriaRuleTypeVMTag  = "VM_TAG"

	DynamicSecurityGroupCriteriaRuleOperatorEquals     = "EQUALS"
	DynamicSecurityGroupCriteriaRuleOperatorContains   = "CONTAINS"
	DynamicSecurityGroupCriteriaRuleOperatorStartsWith = "STARTS_WITH"
	DynamicSecurityGroupCriteriaRuleOperatorEndsWith   = "ENDS_WITH"

	// DynamicSecurityGroupMaxCriteria is the maximum number of criteria allowed on a Dynamic
	// Security Group. Each criteria is combined with the others using a logical OR.
	DynamicSecurityGroupMaxCriteria = 3

	// DynamicSecurityGroupMaxRulesPerCriteria is the maximum number of rules allowed within a
	// single criteria. Rules within a criteria are combined using a logical AND.
	DynamicSecurityGroupMaxRulesPerCriteria = 4
)

// DynamicSecurityGroupAllowedOperators maps a rule type to the operators allowed for that rule
// type. VM_NAME only supports STARTS_WITH and CONTAINS. VM_TAG supports EQUALS, CONTAINS,
// STARTS_WITH, and ENDS_WITH.
var DynamicSecurityGroupAllowedOperators = map[string][]string{
	DynamicSecurityGroupCriteriaRuleTypeVMName: {
		DynamicSecurityGroupCriteriaRuleOperatorStartsWith,
		DynamicSecurityGroupCriteriaRuleOperatorContains,
	},
	DynamicSecurityGroupCriteriaRuleTypeVMTag: {
		DynamicSecurityGroupCriteriaRuleOperatorEquals,
		DynamicSecurityGroupCriteriaRuleOperatorContains,
		DynamicSecurityGroupCriteriaRuleOperatorStartsWith,
		DynamicSecurityGroupCriteriaRuleOperatorEndsWith,
	},
}

type (
	ParamsListDynamicSecurityGroup struct {
		// VDCGroupID is the ID of the Vdc Group owning the Dynamic Security Groups to list.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning the Dynamic Security Groups to list.
		VDCGroupName string
	}

	ParamsGetDynamicSecurityGroup struct {
		// ID is the unique identifier of the Dynamic Security Group.
		ID string

		// Name is the name of the Dynamic Security Group.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this Dynamic Security Group.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Dynamic Security Group.
		VDCGroupName string
	}

	ParamsCreateDynamicSecurityGroup struct {
		// Name is the name of the Dynamic Security Group.
		Name string

		// Description is the description of the Dynamic Security Group.
		Description string

		// VDCGroupID is the ID of the Vdc Group owning this Dynamic Security Group.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Dynamic Security Group.
		VDCGroupName string

		// Criteria is the list of criteria (max 3) that determine VM membership. A VM belongs to
		// the group if it matches at least one criteria (logical OR between criteria).
		Criteria []ParamsDynamicSecurityGroupCriteria
	}

	ParamsDynamicSecurityGroupCriteria struct {
		// Rules is the list of rules (max 4) within this criteria. A VM must match all rules to
		// satisfy this criteria (logical AND between rules).
		Rules []ParamsDynamicSecurityGroupCriteriaRule
	}

	ParamsDynamicSecurityGroupCriteriaRule struct {
		// RuleType is the type of the rule (VM_NAME or VM_TAG).
		RuleType string

		// Operator is the operator of the rule. Allowed values depend on RuleType: VM_NAME allows
		// STARTS_WITH and CONTAINS; VM_TAG allows EQUALS, CONTAINS, STARTS_WITH, and ENDS_WITH.
		Operator string

		// Value is the value to match against.
		Value string
	}

	ParamsUpdateDynamicSecurityGroup struct {
		// ID is the unique identifier of the Dynamic Security Group to update.
		ID string

		// Name is the new name of the Dynamic Security Group.
		Name string

		// Description is the new description of the Dynamic Security Group.
		Description string

		// VDCGroupID is the ID of the Vdc Group owning this Dynamic Security Group.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Dynamic Security Group.
		VDCGroupName string

		// Criteria is the list of criteria (max 3) that determine VM membership.
		Criteria []ParamsDynamicSecurityGroupCriteria
	}

	ParamsDeleteDynamicSecurityGroup struct {
		// ID is the unique identifier of the Dynamic Security Group to delete.
		ID string

		// Name is the name of the Dynamic Security Group to delete.
		Name string

		// VDCGroupID is the ID of the Vdc Group owning this Dynamic Security Group.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning this Dynamic Security Group.
		VDCGroupName string
	}
)
