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
	FieldCallback(structParser *StructParser)
}

//counterfeiter:generate . StructWalker
type StructWalker interface {
	// Interface definition for the *structWalker struct
	//
	// PARAMS:
	// @anyStruct - Struct or pointer to a Struct that will be walked.
	//
	// RETURNS:
	// @error - nil if no errors. An immediate error if passed a non struct object
	Walk(anyStruct interface{}) error
}

type structWalker struct {
	tagParser *tagParser
	walker    Walker
}

func New(tagKey string, walker Walker) (*structWalker, error) {
	if tagKey == "" {
		return nil, fmt.Errorf("required field 'tagKey' to be set")
	}

	if walker == nil {
		return nil, fmt.Errorf("Received a nil walker")
	}

	return &structWalker{
		tagParser: &tagParser{TagKey: tagKey},
		walker:    walker,
	}, nil
}

func (s *structWalker) Walk(anyStruct interface{}) error {
	if anyStruct == nil {
		return fmt.Errorf("recieved a nil struct")
	}

	reflectValueDereference := pointerDereference(reflect.ValueOf(anyStruct))

	switch reflectValueDereference.Kind() {
	case reflect.Struct:
		return s.walkFields(nil, reflectValueDereference)
	default:
		return fmt.Errorf("Expected a struct or pointer to struct, but received a '%s'", reflectValueDereference.Kind().String())
	}
}
