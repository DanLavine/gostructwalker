package gostructwalker

import (
	"reflect"
)

func (s *structWalker) walkFields(structParserParent *StructParser, anyStruct reflect.Value) error {
	for i := 0; i < anyStruct.NumField(); i++ {
		structFieldValue := anyStruct.Field(i)

		// This is a check for private fields that cannot return an interface. Ignore those fields
		if !structFieldValue.CanInterface() {
			continue
		}

		// parse tags for the struct's field
		structField := reflect.TypeOf(anyStruct.Interface()).Field(i)
		fieldTag := s.tagParser.getTag(structField)
		tags, err := s.tagParser.filterTags(fieldTag)
		if err != nil {
			return err
		}

		// set the tags for the truct ield
		structParser := NewDefaultStructParser()
		structParser.FullTag = fieldTag
		structParser.ParsedTags, err = s.tagParser.splitTags(tags.fieldTags)

		// finish setting up structParser fields
		structParser.Parent = structParserParent
		structParser.StructState = StructStateStruct
		structParser.StructField = structField
		structParser.StructValue = structFieldValue

		s.walker.FieldCallback(structParser)

		if err = s.traverse(structParser, tags, structFieldValue); err != nil {
			return err
		}
	}

	return nil
}

func (s *structWalker) walkIterable(structParserParent *StructParser, tags *tags, anyValue reflect.Value) error {
	for i := 0; i < anyValue.Len(); i++ {
		indexedValue := anyValue.Index(i)

		structParser := NewDefaultStructParser()
		structParser.Parent = structParserParent
		structParser.StructState = StructStateIterable
		structParser.Index = i
		structParser.StructField = structParserParent.StructField
		structParser.StructValue = indexedValue

		// On each index, we want to use tags in the `iterable:[...]` section
		structParser.FullTag = tags.iterableTags
		tags, err := s.tagParser.filterTags(tags.iterableTags)
		if err != nil {
			return err
		}

		structParser.ParsedTags, err = s.tagParser.splitTags(tags.fieldTags)
		if err != nil {
			return err
		}

		s.walker.FieldCallback(structParser)

		if err = s.traverse(structParser, tags, indexedValue); err != nil {
			return err
		}
	}

	return nil
}

//func (s *structWalker) walkMap(structParserParent *StructParser, anyValue reflect.Value) {
//	mapIter := anyValue.MapRange()
//
//	for mapIter.Next() {
//		// this is eah key value in the map
//		key := mapIter.Key()
//
//		// this is the entire value in the map. So this could be an array, or another map
//		value := mapIter.Value()
//
//		structParser := NewDefaultStructParser()
//		structParser.Parent = structParserParent
//		structParser.StructState = Map
//		structParser.MapKey = key
//		structParser.MapValue = value
//		structParser.parseTags()
//
//		s.walker.FieldCallback(structParser)
//	}
//}

func (s *structWalker) traverse(structParserParent *StructParser, tags *tags, anyValue reflect.Value) error {
	valueDereference := pointerDereference(anyValue)

	switch valueDereference.Kind() {
	case reflect.Struct:
		return s.walkFields(structParserParent, valueDereference)
	case reflect.Slice, reflect.Array:
		return s.walkIterable(structParserParent, tags, valueDereference)
	}

	return nil
}
