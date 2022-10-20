package custom_tags_test

import (
	"testing"

	"github.com/DanLavine/gostructwalker"
	"github.com/DanLavine/gostructwalker/gostructwalkerfakes"

	. "github.com/onsi/gomega"
)

type customTagsSimilarKeyWord struct {
	CustomIterable []string `valdate:"m_iterable[required=true]"`
}

func TestTagParsing_Custom_Keywords(t *testing.T) {
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
	testStruct := customTagsSimilarKeyWord{
		CustomIterable: []string{"one", "two"},
	}

	g.Expect(structWalker.Walk(testStruct)).ToNot(HaveOccurred())
	g.Expect(walker.FieldCallbackCallCount()).To(Equal(7))
}
