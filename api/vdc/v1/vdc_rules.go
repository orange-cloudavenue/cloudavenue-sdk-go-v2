/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"regexp"

	"github.com/orange-cloudavenue/common-go/utils"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

var vdcRules = commands.NewRules([]commands.ConditionalRule{
	// * ----------- disponibility_class ----------- *
	{
		Consoles: []consoles.ConsoleName{
			consoles.Console1,
			consoles.Console2,
		},
		When: commands.And(
			commands.NewCondition("service_class", "ECO"),
		).Build(),
		Target: "disponibility_class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"ONE-ROOM", "DUAL-ROOM"},
			Description: "Disponibility class allowed for Service Class ECO",
		},
	},
	{
		Consoles: []consoles.ConsoleName{
			consoles.Console4,
			consoles.Console5,
		},
		Target: "disponibility_class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"ONE-ROOM"},
			Description: "Disponibility class allowed for Service Class ECO",
		},
	},
	{
		Consoles: []consoles.ConsoleName{
			consoles.Console1,
			consoles.Console2,
		},
		When: commands.Or(
			commands.NewCondition("service_class", "STD"),
			commands.NewCondition("service_class", "HP"),
			commands.NewCondition("service_class", "VOIP"),
		).Build(),
		Target: "disponibility_class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"ONE-ROOM", "DUAL-ROOM", "HA-DUAL-ROOM"},
			Description: "Disponibility class allowed for Service Class STD, HP, VOIP",
		},
	},

	// * ----------- billing_model ----------- *
	{
		When: commands.Or(
			commands.NewCondition("service_class", "ECO"),
			commands.NewCondition("service_class", "STD"),
		).Build(),
		Target: "billing_model",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"PAYG", "DRAAS", "RESERVED"},
			Description: "Billing model allowed for Service Class ECO, STD",
		},
	},
	{
		When: commands.Or(
			commands.NewCondition("service_class", "HP"),
		).Build(),
		Target: "billing_model",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"PAYG", "RESERVED"},
			Description: "Billing model allowed for Service Class HP",
		},
	},
	{
		When: commands.Or(
			commands.NewCondition("service_class", "VOIP"),
		).Build(),
		Target: "billing_model",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"RESERVED"},
			Description: "Billing model allowed for Service Class VOIP",
		},
	},

	// * ----------- storage_billing_model ----------- *
	{
		Target: "storage_billing_model",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"PAYG", "RESERVED"},
			Description: "Storage billing model allowed for Service Class ECO",
		},
	},

	// * ----------- vcpu -----------

	{
		When: commands.Or(
			commands.NewCondition("billing_model", "PAYG"),
			commands.NewCondition("billing_model", "DRAAS"),
		).Build(),
		Target: "vcpu",
		Rule: commands.RuleValues{
			Editable:    true,
			Min:         utils.ToPTR(5),
			Max:         utils.ToPTR(200),
			Description: "VCPU allowed for Service Class ECO with PAYG or DRAAS billing model",
		},
	},
	{
		When: commands.Or(
			commands.NewCondition("billing_model", "RESERVED"),
		).Build(),
		Target: "vcpu",
		Rule: commands.RuleValues{
			Editable:    true,
			Min:         utils.ToPTR(2),
			Max:         utils.ToPTR(1136),
			Description: "VCPU allowed for Service Class ECO with RESERVED billing model",
		},
	},

	// * ----------- memory ----------- *

	{
		Target: "memory",
		Rule: commands.RuleValues{
			Editable:    true,
			Unit:        "GB",
			Min:         utils.ToPTR(1),
			Max:         utils.ToPTR(5120),
			Description: "Memory allowed for Service Class ECO",
		},
	},

	// * ----------- storage_profiles ----------- *
	{
		When: commands.Or(
			commands.NewCondition("disponibility_class", "ONE-ROOM"),
		).Build(),
		Target: "storage_profiles.{index}.class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"silver", "gold", "platinum3k", "platinum7k", regexp.MustCompile("^(silver|gold|platinum[3|7]{1}k)_(ocb[0-9]{1,7})$")},
			Description: "Storage profile class allowed for Disponibility Class ONE-ROOM",
		},
	},
	{
		Consoles: []consoles.ConsoleName{
			consoles.Console1,
			consoles.Console2,
		},
		When: commands.Or(
			commands.NewCondition("disponibility_class", "DUAL-ROOM"),
		).Build(),
		Target: "storage_profiles.{index}.class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"silver_r1", "silver_r2", "gold_r1", "gold_r2", "platinum3k_r1", "platinum3k_r2", "platinum7k_r1", "platinum7k_r2", regexp.MustCompile("^(silver|gold|platinum[3|7]{1}k)_(ocb[0-9]{1,7})_(r1|r2)$")},
			Description: "Storage profile class allowed for Disponibility Class DUAL-ROOM",
		},
	},
	{
		Consoles: []consoles.ConsoleName{
			consoles.Console1,
			consoles.Console2,
		},
		When: commands.Or(
			commands.NewCondition("disponibility_class", "HA-DUAL-ROOM"),
		).Build(),
		Target: "storage_profiles.{index}.class",
		Rule: commands.RuleValues{
			Editable:    false,
			Enum:        []interface{}{"gold_hm", "platinum3k_hm", "platinum7k_hm", regexp.MustCompile("^(gold|platinum[3|7]{1}k)_(ocb[0-9]{1,7})_(hm)$")},
			Description: "Storage profile class allowed for Disponibility Class HA-DUAL-ROOM",
		},
	},
},
)
