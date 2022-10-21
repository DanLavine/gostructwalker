package gostructwalker

import (
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

type StructParser struct {
	// Struct State we are currently parsing
	// NOTE: Checking this type can determine if we are working on:
	//	- an iterable: Array, Slice
	//  - a mapKey: the key of the map (should always be a simple type)
	//  - a mapValue: the actuual value of a map key
	StructState StructState

	// Name of the Field we are currently working on
	FieldName string

	// Key in the struct we are working on.
	StructField reflect.StructField

	// StrcutValue of the field we are working on
	StructValue reflect.Value

	// Tag value for the struct location we are working on
	ParsedTags Tags

	// Index for an itterable (array, slice) we are processing
	Index int
}

func NewDefaultStructParser() *StructParser {
	return &StructParser{}
}

func (sp *StructParser) GenerateFieldName(parentName string) {
	//TODO
}
