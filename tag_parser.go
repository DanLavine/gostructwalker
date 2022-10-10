package gostructwalker

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	iterable    = "iterable:["
	iterableLen = len(iterable)

	mapKey    = "mapKey:["
	mapKeyLen = len(mapKey)

	mapValue    = "mapValue:["
	mapValueLen = len(mapValue)
)

type tagParser struct {
	TagKey string
}

type tags struct {
	field     string
	iterable  string
	mapKeys   string
	mapValues string
}

func (t *tagParser) getTag(field reflect.StructField) string {
	return field.Tag.Get(t.TagKey)
}

// Split out tags tags
//
// RETURNS
// * tags  - struct containing all parsed tags
// * error - any error encountered when parsing tags
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

				tags.iterable = tag[i+iterableLen : matchedBracket]

				//advance 'i' to end of matched bracket
				i = matchedBracket

				// skip adding to struct tag
				continue
			}
		case 'm':
			if len(tag) >= i+mapKeyLen && tag[i:i+mapKeyLen] == mapKey {
				// found the 'mapKey:[' tags
				matchedBracket, err := matchBrackets(i+mapKeyLen, tag)
				if err != nil {
					return tags, err
				}

				tags.mapKeys = tag[i+mapKeyLen : matchedBracket]

				//advance 'i' to end of matched bracket
				i = matchedBracket

				// skip adding to struct tag
				continue
			} else if len(tag) >= i+mapValueLen && tag[i:i+mapValueLen] == mapValue {
				// found the 'mapValue:[' tags
				matchedBracket, err := matchBrackets(i+mapValueLen, tag)
				if err != nil {
					return tags, err
				}

				tags.mapValues = tag[i+mapValueLen : matchedBracket]

				//advance 'i' to end of matched bracket
				i = matchedBracket

				// skip adding to struct tag
				continue

			}
		}

		// everything else is added to
		tags.field += string(tag[i])
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
