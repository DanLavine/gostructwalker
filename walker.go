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
	// @structParser - Current field or struct we are working on. Contains any relevent info and struct location
	//
	// NOTE: We shouldn't change any of the fields, but thats possibly a lot of data to copy in a param. So pass a reference. Might regret this
	FieldCallback(structParser *StructParser)
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
		return fmt.Errorf("Expected a struct or pointer to struct, but received a '%s'", reflectValueDereference.Kind().String())
	}

	return nil
}
