package maps_test

import (
	"testing"

	"github.com/DanLavine/gostructwalker"
	"github.com/DanLavine/gostructwalker/gostructwalkerfakes"

	. "github.com/onsi/gomega"
)

type simpleMapValues struct {
	MapValues map[string]string `validate:"minLength=100,mapValue[maxLength=200]"`
}

type complexMapValues struct {
	Name             string                      `validate:"minLength=100"`
	ComplexMapValues map[string]complexMapValues `validate:"minLength=200,mapValue[required=true]"`
}

func TestWalkerMapValues_simple_types(t *testing.T) {
	g := NewGomegaWithT(t)

	walker := &gostructwalkerfakes.FakeWalker{}
	structWalker, err := gostructwalker.New("validate", walker)
	g.Expect(err).ToNot(HaveOccurred())

	// Number of checks:
	// field1. MapValues                     - minLength=100 check
	// field2. MapValues[Key "one"]          - no fields to check
	// field3. MapValues["one"] -> value "1" - required=true
	testStruct := simpleMapValues{
		MapValues: map[string]string{"one": "1"},
	}

	structWalker.Walk(testStruct)
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(3))

	field1 := walker.FieldCallbackArgsForCall(0)
	g.Expect(field1.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field1.ParsedTags).To(Equal(map[string]string{"minLength": "100"}))
	g.Expect(field1.StructField.Name).To(Equal("MapValues"))
	g.Expect(field1.StructValue.Interface()).To(Equal(map[string]string{"one": "1"}))

	field2 := walker.FieldCallbackArgsForCall(1)
	g.Expect(field2.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field2.ParsedTags).To(BeNil()) // NOTE this is empty since this is the key value and we have no tags for this
	g.Expect(field2.StructField.Name).To(Equal("MapValues"))
	g.Expect(field2.StructValue.Interface()).To(Equal("one"))

	field3 := walker.FieldCallbackArgsForCall(2)
	g.Expect(field3.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field3.ParsedTags).To(Equal(map[string]string{"maxLength": "200"}))
	g.Expect(field3.StructField.Name).To(Equal("MapValues"))
	g.Expect(field3.StructValue.Interface()).To(Equal("1"))
}

func TestWalkerMapValues_complex_types(t *testing.T) {
	g := NewGomegaWithT(t)

	walker := &gostructwalkerfakes.FakeWalker{}
	structWalker, err := gostructwalker.New("validate", walker)
	g.Expect(err).ToNot(HaveOccurred())

	childStruct := complexMapValues{
		Name: "child one",
		ComplexMapValues: map[string]complexMapValues{
			"grandchild one map key": {Name: "grandchild one"},
		},
	}

	// Number of checks:
	// field1. Name                          - minLength=100
	// field2. ComplexMapValues              - minLength-200
	// field3. ComplexMapValues[0 Key]       - should be empty tags
	// field4. ComplexMapValues[0] aka value - required=true
	// field5. ComplexMapValues[0].Name
	// field6. ComplexMapValues[0].ComplexMapValues
	// field7. ComplexMapValues[0].ComplexMapValues[0 Key]
	// field8. ComplexMapValues[0].ComplexMapValues[0] aka value
	// field9. ComplexMapValues[0].ComplexMapValues[0].Name
	// field10. ComplexMapValues[0].ComplexMapValues[0].ComplexMapValues - this is nil, but still checked
	testStruct := complexMapValues{
		Name: "parent",
		ComplexMapValues: map[string]complexMapValues{
			"child one map key": childStruct,
		},
	}

	err = structWalker.Walk(testStruct)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(10))

	field1 := walker.FieldCallbackArgsForCall(0)
	g.Expect(field1.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field1.ParsedTags).To(Equal(map[string]string{"minLength": "100"}))
	g.Expect(field1.StructField.Name).To(Equal("Name"))
	g.Expect(field1.StructValue.Interface()).To(Equal("parent"))

	field2 := walker.FieldCallbackArgsForCall(1)
	g.Expect(field2.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field2.ParsedTags).To(Equal(map[string]string{"minLength": "200"}))
	g.Expect(field2.StructField.Name).To(Equal("ComplexMapValues"))
	g.Expect(field2.StructValue.Interface()).To(Equal(map[string]complexMapValues{"child one map key": childStruct}))

	field3 := walker.FieldCallbackArgsForCall(2)
	g.Expect(field3.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field3.ParsedTags).To(BeNil())
	g.Expect(field3.StructField.Name).To(Equal("ComplexMapValues"))
	g.Expect(field3.StructValue.Interface()).To(Equal("child one map key"))

	field4 := walker.FieldCallbackArgsForCall(3)
	g.Expect(field4.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field4.ParsedTags).To(Equal(map[string]string{"required": "true"}))
	g.Expect(field4.StructField.Name).To(Equal("ComplexMapValues"))
	g.Expect(field4.StructValue.Interface()).To(Equal(childStruct))

	field5 := walker.FieldCallbackArgsForCall(4)
	g.Expect(field5.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field5.ParsedTags).To(Equal(map[string]string{"minLength": "100"}))
	g.Expect(field5.StructField.Name).To(Equal("Name"))
	g.Expect(field5.StructValue.Interface()).To(Equal("child one"))

	field6 := walker.FieldCallbackArgsForCall(5)
	g.Expect(field6.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field6.ParsedTags).To(Equal(map[string]string{"minLength": "200"}))
	g.Expect(field6.StructField.Name).To(Equal("ComplexMapValues"))
	g.Expect(field6.StructValue.Interface()).To(Equal(map[string]complexMapValues{"grandchild one map key": {Name: "grandchild one"}}))

	field7 := walker.FieldCallbackArgsForCall(6)
	g.Expect(field7.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field7.ParsedTags).To(BeNil())
	g.Expect(field7.StructField.Name).To(Equal("ComplexMapValues"))
	g.Expect(field7.StructValue.Interface()).To(Equal("grandchild one map key"))

	field8 := walker.FieldCallbackArgsForCall(7)
	g.Expect(field8.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field8.ParsedTags).To(Equal(map[string]string{"required": "true"}))
	g.Expect(field8.StructField.Name).To(Equal("ComplexMapValues"))
	g.Expect(field8.StructValue.Interface()).To(Equal(complexMapValues{Name: "grandchild one"}))

	field9 := walker.FieldCallbackArgsForCall(8)
	g.Expect(field9.ParsedTags).To(Equal(map[string]string{"minLength": "100"}))
	g.Expect(field9.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field9.StructField.Name).To(Equal("Name"))
	g.Expect(field9.StructValue.Interface()).To(Equal("grandchild one"))

	field10 := walker.FieldCallbackArgsForCall(9)
	g.Expect(field10.ParsedTags).To(Equal(map[string]string{"minLength": "200"}))
	g.Expect(field10.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field10.StructField.Name).To(Equal("ComplexMapValues"))
	g.Expect(field10.StructValue.Interface()).To(BeNil())
}
