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
	"reflect"
	"strings"
	"testing"
)

type StructWithBoolMap struct {
	BMap map[bool]*Nested
}

func TestGetAllValuesAtTarget_MapIntKeyPattern(t *testing.T) {
	obj := TestStruct{
		MapInt: map[int]*Nested{7: {Value: "sept"}, 42: {Value: "quarantine-deux"}},
	}
	vals, err := GetAllValuesAtTarget(obj, "mapint.42.value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vals) != 1 || vals[0] != "quarantine-deux" {
		t.Errorf("expected ['quarante-deux'], got %v", vals)
	}
	// Test erreur conversion : vérifier 'invalid syntax'
	_, err = GetAllValuesAtTarget(obj, "mapint.notanint.value")
	if err == nil || !strings.Contains(err.Error(), "invalid syntax") {
		t.Errorf("expected conversion error, got %v", err)
	}
	// Test clé non trouvée
	_, err = GetAllValuesAtTarget(obj, "mapint.99.value")
	if err == nil || !strings.Contains(err.Error(), "map key '99' not found") {
		t.Errorf("expected key not found error, got %v", err)
	}
}

func TestGetAllValuesAtTarget_MapBoolKeyPattern(t *testing.T) {
	obj := StructWithBoolMap{
		BMap: map[bool]*Nested{true: {Value: "oui"}, false: {Value: "non"}},
	}
	vals, err := GetAllValuesAtTarget(obj, "bmap.true.value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vals) != 1 || vals[0] != "oui" {
		t.Errorf("expected ['oui'], got %v", vals)
	}
	vals, err = GetAllValuesAtTarget(obj, "bmap.false.value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vals) != 1 || vals[0] != "non" {
		t.Errorf("expected ['non'], got %v", vals)
	}
	// Test erreur de conversion : vérifier 'invalid syntax'
	_, err = GetAllValuesAtTarget(obj, "bmap.notabool.value")
	if err == nil || !strings.Contains(err.Error(), "invalid syntax") {
		t.Errorf("expected conversion error, got %v", err)
	}
}

func TestGetAllValuesAtTarget_SliceIndex(t *testing.T) {
	obj := TestStruct{
		Slice: []Nested{{Value: "A"}, {Value: "B"}, {Value: "C"}},
	}
	vals, err := GetAllValuesAtTarget(obj, "slice.{index}.value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []interface{}{"A", "B", "C"}
	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("expected %v, got %v", expected, vals)
	}
}

func TestGetAllValuesAtTarget_SlicePtrIndex(t *testing.T) {
	obj := TestStruct{
		SlicePtr: []*Nested{{Value: "X"}, {Value: "Y"}},
	}
	vals, err := GetAllValuesAtTarget(obj, "sliceptr.{index}.value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []interface{}{"X", "Y"}
	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("expected %v, got %v", expected, vals)
	}
}

func TestGetAllValuesAtTarget_MapStringKey(t *testing.T) {
	obj := TestStruct{
		MapStr: map[string]Nested{"foo": {Value: "bar"}, "baz": {Value: "qux"}},
	}
	vals, err := GetAllValuesAtTarget(obj, "mapstr.{key}.value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := map[interface{}]bool{}
	for _, v := range vals {
		got[v] = true
	}
	if !got["bar"] || !got["qux"] || len(got) != 2 {
		t.Errorf("expected [bar qux], got %v", vals)
	}
}

func TestGetAllValuesAtTarget_MapPtrKey(t *testing.T) {
	obj := TestStruct{
		MapPtr: map[string]*Nested{"foo": {Value: "baz"}, "bar": {Value: "qux"}},
	}
	vals, err := GetAllValuesAtTarget(obj, "mapptr.{key}.value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := map[interface{}]bool{}
	for _, v := range vals {
		got[v] = true
	}
	if !got["baz"] || !got["qux"] || len(got) != 2 {
		t.Errorf("expected [baz qux], got %v", vals)
	}
}

func TestGetAllValuesAtTarget_NestedSliceMap(t *testing.T) {
	obj := struct {
		List []map[string]*Nested
	}{
		List: []map[string]*Nested{
			{"foo": &Nested{Value: "a1"}, "bar": &Nested{Value: "b1"}},
			{"foo": &Nested{Value: "a2"}, "bar": &Nested{Value: "b2"}},
		},
	}
	vals, err := GetAllValuesAtTarget(obj, "list.{index}.{key}.value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := map[interface{}]bool{}
	for _, v := range vals {
		got[v] = true
	}
	for _, want := range []string{"a1", "b1", "a2", "b2"} {
		if !got[want] {
			t.Errorf("expected %v in result, got %v", want, vals)
		}
	}
	if len(got) != 4 {
		t.Errorf("expected 4 results, got %v", vals)
	}
}

func TestGetAllValuesAtTarget_EmptyPattern(t *testing.T) {
	obj := TestStruct{Count: 123}
	vals, err := GetAllValuesAtTarget(obj, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vals) != 1 || vals[0].(TestStruct).Count != 123 {
		t.Errorf("expected the struct itself, got %v", vals)
	}
}

func TestGetAllValuesAtTarget_InvalidPattern(t *testing.T) {
	obj := TestStruct{}
	_, err := GetAllValuesAtTarget(obj, "notfound.{index}")
	if err == nil {
		t.Errorf("expected error for invalid struct field, got nil")
	}
}
