package slice_map_test

import (
	"testing"

	"github.com/DanLavine/gostructwalker"
	"github.com/DanLavine/gostructwalker/gostructwalkerfakes"

	. "github.com/onsi/gomega"
)

type mapSlice struct {
	MapSlice map[string][]string `validate:"minLength=2,mapKey[isString=true],mapValue[iterable[canCastInt=true],required=true]"`
}

func TestWalkerCommbination_Map_Of_Slice(t *testing.T) {
	g := NewGomegaWithT(t)

	walker := &gostructwalkerfakes.FakeWalker{}
	structWalker, err := gostructwalker.New("validate", walker)
	g.Expect(err).ToNot(HaveOccurred())

	// Number of checks:
	// filed1. MapSlice
	// field2. MapSlice[key "one"]
	// field3. MapSlice["one"]
	// field4. MapSlice["one"][0]
	// field5. MapSlice["one"][1]
	testStruct := mapSlice{
		MapSlice: map[string][]string{
			"one": {"1", "11"},
		},
	}

	g.Expect(structWalker.Walk(testStruct)).ToNot(HaveOccurred())
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(5))

	field1 := walker.FieldCallbackArgsForCall(0)
	g.Expect(field1.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field1.FieldName).To(Equal("MapSlice"))
	g.Expect(field1.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=2", MapKeys: "isString=true", MapValues: "iterable[canCastInt=true],required=true"}))
	g.Expect(field1.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field1.StructValue.Interface()).To(Equal(map[string][]string{"one": {"1", "11"}}))

	field2 := walker.FieldCallbackArgsForCall(1)
	g.Expect(field2.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field2.FieldName).To(Equal("MapSlice[key: one]"))
	g.Expect(field2.ParsedTags).To(Equal(gostructwalker.Tags{Field: "isString=true"}))
	g.Expect(field2.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field2.StructValue.Interface()).To(Equal("one"))

	field3 := walker.FieldCallbackArgsForCall(2)
	g.Expect(field3.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field3.FieldName).To(Equal("MapSlice[one]"))
	g.Expect(field3.ParsedTags).To(Equal(gostructwalker.Tags{Field: "required=true", Iterable: "canCastInt=true"}))
	g.Expect(field3.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field3.StructValue.Interface()).To(Equal([]string{"1", "11"}))

	field4 := walker.FieldCallbackArgsForCall(3)
	g.Expect(field4.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field4.FieldName).To(Equal("MapSlice[one][0]"))
	g.Expect(field4.Index).To(Equal(0))
	g.Expect(field4.ParsedTags).To(Equal(gostructwalker.Tags{Field: "canCastInt=true"}))
	g.Expect(field4.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field4.StructValue.Interface()).To(Equal("1"))

	field5 := walker.FieldCallbackArgsForCall(4)
	g.Expect(field5.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field5.FieldName).To(Equal("MapSlice[one][1]"))
	g.Expect(field5.Index).To(Equal(1))
	g.Expect(field5.ParsedTags).To(Equal(gostructwalker.Tags{Field: "canCastInt=true"}))
	g.Expect(field5.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field5.StructValue.Interface()).To(Equal("11"))
}
