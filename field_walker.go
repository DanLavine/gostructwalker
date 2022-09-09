package gostructwalker

import (
	"reflect"
)

func (s *structWalker) walkFields(structParserParent *StructParser, anyStruct reflect.Value) {
	for i := 0; i < anyStruct.NumField(); i++ {

		structFieldValue := anyStruct.Field(i)

		// This is a check for private fields that cannot return an interface. Ignore those fields
		if !structFieldValue.CanInterface() {
			continue
		}

		structParser := &StructParser{
			Parent: structParserParent,
			Field:  reflect.TypeOf(anyStruct.Interface()).Field(i),
			Value:  structFieldValue,
		}

		s.walker.FieldCallback(structParser)

		s.traverse(structParser, structFieldValue)
	}
}

func (s *structWalker) traverse(structParserParent *StructParser, anyValue reflect.Value) {
	valueDereference := pointerDereference(anyValue)

	switch valueDereference.Kind() {
	case reflect.Struct:
		s.walkFields(structParserParent, valueDereference)
	}
}
