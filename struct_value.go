package gostructwalker

import (
	"fmt"
	"reflect"
)

type StructParser struct {
	// Parent for the current field we are looking at
	Parent *StructParser

	// Key in the struct we are working on.
	// NOTE: Checking this type can determine if we are working on:
	//	- an iterable: Array, Slice, Map
	Field reflect.StructField
	// Value of the field we are working on
	Value reflect.Value

	// Index for an itterable (array, slice) we are processing
	Index int

	// Map Key value we are processing
	MapKey reflect.Value
	// Map Value we are processing. Can b 0-N values for each MapKey
	MapValue reflect.Value
}

func (sp *StructParser) GetFieldName() string {
	field := ""

	currentSP := sp

	for {
		if currentSP == nil {
			break
		}

		if field == "" {
			field = sp.Field.Name
		} else {
			field = fmt.Sprintf("%s.%s", sp.Field.Name, field)
		}

		currentSP = sp.Parent
	}

	return field
}
