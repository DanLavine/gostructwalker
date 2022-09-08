package gostructwalker_test

import (
	"testing"

	"github.com/DanLavine/gostructwalker"
	"github.com/DanLavine/gostructwalker/gostructwalkerfakes"

	. "github.com/onsi/gomega"
)

func TestWalk(t *testing.T) {
	g := NewGomegaWithT(t)

	String := "string"

	t.Run("can validate a basic types like 'string'", func(t *testing.T) {
		walker := &gostructwalkerfakes.FakeWalker{}
		structWalker, err := gostructwalker.New(walker)
		g.Expect(err).ToNot(HaveOccurred())

		testStruct := struct {
			String        string
			StringPointer *string
		}{
			String:        String,
			StringPointer: &String,
		}

		structWalker.Walk(testStruct)
		g.Expect(walker.FieldCallbackCallCount()).To(Equal(2))

		field1 := walker.FieldCallbackArgsForCall(0)
		g.Expect(field1.Parent).To(BeNil())
		g.Expect(field1.Field.Name).To(Equal("String"))
		g.Expect(field1.Value.Interface()).To(Equal("string"))

		field2 := walker.FieldCallbackArgsForCall(1)
		g.Expect(field2.Parent).To(BeNil())
		g.Expect(field2.Field.Name).To(Equal("StringPointer"))
		g.Expect(field2.Value.Interface()).To(Equal(&String))
	})

	t.Run("can validate named structs", func(t *testing.T) {
		walker := &gostructwalkerfakes.FakeWalker{}
		structWalker, err := gostructwalker.New(walker)
		g.Expect(err).ToNot(HaveOccurred())

		testStruct := struct {
			NamedStruct struct {
				String string
			}
		}{
			NamedStruct: struct {
				String string
			}{
				String: String,
			},
		}

		structWalker.Walk(testStruct)
		g.Expect(walker.FieldCallbackCallCount()).To(Equal(2))

		field1 := walker.FieldCallbackArgsForCall(0)
		g.Expect(field1.Parent).To(BeNil())
		g.Expect(field1.Field.Name).To(Equal("NamedStruct"))
		g.Expect(field1.Value.Interface()).To(Equal(struct{ String string }{String: "string"}))

		field2 := walker.FieldCallbackArgsForCall(1)
		g.Expect(field2.Parent).ToNot(BeNil())
		g.Expect(field2.Parent).To(Equal(field1))
		g.Expect(field2.Field.Name).To(Equal("String"))
		g.Expect(field2.Value.Interface()).To(Equal("string"))
	})
}
