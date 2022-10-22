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
	// field6. MapSlice[key "two"]
	// field7. MapSlice["two"]
	// field8. MapSlice["two"][0]
	// field9. MapSlice["two"][1]
	testStruct := mapSlice{
		MapSlice: map[string][]string{
			"one": {"1", "11"},
			"two": {"2", "22"},
		},
	}

	g.Expect(structWalker.Walk(testStruct)).ToNot(HaveOccurred())
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(9))

	field1 := walker.FieldCallbackArgsForCall(0)
	g.Expect(field1.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field1.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=2", MapKeys: "isString=true", MapValues: "iterable[canCastInt=true],required=true"}))
	g.Expect(field1.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field1.StructValue.Interface()).To(Equal(map[string][]string{"one": {"1", "11"}, "two": {"2", "22"}}))

	field2 := walker.FieldCallbackArgsForCall(1)
	g.Expect(field2.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field2.ParsedTags).To(Equal(gostructwalker.Tags{Field: "isString=true"}))
	g.Expect(field2.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field2.StructValue.Interface()).To(Equal("one"))

	field3 := walker.FieldCallbackArgsForCall(2)
	g.Expect(field3.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field3.ParsedTags).To(Equal(gostructwalker.Tags{Field: "required=true", Iterable: "canCastInt=true"}))
	g.Expect(field3.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field3.StructValue.Interface()).To(Equal([]string{"1", "11"}))

	field4 := walker.FieldCallbackArgsForCall(3)
	g.Expect(field4.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field4.Index).To(Equal(0))
	g.Expect(field4.ParsedTags).To(Equal(gostructwalker.Tags{Field: "canCastInt=true"}))
	g.Expect(field4.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field4.StructValue.Interface()).To(Equal("1"))

	field5 := walker.FieldCallbackArgsForCall(4)
	g.Expect(field5.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field5.Index).To(Equal(1))
	g.Expect(field5.ParsedTags).To(Equal(gostructwalker.Tags{Field: "canCastInt=true"}))
	g.Expect(field5.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field5.StructValue.Interface()).To(Equal("11"))

	field6 := walker.FieldCallbackArgsForCall(5)
	g.Expect(field6.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field6.ParsedTags).To(Equal(gostructwalker.Tags{Field: "isString=true"}))
	g.Expect(field6.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field6.StructValue.Interface()).To(Equal("two"))

	field7 := walker.FieldCallbackArgsForCall(6)
	g.Expect(field7.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field7.ParsedTags).To(Equal(gostructwalker.Tags{Field: "required=true", Iterable: "canCastInt=true"}))
	g.Expect(field7.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field7.StructValue.Interface()).To(Equal([]string{"2", "22"}))

	field8 := walker.FieldCallbackArgsForCall(7)
	g.Expect(field8.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field8.Index).To(Equal(0))
	g.Expect(field8.ParsedTags).To(Equal(gostructwalker.Tags{Field: "canCastInt=true"}))
	g.Expect(field8.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field8.StructValue.Interface()).To(Equal("2"))

	field9 := walker.FieldCallbackArgsForCall(8)
	g.Expect(field9.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field9.Index).To(Equal(1))
	g.Expect(field9.ParsedTags).To(Equal(gostructwalker.Tags{Field: "canCastInt=true"}))
	g.Expect(field9.StructField.Name).To(Equal("MapSlice"))
	g.Expect(field9.StructValue.Interface()).To(Equal("22"))
}
