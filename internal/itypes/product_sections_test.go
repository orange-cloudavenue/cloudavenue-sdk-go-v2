/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// richProductSectionList uses the elements of the OVF ProductSection schema (DSP8023)
// that the public model does not carry.
const richProductSectionList = `<ProductSectionList xmlns="http://www.vmware.com/vcloud/v1.5" xmlns:ovf="http://schemas.dmtf.org/ovf/envelope/1" href="https://console.example.org/api/vApp/vapp-x/productSections/" type="application/vnd.vmware.vcloud.productSections+xml">
  <ovf:ProductSection ovf:required="false">
    <ovf:Info>Custom properties</ovf:Info>
    <ovf:Product>My product</ovf:Product>
    <ovf:Category>Network</ovf:Category>
    <ovf:Property ovf:key="keep" ovf:type="string" ovf:userConfigurable="true" ovf:value="" ovf:password="true" ovf:qualifiers="MaxLen(10)">
      <ovf:Label>keep</ovf:Label>
      <ovf:Description>kept description</ovf:Description>
      <ovf:Value ovf:value="old"/>
      <ovf:Value ovf:value="alt" ovf:configuration="big"/>
    </ovf:Property>
    <ovf:Property ovf:key="drop" ovf:type="string" ovf:userConfigurable="true" ovf:value="d"/>
    <ovf:Category>Other</ovf:Category>
    <vmw:Extension xmlns:vmw="http://example.org/vmw">ext</vmw:Extension>
  </ovf:ProductSection>
  <ovf:ProductSection ovf:class="vami" ovf:instance="x">
    <ovf:Info>vami</ovf:Info>
    <ovf:Property ovf:key="other" ovf:type="string" ovf:userConfigurable="false" ovf:value="def"/>
  </ovf:ProductSection>
</ProductSectionList>`

func decodeProductSectionList(t *testing.T, document string) *APIProductSectionList {
	t.Helper()

	var list APIProductSectionList
	require.NoError(t, xml.Unmarshal([]byte(document), &list))

	return &list
}

func TestAPIProductSectionList_ToModel(t *testing.T) {
	model := decodeProductSectionList(t, richProductSectionList).ToModel()

	require.Len(t, model.Sections, 2)

	assert.Empty(t, model.Sections[0].Class)
	assert.Equal(t, map[string]types.ModelVAppGuestProperty{
		// the Value without configuration is the value, the value attribute is only the default
		"keep": {Key: "keep", Label: "keep", Type: "string", Value: "old", UserConfigurable: true},
		"drop": {Key: "drop", Type: "string", Value: "d", UserConfigurable: true},
	}, model.Sections[0].Properties)

	assert.Equal(t, "vami", model.Sections[1].Class)
	assert.Equal(t, "x", model.Sections[1].Instance)
	assert.Equal(t, "def", model.Sections[1].Properties["other"].Value)
}

func TestMergeProductSectionList_KeepsWhatTheModelDoesNotCarry(t *testing.T) {
	merged := MergeProductSectionList(decodeProductSectionList(t, richProductSectionList), types.ModelVAppGuestProperties{Sections: []types.ModelVAppGuestPropertySection{
		{Properties: map[string]types.ModelVAppGuestProperty{
			"keep":  {Value: "new", Type: "string", UserConfigurable: true},
			"added": {Value: "a", Type: "string", UserConfigurable: true},
		}},
		{Class: "vami", Instance: "x", Properties: map[string]types.ModelVAppGuestProperty{
			"other": {Value: "def2", Type: "string"},
		}},
	}})

	out, err := xml.Marshal(merged)
	require.NoError(t, err)
	document := string(out)

	// attributes and elements of the schema that the model does not carry are kept
	assert.Contains(t, document, `:password="true"`)
	assert.Contains(t, document, `:qualifiers="MaxLen(10)"`)
	assert.Contains(t, document, "kept description")
	assert.Contains(t, document, `:configuration="big"`)
	assert.Contains(t, document, "My product")
	assert.Contains(t, document, ">Network<")
	assert.Contains(t, document, ">Other<")
	assert.Contains(t, document, ">ext<")
	assert.Contains(t, document, `:required="false"`)

	// the property that is not in the model is removed
	assert.NotContains(t, document, `:key="drop"`)

	// order: the new property is after the last property / category and before the extension element
	assert.Less(t, strings.Index(document, `:key="keep"`), strings.Index(document, `:key="added"`))
	assert.Less(t, strings.Index(document, ">Other<"), strings.Index(document, `:key="added"`))
	assert.Less(t, strings.Index(document, `:key="added"`), strings.Index(document, ">ext<"))

	// what was written can be read back with the new values
	model := decodeProductSectionList(t, document).ToModel()
	require.Len(t, model.Sections, 2)
	assert.Equal(t, "new", model.Sections[0].Properties["keep"].Value)
	assert.Equal(t, "a", model.Sections[0].Properties["added"].Value)
	assert.NotContains(t, model.Sections[0].Properties, "drop")
	assert.Equal(t, "def2", model.Sections[1].Properties["other"].Value)
}

func TestMergeProductSectionList_EmptyModelRemovesOnlyTheProperties(t *testing.T) {
	merged := MergeProductSectionList(decodeProductSectionList(t, richProductSectionList), types.ModelVAppGuestProperties{Sections: []types.ModelVAppGuestPropertySection{{Properties: map[string]types.ModelVAppGuestProperty{}}}})

	out, err := xml.Marshal(merged)
	require.NoError(t, err)

	assert.NotContains(t, string(out), "<Property")
	assert.Contains(t, string(out), "My product")
	assert.Contains(t, string(out), ">ext<")
	assert.Contains(t, string(out), "<ProductSection")
}

func TestMergeProductSectionList_NewSectionForAnotherClass(t *testing.T) {
	merged := MergeProductSectionList(decodeProductSectionList(t, richProductSectionList), types.ModelVAppGuestProperties{Sections: []types.ModelVAppGuestPropertySection{
		{Class: "other", Instance: "y", Properties: map[string]types.ModelVAppGuestProperty{"new": {Value: "n"}}},
	}})

	model := merged.ToModel()
	var found bool
	for _, section := range model.Sections {
		if section.Class == "other" && section.Instance == "y" {
			found = true
			assert.Equal(t, "n", section.Properties["new"].Value)
		}
	}
	assert.True(t, found, "the section of the new class is created")
}

func TestMergeProductSectionList_FromNothing(t *testing.T) {
	merged := MergeProductSectionList(nil, types.ModelVAppGuestProperties{Sections: []types.ModelVAppGuestPropertySection{{Properties: map[string]types.ModelVAppGuestProperty{"k": {Value: "v", UserConfigurable: true}}}}})

	out, err := xml.Marshal(merged)
	require.NoError(t, err)

	document := string(out)
	assert.True(t, strings.HasPrefix(document, "<ProductSectionList "), document)
	assert.Contains(t, document, `xmlns="http://www.vmware.com/vcloud/v1.5"`)
	assert.Contains(t, document, `xmlns:ovf="http://schemas.dmtf.org/ovf/envelope/1"`)
	assert.Contains(t, document, "<Info")
	assert.Contains(t, document, "Custom properties")

	model := decodeProductSectionList(t, document).ToModel()
	require.Len(t, model.Sections, 1)
	assert.Equal(t, types.ModelVAppGuestProperty{Key: "k", Label: "k", Type: "string", Value: "v", UserConfigurable: true}, model.Sections[0].Properties["k"])
}
