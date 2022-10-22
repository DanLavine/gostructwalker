package gostructwalker

import (
	"fmt"
	"reflect"
)

type StructState int

const (
	// first one rename to field?
	StructStateStruct StructState = iota
	StructStateIterable
	StructStateMapKey
	StructStateMapValue
)

// StructParser contains all the information about the current field we are parsing through
type StructParser struct {
	// Name of the Field we are currently working on
	FieldName string

	// Struct State we are currently parsing
	// NOTE: Checking this type can determine if we are working on:
	//	- an iterable: Array, Slice
	//  - a mapKey: the key of the map (should always be a simple type)
	//  - a mapValue: the actuual value of a map key
	StructState StructState

	// Key in the struct we are working on.
	StructField reflect.StructField

	// StrcutValue of the field we are working on
	StructValue reflect.Value

	// Tag value for the struct location we are working on
	ParsedTags Tags

	// Index for an itterable (array, slice) we are processing
	Index int
}

// Generate an empty *StructParser
func NewDefaultStructParser() *StructParser {
	return &StructParser{}
}

// TODO add an error
func (sp *StructParser) generateCurrentName(parentName string) error {
	switch sp.StructState {
	case StructStateStruct:
		if parentName == "" {
			sp.FieldName = sp.StructField.Name
		} else {
			sp.FieldName = fmt.Sprintf("%v.%v", parentName, sp.StructField.Name)
		}
	case StructStateIterable:
		// iterables cannot ever be first, so don't need to check the empty case
		sp.FieldName = fmt.Sprintf("%v[%d]", parentName, sp.Index)
	default:
		return fmt.Errorf("Interanl struct parsing error. Unexpected struct state. Required Struct or Iterable")
	}

	return nil
}

// TODO add an error
func (sp *StructParser) generateCurrentNameMap(parentName string, key interface{}) error {
	switch sp.StructState {
	case StructStateMapKey:
		// map keys cannot ever be first, so don't need to check the empty case
		sp.FieldName = fmt.Sprintf("%s[key: %v]", parentName, key)
	case StructStateMapValue:
		// map values cannot ever be first, so don't need to check the empty case
		sp.FieldName = fmt.Sprintf("%v[%v]", parentName, key)
	default:
		return fmt.Errorf("Interanl struct parsing error. Unexpected struct state. Required MapKey  or MapValue")
	}

	return nil
}
