# gostructwalker

A helper library for traversing golang sttructs

# Tag Parsing

This package comes with a nice ability to parse complex tags setup for specific fields,
but we need to follow a few rules. Take this example:

```
type MyStruct struct {
  ComplexType []map[string][]interface{} `tags:"required=true,iterable[minLength=2,maxLength=10,mapKey[isCapitalized=true],mapValue[isString=true]]"`
}
```

In this case we have quite a complex nested object. The tags can be split into the following pieces:

* `required=true`
  - This at the field level (against the entire []map[string][]int object). So we could do some
    validation to ensure that this field is not set to nil
* `iterable[...]`
  - This denotes that each index in the Slice will follow the rules listed here. This is because we recurse through
    all values assigned to this variable
* `minLength=2,maxLength=10` (from iterable[...])
  - So each index in the slice can check to make sure that a map has a minLength=2 and a maxLength=10
* `mapKey[isCapitalized=true]` (from iterable[...])
  - Defines the recursive tags through our Map's keys
* `isCapitalzed=true` (from mapKey[...])
  - A check can be used here to ensure that each of the keys in our map begins with a capital letter
* `mapValue[isString=true]` (from iterable[...])
  - Is now the recursive tags for each of the Map's values. So in this case with it being an interface{}
    we could write a custom check to ensure that the field's type is actually a string

With such a complicated way of writting nested objects and types, we need to follow some specifi rules

### Tag Rules

1. Each `[` character must have a matching `]` bracket. The `[...]` define recursive tag operations
1. `iterable[...]` defines the tags to provide to each index in an iterable (Slice, Array)
1. `mapKey[...]` defines the tags to provide for each index of a Map's keys
1. `mapValue[...]` defines the tags to provide for each index of a Map's values
