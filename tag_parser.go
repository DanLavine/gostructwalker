package gostructwalker

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	iterable    = "iterable["
	iterableLen = len(iterable)

	mapKey    = "mapKey["
	mapKeyLen = len(mapKey)

	mapValue    = "mapValue["
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

// TODO. This also isn't quite right. We expect `tag:"m_mapKey[...]"` to be a custome user key. This should be in the "fields", but will be parsed out
func (t *tagParser) filterTags(tag string) (*tags, error) {
	tags := &tags{}

	for i := 0; i < len(tag); i++ {
		switch tag[i] {

		// TODO add a test with custom tag flimFlam[...] for example. Make sure this still works
		case '[':
			// add the `[` tag for our indexing
			i = i + 1

			matchedBracket, err := matchBrackets(tag[i:])
			if err != nil {
				return tags, err
			}

			finalBracket := matchedBracket + i

			if i >= iterableLen && tag[i-iterableLen:i] == iterable {
				// found the 'iterable[' tags
				tags.iterable = tag[i:finalBracket]

				//remove iterable from tags already set
				// TODO is there a better way so we don't need to remove these? Kind of annoying. But the look ahead
				// is also confusing
				tags.field = tags.field[:len(tags.field)-iterableLen]
			} else if i >= mapKeyLen && tag[i-mapKeyLen:i] == mapKey {
				// found the 'mapKey[' tags
				tags.mapKeys = tag[i:finalBracket]

				//remove mapKey from tags already set
				tags.field = tags.field[:len(tags.field)-mapKeyLen]
			} else if i >= mapValueLen && tag[i-mapValueLen:i] == mapValue {
				// found the 'mapValue[' tags
				tags.mapValues = tag[i:finalBracket]

				//remove mapValue from tags already set
				tags.field = tags.field[:len(tags.field)-mapValueLen]
			} else {
				// everything else is added to
				tags.field += string(tag[i:finalBracket])
			}

			//advance 'i' to the matched bracket
			i = finalBracket
		default:
			// everything else is added to
			tags.field += string(tag[i])
		}
	}

	return tags, nil
}

// when calling this function, is should be everything after a key word:
// - iterable[, mapKey[, mapValue[
func matchBrackets(tag string) (int, error) {
	bracketCounter := 1

	for index, runeVal := range tag {
		switch runeVal {
		case '[':
			bracketCounter++
		case ']':
			bracketCounter--
		}

		if bracketCounter == 0 {
			// note the + 1. This is because we want to account for indexing
			return index, nil
		}
	}

	// this should be an error
	return 0, fmt.Errorf("Did not find a matching bracket for key '%s'", tag)
}

func (t *tagParser) splitTags(tag string) (map[string]string, error) {
	parsedTags := map[string]string{}

	if tag == "" {
		return nil, nil
	}

	tags := strings.Split(tag, ",")
	for _, tag := range tags {
		// this captures complex cases where our parser leaves tags like `isStriing=true,,,` after filtering out the complex tags "iterable:[...]" for example
		if tag == "" {
			continue
		}

		splitTag := strings.Split(tag, "=")

		if len(splitTag) == 2 {
			parsedTags[splitTag[0]] = splitTag[1]
		} else {
			return parsedTags, fmt.Errorf("Invalid tag '%s'", tag)
		}
	}

	return parsedTags, nil
}
