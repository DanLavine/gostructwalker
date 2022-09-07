package gostructwalker

import (
	"fmt"
	"reflect"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Walker
type Walker interface {
	// For each Public field in a struct, this function is called.
	//
	// PARAMS:
	// @field - Current field we are working on. Contains any relevent info and struct location
	FieldCallback(field Field)
}

type structWalker struct {
	walker Walker
}

func New(walker Walker) (*structWalker, error) {
	if walker == nil {
		return nil, fmt.Errorf("Received a nil walker")
	}

	return &structWalker{
		walker: walker,
	}, nil
}

func (s *structWalker) Walk(anyStruct interface{}) error {
	if anyStruct == nil {
		return fmt.Errorf("recieved a nil struct")
	}

	reflectValueDereference := pointerDereference(reflect.ValueOf(anyStruct))

	switch reflectValueDereference.Kind() {
	case reflect.Struct:
		s.walkFields(nil, reflectValueDereference)
	default:
		return fmt.Errorf("Expected a struct or pointer to struct,  but received a '%s'", reflectValueDereference.Kind().String())
	}

	return nil
}

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
	}
}
