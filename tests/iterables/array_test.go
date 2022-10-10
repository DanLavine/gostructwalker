package iterables_test

import (
	"testing"

	"github.com/DanLavine/gostructwalker"
	"github.com/DanLavine/gostructwalker/gostructwalkerfakes"

	. "github.com/onsi/gomega"
)

type simpleTestStruct struct {
	Strings []string `validate:"minLength=100,iterable:[maxLength=200]"`
}

type complexTestStruct struct {
	Name               string              `validate:"minLength=100"`
	ComplexTestStructs []complexTestStruct `validate:"minLength=200,iterable:[required=true]"`
}

func TestWalkerArrays_simple_types(t *testing.T) {
	g := NewGomegaWithT(t)

	walker := &gostructwalkerfakes.FakeWalker{}
	structWalker, err := gostructwalker.New("validate", walker)
	g.Expect(err).ToNot(HaveOccurred())

	testStruct := simpleTestStruct{
		Strings: []string{"one", "two"},
	}

	structWalker.Walk(testStruct)
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(3))

	field1 := walker.FieldCallbackArgsForCall(0)
	g.Expect(field1.ParsedTags).To(Equal(map[string]string{"minLength": "100"}))
	g.Expect(field1.StructField.Name).To(Equal("Strings"))
	g.Expect(field1.StructValue.Interface()).To(Equal([]string{"one", "two"}))

	field2 := walker.FieldCallbackArgsForCall(1)
	g.Expect(field2.ParsedTags).To(Equal(map[string]string{"maxLength": "200"}))
	g.Expect(field2.StructField.Name).To(Equal("Strings"))
	g.Expect(field2.StructValue.Interface()).To(Equal("one"))

	field3 := walker.FieldCallbackArgsForCall(2)
	g.Expect(field3.ParsedTags).To(Equal(map[string]string{"maxLength": "200"}))
	g.Expect(field3.StructField.Name).To(Equal("Strings"))
	g.Expect(field3.StructValue.Interface()).To(Equal("two"))
}

func TestWalkerArrays_complex_types(t *testing.T) {
	g := NewGomegaWithT(t)

	walker := &gostructwalkerfakes.FakeWalker{}
	structWalker, err := gostructwalker.New("validate", walker)
	g.Expect(err).ToNot(HaveOccurred())

	childStruct := []complexTestStruct{
		{
			Name: "child one",
			ComplexTestStructs: []complexTestStruct{
				{Name: "grandchild one"},
			},
		},
	}

	// Number of checks:
	// field1. Name                 - minLength=100 check
	// field2. ComplexTestStructs   - minLength-200 check
	// field3. ComplexTestStructs[0] - required=true check. Also Parent is #2
	// field4. ComplexTestStructs[0].Name
	// field5. ComplexTestStructs[0].ComplexTestStructs
	// field6. ComplexTestStructs[0].ComplexTestStructs[0]
	// field7. ComplexTestStructs[0].ComplexTestStructs[0].Name
	// field8. ComplexTestStructs[0].ComplexTestStructs[0].ComplexTestStructs - this is nil, but still checked
	testStruct := complexTestStruct{
		Name:               "parent",
		ComplexTestStructs: childStruct,
	}

	err = structWalker.Walk(testStruct)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(8))

	field1 := walker.FieldCallbackArgsForCall(0)
	g.Expect(field1.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field1.ParsedTags).To(Equal(map[string]string{"minLength": "100"}))
	g.Expect(field1.StructField.Name).To(Equal("Name"))
	g.Expect(field1.StructValue.Interface()).To(Equal("parent"))

	field2 := walker.FieldCallbackArgsForCall(1)
	g.Expect(field2.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field2.ParsedTags).To(Equal(map[string]string{"minLength": "200"}))
	g.Expect(field2.StructField.Name).To(Equal("ComplexTestStructs"))
	g.Expect(field2.StructValue.Interface()).To(Equal(childStruct))

	field3 := walker.FieldCallbackArgsForCall(2)
	g.Expect(field3.ParsedTags).To(Equal(map[string]string{"required": "true"}))
	g.Expect(field3.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field3.Index).To(Equal(0))
	g.Expect(field3.StructField.Name).To(Equal("ComplexTestStructs"))
	g.Expect(field3.StructValue.Interface()).To(Equal(childStruct[0]))

	field4 := walker.FieldCallbackArgsForCall(3)
	g.Expect(field4.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field4.ParsedTags).To(Equal(map[string]string{"minLength": "100"}))
	g.Expect(field4.StructField.Name).To(Equal("Name"))
	g.Expect(field4.StructValue.Interface()).To(Equal("child one"))

	field5 := walker.FieldCallbackArgsForCall(4)
	g.Expect(field5.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field5.ParsedTags).To(Equal(map[string]string{"minLength": "200"}))
	g.Expect(field5.StructField.Name).To(Equal("ComplexTestStructs"))
	g.Expect(field5.StructValue.Interface()).To(Equal([]complexTestStruct{{Name: "grandchild one"}}))

	field6 := walker.FieldCallbackArgsForCall(5)
	g.Expect(field6.ParsedTags).To(Equal(map[string]string{"required": "true"}))
	g.Expect(field6.StructState).To(Equal(gostructwalker.StructStateIterable))
	g.Expect(field6.Index).To(Equal(0))
	g.Expect(field6.StructField.Name).To(Equal("ComplexTestStructs"))
	g.Expect(field6.StructValue.Interface()).To(Equal(complexTestStruct{Name: "grandchild one"}))

	field7 := walker.FieldCallbackArgsForCall(6)
	g.Expect(field7.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field7.ParsedTags).To(Equal(map[string]string{"minLength": "100"}))
	g.Expect(field7.StructField.Name).To(Equal("Name"))
	g.Expect(field7.StructValue.Interface()).To(Equal("grandchild one"))

	field8 := walker.FieldCallbackArgsForCall(7)
	g.Expect(field8.ParsedTags).To(Equal(map[string]string{"minLength": "200"}))
	g.Expect(field8.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field8.StructField.Name).To(Equal("ComplexTestStructs"))
	g.Expect(field8.StructValue.Interface()).To(BeNil())
}
