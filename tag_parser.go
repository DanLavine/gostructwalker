package gostructwalker

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	iterable = "iterable"
	mapKey   = "mapKey"
	mapValue = "mapValue"
)

type tagParser struct {
	TagKey string
}

type Tags struct {
	Field     string
	Iterable  string
	MapKeys   string
	MapValues string
}

func (t *tagParser) getTag(field reflect.StructField) string {
	return field.Tag.Get(t.TagKey)
}

// Split out tags tags
//
// RETURNS
// * Tags  - struct containing all parsed tags
// * error - any error encountered when parsing tags
func (t *tagParser) filterTags(fullTag string) (Tags, error) {
	currentTagSection := ""
	tags := Tags{}

	for i := 0; i < len(fullTag); i++ {
		switch fullTag[i] {
		case ',':
			if currentTagSection != "" {
				if tags.Field == "" {
					tags.Field = currentTagSection
				} else {
					tags.Field = fmt.Sprintf("%s,%s", tags.Field, currentTagSection)
				}

				currentTagSection = ""
			}
		case '[':
			matchedBracket, err := matchBrackets(fullTag[i:])
			if err != nil {
				return tags, err
			}

			// align to the
			finalBracket := matchedBracket + i

			switch currentTagSection {
			case iterable:
				// found the 'iterable[' tags
				// don't include the '[]'
				tags.Iterable = fullTag[i+1 : finalBracket]
				currentTagSection = ""
			case mapKey:
				// found the 'mapKey[' tags
				// don't include the '[]'
				tags.MapKeys = fullTag[i+1 : finalBracket]
				currentTagSection = ""
			case mapValue:
				// found the 'mapValue[' tags
				// don't include the '[]'
				tags.MapValues = fullTag[i+1 : finalBracket]
				currentTagSection = ""
			default:
				// everything else is added to the field tag, including the final bracket
				currentTagSection += string(fullTag[i : finalBracket+1])
			}

			//advance 'i' to the matcked bracket
			i = finalBracket
		default:
			// everything else is added to
			currentTagSection += string(fullTag[i])
		}
	}

	// add the ending tags
	if currentTagSection != "" {
		if tags.Field == "" {
			tags.Field = currentTagSection
		} else {
			tags.Field = fmt.Sprintf("%s,%s", tags.Field, currentTagSection)
		}
	}

	return tags, nil
}

// when calling this function, is should be everything after a key word:
// - iterable[, mapKey[, mapValue[
func matchBrackets(tag string) (int, error) {
	bracketCounter := 0

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
