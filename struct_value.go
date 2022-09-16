package gostructwalker

import (
	"fmt"
	"reflect"
)

const (
	iterable    = "iterable:["
	iterableLen = len(iterable)

	mapKey    = "mapKey:["
	mapKeyLen = len(mapKey)

	mapValue    = "mapValue:["
	mapValueLen = len(mapValue)
)

type StructState int

const (
	StructStateStruct StructState = iota
	StructStateIterable
	StructStateMap
	StructStateMapKey
	StructStateMapValue
)

type StructParser struct {
	// Parent for the current field we are looking at
	Parent *StructParser

	// Struct State we are currently parsing
	StructState StructState

	// Key in the struct we are working on.
	// NOTE: Checking this type can determine if we are working on:
	//	- an iterable: Array, Slice, Map
	StructField reflect.StructField

	// StrcutValue of the field we are working on
	// NOTE:
	//	- If this is nil, then we are working on a Map.
	//  - If this has an Index, we are wokring on an Array or Slice.
	StructValue reflect.Value

	// Tag value for the struct location we are working on
	FullTag    string
	ParsedTags map[string]string

	// Index for an itterable (array, slice) we are processing
	Index int

	// Map Key value we are processing
	MapKey reflect.Value

	// Map Value we are processing. Can b 0-N values for each MapKey
	MapValue reflect.Value
}

func NewDefaultStructParser() *StructParser {
	return &StructParser{
		ParsedTags: map[string]string{},
	}
}

func (sp *StructParser) GetFieldName() string {
	iterable := false
	field := ""

	currentSP := sp

	for {
		if currentSP == nil {
			break
		}

		switch currentSP.StructState {
		case StructStateIterable:
			iterable = true
			field = fmt.Sprintf("%s[%d]", field, currentSP.Index)
		default: // struct
			if iterable {
				// last call was iterable. Concat like `fieldName[1]`
				field = fmt.Sprintf("%s%s", currentSP.StructField.Name, field)
			} else {
				if field == "" {
					field = fmt.Sprintf("%s", currentSP.StructField.Name)
				} else {
					field = fmt.Sprintf("%s.%s", currentSP.StructField.Name, field)
				}
			}

			iterable = false
		}

		//TODO prepend parent appropriately

		currentSP = currentSP.Parent
	}

	return field
}
