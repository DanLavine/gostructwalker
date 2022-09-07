package gostructwalker

import "reflect"

type Field struct {
	// Parent for the current field we are looking at
	Parent *Field

	// Key in the struct we are working on.
	// NOTE: Checking this type can determine if we are working on:
	//	- an iterable: Array, Slice, Map
	StructField reflect.StructField
	// Value of the field we are working on
	StructFieldValue reflect.Value

	// Index for an array element we are processing
	Index int

	// Map Key value we are processing
	MapKey reflect.Value
	// Map Value we are processing. Can b 0-N values for each MapKey
	MapValue reflect.Value
}
