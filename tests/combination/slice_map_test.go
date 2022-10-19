package slice_map_test

import (
	"testing"

	"github.com/DanLavine/gostructwalker"
	"github.com/DanLavine/gostructwalker/gostructwalkerfakes"

	. "github.com/onsi/gomega"
)

type sliceMap struct {
	SliceMap []map[string]string `validate:"minLength=2,iterable:[required=true,mapKey[isString=true],mapValue[canCastInt=true]]"`
}

func TestWalkerCommbination_Array_of_Maps(t *testing.T) {
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
	testStruct := sliceMap{
		SliceMap: []map[string]string{
			{"one": "1"},
			{"two": "2"},
		},
	}

	g.Expect(structWalker.Walk(testStruct)).ToNot(HaveOccurred())
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(7))

	//field1 := walker.FieldCallbackArgsForCall(0)
	//g.Expect(field1.StructState).To(Equal(gostructwalker.StructStateStruct))
	//g.Expect(field1.ParsedTags).To(Equal(map[string]string{"minLength": "100"}))
	//g.Expect(field1.StructField.Name).To(Equal("MapKeys"))
	//g.Expect(field1.StructValue.Interface()).To(Equal(map[string]string{"one": "1"}))

	//field2 := walker.FieldCallbackArgsForCall(1)
	//g.Expect(field2.StructState).To(Equal(gostructwalker.StructStateMapKey))
	//g.Expect(field2.ParsedTags).To(Equal(map[string]string{"maxLength": "200"}))
	//g.Expect(field2.StructField.Name).To(Equal("MapKeys"))
	//g.Expect(field2.StructValue.Interface()).To(Equal("one"))

	//field3 := walker.FieldCallbackArgsForCall(2)
	//g.Expect(field3.StructState).To(Equal(gostructwalker.StructStateMapValue))
	//g.Expect(field3.ParsedTags).To(BeNil())
	//g.Expect(field3.StructField.Name).To(Equal("MapKeys"))
	//g.Expect(field3.StructValue.Interface()).To(Equal("1"))
}
