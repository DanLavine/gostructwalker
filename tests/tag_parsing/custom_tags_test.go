package custom_tags_test

import (
	"testing"

	"github.com/DanLavine/gostructwalker"
	"github.com/DanLavine/gostructwalker/gostructwalkerfakes"

	. "github.com/onsi/gomega"
)

type customTagsSimilarKeyWord struct {
	CustomIterable []string       `validate:"minLength=2,m_iterable[required=true]"`
	CustomMap      map[string]int `validate:"minLength=2,m_mapKey[required=true],mapKey[isString=true]"`
	CustomMap2     map[string]int `validate:"minLength=2,mapKey[isString=true],m_mapKey[required=true]"`
}

func TestTagParsing_Custom_Keywords(t *testing.T) {
	g := NewGomegaWithT(t)

	walker := &gostructwalkerfakes.FakeWalker{}
	structWalker, err := gostructwalker.New("validate", walker)
	g.Expect(err).ToNot(HaveOccurred())

	// Number of checks:
	// field1. CustomIterable
	// field2. CustomIterable[0]
	// field3. CustomIterable[1]
	// field4. CustomMap
	// field5. CustomMap[key "str one"]
	// field6. CustomMap["str one"]
	// field7. CustomMap2
	// field8. CustomMap2[key "foo"]
	// field9. CustomMap2[2]
	testStruct := customTagsSimilarKeyWord{
		CustomIterable: []string{"one", "two"},
		CustomMap:      map[string]int{"str one": 1},
		CustomMap2:     map[string]int{"foo": 2},
	}

	g.Expect(structWalker.Walk(testStruct)).ToNot(HaveOccurred())
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(9))

	field1 := walker.FieldCallbackArgsForCall(0)
	g.Expect(field1.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field1.FieldName).To(Equal("CustomIterable"))
	g.Expect(field1.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=2,m_iterable[required=true]"}))
	g.Expect(field1.StructField.Name).To(Equal("CustomIterable"))
	g.Expect(field1.StructValue.Interface()).To(Equal([]string{"one", "two"}))

	field2 := walker.FieldCallbackArgsForCall(1)
	g.Expect(field2.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field2.FieldName).To(Equal("CustomIterable[0]"))
	g.Expect(field2.ParsedTags).To(Equal(gostructwalker.Tags{}))
	g.Expect(field2.StructField.Name).To(Equal("CustomIterable"))
	g.Expect(field2.StructValue.Interface()).To(Equal("one"))

	field3 := walker.FieldCallbackArgsForCall(2)
	g.Expect(field3.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field3.FieldName).To(Equal("CustomIterable[1]"))
	g.Expect(field3.ParsedTags).To(Equal(gostructwalker.Tags{}))
	g.Expect(field3.StructField.Name).To(Equal("CustomIterable"))
	g.Expect(field3.StructValue.Interface()).To(Equal("two"))

	field4 := walker.FieldCallbackArgsForCall(3)
	g.Expect(field4.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field4.FieldName).To(Equal("CustomMap"))
	g.Expect(field4.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=2,m_mapKey[required=true]", MapKeys: "isString=true"}))
	g.Expect(field4.StructField.Name).To(Equal("CustomMap"))
	g.Expect(field4.StructValue.Interface()).To(Equal(map[string]int{"str one": 1}))

	field5 := walker.FieldCallbackArgsForCall(4)
	g.Expect(field5.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field5.FieldName).To(Equal("CustomMap[key: str one]"))
	g.Expect(field5.ParsedTags).To(Equal(gostructwalker.Tags{Field: "isString=true"}))
	g.Expect(field5.StructField.Name).To(Equal("CustomMap"))
	g.Expect(field5.StructValue.Interface()).To(Equal("str one"))

	field6 := walker.FieldCallbackArgsForCall(5)
	g.Expect(field6.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field6.FieldName).To(Equal("CustomMap[str one]"))
	g.Expect(field6.ParsedTags).To(Equal(gostructwalker.Tags{}))
	g.Expect(field6.StructField.Name).To(Equal("CustomMap"))
	g.Expect(field6.StructValue.Interface()).To(Equal(1))

	field7 := walker.FieldCallbackArgsForCall(6)
	g.Expect(field7.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field7.FieldName).To(Equal("CustomMap2"))
	g.Expect(field7.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=2,m_mapKey[required=true]", MapKeys: "isString=true"}))
	g.Expect(field7.StructField.Name).To(Equal("CustomMap2"))
	g.Expect(field7.StructValue.Interface()).To(Equal(map[string]int{"foo": 2}))

	field8 := walker.FieldCallbackArgsForCall(7)
	g.Expect(field8.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field8.FieldName).To(Equal("CustomMap2[key: foo]"))
	g.Expect(field8.ParsedTags).To(Equal(gostructwalker.Tags{Field: "isString=true"}))
	g.Expect(field8.StructField.Name).To(Equal("CustomMap2"))
	g.Expect(field8.StructValue.Interface()).To(Equal("foo"))

	field9 := walker.FieldCallbackArgsForCall(8)
	g.Expect(field9.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field9.FieldName).To(Equal("CustomMap2[foo]"))
	g.Expect(field9.ParsedTags).To(Equal(gostructwalker.Tags{}))
	g.Expect(field9.StructField.Name).To(Equal("CustomMap2"))
	g.Expect(field9.StructValue.Interface()).To(Equal(2))
}
