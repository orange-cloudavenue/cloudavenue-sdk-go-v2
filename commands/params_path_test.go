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
	"testing"
)

// ----- Test structures -----

type Nested struct {
	Value string
}

type TestStruct struct {
	Name         string
	Count        int
	Ptr          *Nested
	Slice        []Nested
	SlicePtr     []*Nested
	MapStr       map[string]Nested
	MapInt       map[int]*Nested
	MapPtr       map[string]*Nested
	NestedStruct Nested
	PtrNil       *Nested
	SliceNil     []*Nested
	MapNil       map[string]Nested
}

// ----- Tests -----

func TestGetValueAtPath_SimpleFields(t *testing.T) {
	obj := TestStruct{
		Name:  "hello",
		Count: 42,
	}
	val, err := GetValueAtPath(obj, "name")
	if err != nil || val != "hello" {
		t.Errorf("expected 'hello', got %v, err=%v", val, err)
	}
	val, err = GetValueAtPath(&obj, "count")
	if err != nil || val != 42 {
		t.Errorf("expected 42, got %v, err=%v", val, err)
	}
}

func TestGetValueAtPath_NestedStruct(t *testing.T) {
	obj := TestStruct{
		NestedStruct: Nested{Value: "nested!"},
	}
	val, err := GetValueAtPath(obj, "nestedstruct.value")
	if err != nil || val != "nested!" {
		t.Errorf("expected 'nested!', got %v, err=%v", val, err)
	}
}

func TestGetValueAtPath_Pointer(t *testing.T) {
	obj := TestStruct{
		Ptr: &Nested{Value: "ptr"},
	}
	val, err := GetValueAtPath(obj, "ptr.value")
	if err != nil || val != "ptr" {
		t.Errorf("expected 'ptr', got %v, err=%v", val, err)
	}
}

func TestGetValueAtPath_NilPointer(t *testing.T) {
	obj := TestStruct{}
	_, err := GetValueAtPath(obj, "ptr.value")
	if err == nil || err.Error() != "nil pointer at 'ptr'" {
		t.Errorf("expected nil pointer error, got %v", err)
	}
}

func TestGetValueAtPath_SliceIndex(t *testing.T) {
	obj := TestStruct{
		Slice: []Nested{{Value: "A"}, {Value: "B"}},
	}
	val, err := GetValueAtPath(obj, "slice.1.value")
	if err != nil || val != "B" {
		t.Errorf("expected 'B', got %v, err=%v", val, err)
	}
}

func TestGetValueAtPath_SlicePointerIndex(t *testing.T) {
	obj := TestStruct{
		SlicePtr: []*Nested{{Value: "X"}, {Value: "Y"}},
	}
	val, err := GetValueAtPath(obj, "sliceptr.0.value")
	if err != nil || val != "X" {
		t.Errorf("expected 'X', got %v, err=%v", val, err)
	}
}

func TestGetValueAtPath_SliceIndexOutOfBounds(t *testing.T) {
	obj := TestStruct{
		Slice: []Nested{{Value: "A"}},
	}
	_, err := GetValueAtPath(obj, "slice.5.value")
	if err == nil || err.Error() != "index 5 out of bounds at 'slice'" {
		t.Errorf("expected out of bounds error, got %v", err)
	}
}

func TestGetValueAtPath_MapStringKey(t *testing.T) {
	obj := TestStruct{
		MapStr: map[string]Nested{"foo": {Value: "bar"}},
	}
	val, err := GetValueAtPath(obj, "mapstr.foo.value")
	if err != nil || val != "bar" {
		t.Errorf("expected 'bar', got %v, err=%v", val, err)
	}
}

func TestGetValueAtPath_MapPointerValue(t *testing.T) {
	obj := TestStruct{
		MapPtr: map[string]*Nested{"foo": {Value: "baz"}},
	}
	val, err := GetValueAtPath(obj, "mapptr.foo.value")
	if err != nil || val != "baz" {
		t.Errorf("expected 'baz', got %v, err=%v", val, err)
	}
}

func TestGetValueAtPath_MapIntKey(t *testing.T) {
	obj := TestStruct{
		MapInt: map[int]*Nested{7: {Value: "lucky"}},
	}
	val, err := GetValueAtPath(obj, "mapint.7.value")
	if err != nil || val != "lucky" {
		t.Errorf("expected 'lucky', got %v, err=%v", val, err)
	}
}

func TestGetValueAtPath_KindError(t *testing.T) {
	obj := TestStruct{
		Name: "onlystring",
	}
	_, err := GetValueAtPath(obj, "name.value")
	if err == nil || err.Error() != "cannot traverse kind string at 'name'" {
		t.Errorf("expected kind error, got %v", err)
	}
}

func TestGetValueAtPath_FieldNotFound(t *testing.T) {
	obj := TestStruct{}
	_, err := GetValueAtPath(obj, "nope")
	if err == nil || err.Error() != "field 'nope' not found in struct at ''" {
		t.Errorf("expected field not found error, got %v", err)
	}
}

func TestGetValueAtPath_MapKeyNotFound(t *testing.T) {
	obj := TestStruct{
		MapStr: map[string]Nested{"foo": {Value: "bar"}},
	}
	_, err := GetValueAtPath(obj, "mapstr.bar.value")
	if err == nil || err.Error() != "map key 'bar' not found at 'mapstr'" {
		t.Errorf("expected map key not found error, got %v", err)
	}
}

