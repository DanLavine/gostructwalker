package maps_test

import (
	"testing"

	"github.com/DanLavine/gostructwalker"
	"github.com/DanLavine/gostructwalker/gostructwalkerfakes"

	. "github.com/onsi/gomega"
)

type simpleMapKeys struct {
	MapKeys map[string]string `validate:"minLength=100,mapKey[maxLength=200]"`
}

type complexMapKeys struct {
	Name           string                    `validate:"minLength=100"`
	ComplexMapKeys map[string]complexMapKeys `validate:"minLength=200,mapKey[required=true]"`
}

func TestWalkerMapKeys_simple_types(t *testing.T) {
	g := NewGomegaWithT(t)

	walker := &gostructwalkerfakes.FakeWalker{}
	structWalker, err := gostructwalker.New("validate", walker)
	g.Expect(err).ToNot(HaveOccurred())

	// Number of checks:
	// field1. MapKeys                     - minLength=100 check
	// field2. MapKeys[Key "one"]          - required=true
	// field3. MapKeys["one"] -> value "1" - no fields to check
	testStruct := simpleMapKeys{
		MapKeys: map[string]string{"one": "1"},
	}

	structWalker.Walk(testStruct)
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(3))

	field1 := walker.FieldCallbackArgsForCall(0)
	g.Expect(field1.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field1.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=100", MapKeys: "maxLength=200"}))
	g.Expect(field1.StructField.Name).To(Equal("MapKeys"))
	g.Expect(field1.StructValue.Interface()).To(Equal(map[string]string{"one": "1"}))

	field2 := walker.FieldCallbackArgsForCall(1)
	g.Expect(field2.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field2.ParsedTags).To(Equal(gostructwalker.Tags{Field: "maxLength=200"}))
	g.Expect(field2.StructField.Name).To(Equal("MapKeys"))
	g.Expect(field2.StructValue.Interface()).To(Equal("one"))

	field3 := walker.FieldCallbackArgsForCall(2)
	g.Expect(field3.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field3.ParsedTags).To(Equal(gostructwalker.Tags{}))
	g.Expect(field3.StructField.Name).To(Equal("MapKeys"))
	g.Expect(field3.StructValue.Interface()).To(Equal("1"))
}

func TestWalkerMapKeys_complex_types(t *testing.T) {
	g := NewGomegaWithT(t)

	walker := &gostructwalkerfakes.FakeWalker{}
	structWalker, err := gostructwalker.New("validate", walker)
	g.Expect(err).ToNot(HaveOccurred())

	childStruct := complexMapKeys{
		Name: "child one",
		ComplexMapKeys: map[string]complexMapKeys{
			"grandchild one map key": {Name: "grandchild one"},
		},
	}

	// Number of checks:
	// field1. Name                        - minLength=100
	// field2. ComplexMapKeys              - minLength-200
	// field3. ComplexMapKeys[0 Key]       - required=true
	// field4. ComplexMapKeys[0] aka value - should be empty tags
	// field5. ComplexMapKeys[0].Name
	// field6. ComplexMapKeys[0].ComplexMapKeys
	// field7. ComplexMapKeys[0].ComplexMapKeys[0 Key]
	// field8. ComplexMapKeys[0].ComplexMapKeys[0] aka value
	// field9. ComplexMapKeys[0].ComplexMapKeys[0].Name
	// field10. ComplexMapKeys[0].ComplexMapKeys[0].ComplexMapKeys - this is nil, but still checked
	testStruct := complexMapKeys{
		Name: "parent",
		ComplexMapKeys: map[string]complexMapKeys{
			"child one map key": childStruct,
		},
	}

	err = structWalker.Walk(testStruct)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(10))

	field1 := walker.FieldCallbackArgsForCall(0)
	g.Expect(field1.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field1.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=100"}))
	g.Expect(field1.StructField.Name).To(Equal("Name"))
	g.Expect(field1.StructValue.Interface()).To(Equal("parent"))

	field2 := walker.FieldCallbackArgsForCall(1)
	g.Expect(field2.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field2.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=200", MapKeys: "required=true"}))
	g.Expect(field2.StructField.Name).To(Equal("ComplexMapKeys"))
	g.Expect(field2.StructValue.Interface()).To(Equal(map[string]complexMapKeys{"child one map key": childStruct}))

	field3 := walker.FieldCallbackArgsForCall(2)
	g.Expect(field3.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field3.ParsedTags).To(Equal(gostructwalker.Tags{Field: "required=true"}))
	g.Expect(field3.StructField.Name).To(Equal("ComplexMapKeys"))
	g.Expect(field3.StructValue.Interface()).To(Equal("child one map key"))

	field4 := walker.FieldCallbackArgsForCall(3)
	g.Expect(field4.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field4.ParsedTags).To(Equal(gostructwalker.Tags{}))
	g.Expect(field4.StructField.Name).To(Equal("ComplexMapKeys"))
	g.Expect(field4.StructValue.Interface()).To(Equal(childStruct))

	field5 := walker.FieldCallbackArgsForCall(4)
	g.Expect(field5.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field5.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=100"}))
	g.Expect(field5.StructField.Name).To(Equal("Name"))
	g.Expect(field5.StructValue.Interface()).To(Equal("child one"))

	field6 := walker.FieldCallbackArgsForCall(5)
	g.Expect(field6.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field6.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=200", MapKeys: "required=true"}))
	g.Expect(field6.StructField.Name).To(Equal("ComplexMapKeys"))
	g.Expect(field6.StructValue.Interface()).To(Equal(map[string]complexMapKeys{"grandchild one map key": {Name: "grandchild one"}}))

	field7 := walker.FieldCallbackArgsForCall(6)
	g.Expect(field7.StructState).To(Equal(gostructwalker.StructStateMapKey))
	g.Expect(field7.ParsedTags).To(Equal(gostructwalker.Tags{Field: "required=true"}))
	g.Expect(field7.StructField.Name).To(Equal("ComplexMapKeys"))
	g.Expect(field7.StructValue.Interface()).To(Equal("grandchild one map key"))

	field8 := walker.FieldCallbackArgsForCall(7)
	g.Expect(field8.StructState).To(Equal(gostructwalker.StructStateMapValue))
	g.Expect(field8.ParsedTags).To(Equal(gostructwalker.Tags{}))
	g.Expect(field8.StructField.Name).To(Equal("ComplexMapKeys"))
	g.Expect(field8.StructValue.Interface()).To(Equal(complexMapKeys{Name: "grandchild one"}))

	field9 := walker.FieldCallbackArgsForCall(8)
	g.Expect(field9.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field9.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=100"}))
	g.Expect(field9.StructField.Name).To(Equal("Name"))
	g.Expect(field9.StructValue.Interface()).To(Equal("grandchild one"))

	field10 := walker.FieldCallbackArgsForCall(9)
	g.Expect(field10.StructState).To(Equal(gostructwalker.StructStateStruct))
	g.Expect(field10.ParsedTags).To(Equal(gostructwalker.Tags{Field: "minLength=200", MapKeys: "required=true"}))
	g.Expect(field10.StructField.Name).To(Equal("ComplexMapKeys"))
	g.Expect(field10.StructValue.Interface()).To(BeNil())
}
