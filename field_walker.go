package gostructwalker

import "reflect"

func (s *structWalker) walkFields(field *Field, anyStruct reflect.Value) {
	for i := 0; i < anyStruct.NumField(); i++ {
		structFieldValue := anyStruct.Field(i)

		// This is a check for private fields that cannot return an interface. Ignore those fields
		if !structFieldValue.CanInterface() {
			continue
		}

		field := Field{
			Parent:           field,
			StructField:      reflect.TypeOf(anyStruct.Interface()).Field(i),
			StructFieldValue: structFieldValue,
		}

		s.walker.FieldCallback(field)

		s.traverse(&field, structFieldValue)
	}
}

func (s *structWalker) traverse(parentField *Field, anyValue reflect.Value) {

	valueDereference := pointerDereference(anyValue)

	switch valueDereference.Kind() {
	case reflect.Struct:
		s.walkFields(parentField, valueDereference)
	}
}