func TestGetValueAtPath_MapKeyCannotConvert(t *testing.T) {
	obj := TestStruct{
		MapInt: map[int]*Nested{1: {Value: "one"}},
	}
	_, err := GetValueAtPath(obj, "mapint.notanint.value")
	if err == nil || err.Error() != "cannot convert key 'notanint' to int at 'mapint': strconv.ParseInt: parsing \"notanint\": invalid syntax" {
		t.Errorf("expected cannot convert key error, got %v", err)
	}
}

func TestGetValueAtPath_EmptyPath(t *testing.T) {
	obj := TestStruct{Count: 9}
	val, err := GetValueAtPath(obj, "")
	if err != nil || val.(TestStruct).Count != 9 {
		t.Errorf("expected entire struct, got %v, err=%v", val, err)
	}
}

func TestGetValueAtPath_NilParams(t *testing.T) {
	val, err := GetValueAtPath(nil, "any")
	if err == nil || err.Error() != "params is nil" || val != nil {
		t.Errorf("expected nil params error, got val=%v, err=%v", val, err)
	}
}

func TestGetValueAtPath_FinalNilPointer(t *testing.T) {
	obj := TestStruct{PtrNil: nil}
	_, err := GetValueAtPath(obj, "ptrnil")
	if err == nil || err.Error() != "final value is nil pointer" {
		t.Errorf("expected final value is nil pointer error, got %v", err)
	}
}

func TestStoreValueAtPath_SimpleFields(t *testing.T) {
	obj := &TestStruct{}
	err := StoreValueAtPath(obj, "name", "world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Name != "world" {
		t.Errorf("expected 'world', got %v", obj.Name)
	}
	err = StoreValueAtPath(obj, "count", "99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Count != 99 {
		t.Errorf("expected 99, got %v", obj.Count)
	}
}

func TestStoreValueAtPath_NestedStruct(t *testing.T) {
	obj := &TestStruct{}
	err := StoreValueAtPath(obj, "nestedstruct.value", "deep")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.NestedStruct.Value != "deep" {
		t.Errorf("expected 'deep', got %v", obj.NestedStruct.Value)
	}
}

func TestStoreValueAtPath_Pointer(t *testing.T) {
	obj := &TestStruct{Ptr: &Nested{}}
	err := StoreValueAtPath(obj, "ptr.value", "ptrval")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Ptr.Value != "ptrval" {
		t.Errorf("expected 'ptrval', got %v", obj.Ptr.Value)
	}
}

func TestStoreValueAtPath_SliceIndex(t *testing.T) {
	obj := &TestStruct{Slice: []Nested{{}, {}}}
	err := StoreValueAtPath(obj, "slice.1.value", "B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Slice[1].Value != "B" {
		t.Errorf("expected 'B', got %v", obj.Slice[1].Value)
	}
}

func TestStoreValueAtPath_SlicePointerIndex(t *testing.T) {
	obj := &TestStruct{SlicePtr: []*Nested{{}, {}}}
	err := StoreValueAtPath(obj, "sliceptr.0.value", "X")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.SlicePtr[0].Value != "X" {
		t.Errorf("expected 'X', got %v", obj.SlicePtr[0].Value)
	}
}

func TestStoreValueAtPath_MapStringKey(t *testing.T) {
	obj := &TestStruct{MapStr: map[string]Nested{"foo": {}}}
	err := StoreValueAtPath(obj, "mapstr.foo.value", "bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.MapStr["foo"].Value != "bar" {
		t.Errorf("expected 'bar', got %v", obj.MapStr["foo"].Value)
	}
}

func TestStoreValueAtPath_MapPointerValue(t *testing.T) {
	obj := &TestStruct{MapPtr: map[string]*Nested{"foo": {}}}
	err := StoreValueAtPath(obj, "mapptr.foo.value", "baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.MapPtr["foo"].Value != "baz" {
		t.Errorf("expected 'baz', got %v", obj.MapPtr["foo"].Value)
	}
}

func TestStoreValueAtPath_MapIntKey(t *testing.T) {
	obj := &TestStruct{MapInt: map[int]*Nested{7: {}}}
	err := StoreValueAtPath(obj, "mapint.7.value", "lucky")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.MapInt[7].Value != "lucky" {
		t.Errorf("expected 'lucky', got %v", obj.MapInt[7].Value)
	}
}

func TestStoreValueAtPath_Errors(t *testing.T) {
	obj := &TestStruct{}
	// Nil params
	if err := StoreValueAtPath(nil, "name", "x"); err == nil {
		t.Error("expected error for nil params")
	}
	// Empty path
	if err := StoreValueAtPath(obj, "", "x"); err == nil {
		t.Error("expected error for empty path")
	}
	// Field not found
	if err := StoreValueAtPath(obj, "nope", "x"); err == nil {
		t.Error("expected error for field not found")
	}
	// Out of bounds
	// obj2 := &TestStruct{Slice: []Nested{{}}}
	// if err := StoreValueAtPath(obj2, "slice.2.value", "x"); err == nil {
	// 	t.Error("expected error for out of bounds")
	// }
	// Map key not found
	obj3 := &TestStruct{MapStr: map[string]Nested{"foo": {}}}
	if err := StoreValueAtPath(obj3, "mapstr.bar.value", "x"); err == nil {
		t.Error("expected error for map key not found")
	}
	// Nil pointer
	obj4 := &TestStruct{}
	if err := StoreValueAtPath(obj4, "ptr.value", "x"); err == nil {
		t.Error("expected error for nil pointer")
	}
	// Final value is nil pointer
	obj5 := &TestStruct{PtrNil: nil}
	if err := StoreValueAtPath(obj5, "ptrnil.value", "x"); err == nil {
		t.Error("expected error for final value is nil pointer")
	}
}
