package gostructwalker

import (
	"fmt"
	"reflect"
	"strings"
)

type tagParser struct {
	TagKey string
}

type tags struct {
	fieldTags    string
	iterableTags string
	mapKeys      string
	mapValues    string
}

func (t *tagParser) getTag(field reflect.StructField) string {
	return field.Tag.Get(t.TagKey)
}

// Split out tags tags
//
// RETURNS
// - struct
// - iterable
// - mapKeys
// - mapValues
func (t *tagParser) filterTags(tag string) (*tags, error) {
	tags := &tags{}

	for i := 0; i < len(tag); i++ {
		switch tag[i] {
		case 'i':
			// found the 'itrable:[' key
			if len(tag) >= i+iterableLen && tag[i:i+iterableLen] == iterable {
				matchedBracket, err := matchBrackets(i+iterableLen, tag)
				if err != nil {
					return tags, err
				}

				tags.iterableTags = tag[i+iterableLen : matchedBracket]

				//advance 'i' to end of matched bracket
				i = matchedBracket

				// skip adding to struct tag
				continue
			}
		case 'm':
			// TODO mapKeys and mapValues
		}

		// everything else is added to
		tags.fieldTags += string(tag[i])
	}

	return tags, nil
}

func (t *tagParser) splitTags(tag string) (map[string]string, error) {
	parsedTags := map[string]string{}
	tags := strings.Split(tag, ",")

	if tag == "" {
		return nil, nil
	}

	for _, tag := range tags {
		splitTag := strings.Split(tag, "=")

		if len(splitTag) == 2 {
			parsedTags[splitTag[0]] = splitTag[1]
		} else {
			return parsedTags, fmt.Errorf("Invaid tag '%s'", tag)
		}
	}

	return parsedTags, nil
}

// when calling match bracket, each of the parers has already found
// keyWord:[ -- the first bracket is found
func matchBrackets(startIndex int, tag string) (int, error) {
	bracketCounter := 1

	for index, runeVal := range tag[startIndex:] {
		switch runeVal {
		case '[':
			bracketCounter++
		case ']':
			bracketCounter--
		}

		if bracketCounter == 0 {
			return index + startIndex, nil
		}
	}

	// this should be an error
	return 0, fmt.Errorf("Did not find a matching bracket for key '%s'", tag)
}
