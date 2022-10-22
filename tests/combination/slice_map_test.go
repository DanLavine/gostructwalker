package slice_map_test

import (
	"testing"

	"github.com/DanLavine/gostructwalker"
	"github.com/DanLavine/gostructwalker/gostructwalkerfakes"

	. "github.com/onsi/gomega"
)

type sliceMapFormat1 struct {
	SliceMap []map[string]string `validate:"minLength=2,iterable[required=true,mapKey[isString=true],mapValue[canCastInt=true]]"`
}

func TestWalkerCommbination_Slice_of_Maps(t *testing.T) {
	g := NewGomegaWithT(t)

	walker := &gostructwalkerfakes.FakeWalker{}
	structWalker, err := gostructwalker.New("validate", walker)
	g.Expect(err).ToNot(HaveOccurred())

	// Number of checks:
	// field1. SliceMap                   - minLength=2
	// field2. SliceMap[0]                - required=true
	// field3. SliceMap[0][key "one"]     - isString=true
	// field4. SliceMap[0]["one"] aka "1" - canCastInt=true
	// field5. SliceMap[1]                - required=true
	// field6. SliceMap[1][key "two"]     - isString=true
	// field7. SliceMap[1]["two"] aka "2" - canCastInt=true
	testStruct := sliceMapFormat1{
		SliceMap: []map[string]string{
			{"one": "1"},
			{"two": "2"},
		},
	}

	g.Expect(structWalker.Walk(testStruct)).ToNot(HaveOccurred())
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(7))

	field1 := walker.FieldCallbackArgsForCall(0)
	g.Expect(field1.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field1.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=2", Iterable: "required=true,mapKey[isString=true],mapValue[canCastInt=true]"}))
	g.Expect(field1.StructField.Name).To(Equal("SliceMap"))
	g.Expect(field1.StructValue.Interface()).To(Equal([]map[string]string{{"one": "1"}, {"two": "2"}}))

	field2 := walker.FieldCallbackArgsForCall(1)
	g.Expect(field2.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field2.Index).To(Equal(0))
	g.Expect(field2.ParsedTags).To(Equal(gostructwalker.Tags{Field: "required=true", MapKeys: "isString=true", MapValues: "canCastInt=true"}))
	g.Expect(field2.StructField.Name).To(Equal("SliceMap"))
	g.Expect(field2.StructValue.Interface()).To(Equal(map[string]string{"one": "1"}))

	field3 := walker.FieldCallbackArgsForCall(2)
	g.Expect(field3.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field3.ParsedTags).To(Equal(gostructwalker.Tags{Field: "isString=true"}))
	g.Expect(field3.StructField.Name).To(Equal("SliceMap"))
	g.Expect(field3.StructValue.Interface()).To(Equal("one"))

	field4 := walker.FieldCallbackArgsForCall(3)
	g.Expect(field4.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field4.ParsedTags).To(Equal(gostructwalker.Tags{Field: "canCastInt=true"}))
	g.Expect(field4.StructField.Name).To(Equal("SliceMap"))
	g.Expect(field4.StructValue.Interface()).To(Equal("1"))

	field5 := walker.FieldCallbackArgsForCall(4)
	g.Expect(field5.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field5.Index).To(Equal(1))
	g.Expect(field5.ParsedTags).To(Equal(gostructwalker.Tags{Field: "required=true", MapKeys: "isString=true", MapValues: "canCastInt=true"}))
	g.Expect(field5.StructField.Name).To(Equal("SliceMap"))
	g.Expect(field5.StructValue.Interface()).To(Equal(map[string]string{"two": "2"}))

	field6 := walker.FieldCallbackArgsForCall(5)
	g.Expect(field6.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field6.ParsedTags).To(Equal(gostructwalker.Tags{Field: "isString=true"}))
	g.Expect(field6.StructField.Name).To(Equal("SliceMap"))
	g.Expect(field6.StructValue.Interface()).To(Equal("two"))

	field7 := walker.FieldCallbackArgsForCall(6)
	g.Expect(field7.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field7.ParsedTags).To(Equal(gostructwalker.Tags{Field: "canCastInt=true"}))
	g.Expect(field7.StructField.Name).To(Equal("SliceMap"))
	g.Expect(field7.StructValue.Interface()).To(Equal("2"))
}
