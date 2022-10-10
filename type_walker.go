package gostructwalker

import (
	"reflect"
)

// TODO field name generation
func (s *structWalker) walkFields(anyStruct reflect.Value) error {
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
		structParser.ParsedTags, err = s.tagParser.splitTags(tags.field)

		// finish setting up structParser fields
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

// TODO field name generation
// TODO remove strcutParserParent
func (s *structWalker) walkIterable(structParserParent *StructParser, tags *tags, anyValue reflect.Value) error {
	for i := 0; i < anyValue.Len(); i++ {
		indexedValue := anyValue.Index(i)

		structParser := NewDefaultStructParser()
		structParser.StructState = StructStateIterable
		structParser.Index = i
		structParser.StructField = structParserParent.StructField
		structParser.StructValue = indexedValue

		// On each index, we want to use tags in the `iterable:[...]` section
		tags, err := s.tagParser.filterTags(tags.iterable)
		if err != nil {
			return err
		}

		structParser.ParsedTags, err = s.tagParser.splitTags(tags.field)
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

// TODO field name generation
// TODO remove strcutParserParent
func (s *structWalker) walkMap(structParserParent *StructParser, tags *tags, anyValue reflect.Value) error {
	mapIter := anyValue.MapRange()

	for mapIter.Next() {
		// this is the key value in the map
		key := mapIter.Key()

		structParser := NewDefaultStructParser()
		structParser.StructState = StructStateMapKey
		structParser.StructField = structParserParent.StructField
		structParser.StructValue = key

		// On each index, we want to use tags in the `iterable:[...]` section
		mapKeyTags, err := s.tagParser.filterTags(tags.mapKeys)
		if err != nil {
			return err
		}

		structParser.ParsedTags, err = s.tagParser.splitTags(mapKeyTags.field)
		if err != nil {
			return err
		}

		s.walker.FieldCallback(structParser)

		// this is the entire value in the map. So this could be an array, or another map
		value := mapIter.Value()

		structParser = NewDefaultStructParser()
		structParser.StructState = StructStateMapValue
		structParser.StructField = structParserParent.StructField
		structParser.StructValue = value

		// On each index, we want to use tags in the `iterable:[...]` section
		mapValueTags, err := s.tagParser.filterTags(tags.mapValues)
		if err != nil {
			return err
		}

		structParser.ParsedTags, err = s.tagParser.splitTags(mapValueTags.field)
		if err != nil {
			return err
		}

		s.walker.FieldCallback(structParser)

		s.traverse(structParser, mapValueTags, value)
	}

	return nil
}

// TODO field name generation
// TODO remove strcutParserParent
func (s *structWalker) traverse(structParserParent *StructParser, tags *tags, anyValue reflect.Value) error {
	valueDereference := pointerDereference(anyValue)

	switch valueDereference.Kind() {
	case reflect.Struct:
		return s.walkFields(valueDereference)
	case reflect.Slice, reflect.Array:
		return s.walkIterable(structParserParent, tags, valueDereference)
	case reflect.Map:
		return s.walkMap(structParserParent, tags, valueDereference)
	}

	return nil
}
