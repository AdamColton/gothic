package gothic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testSnippet struct {
	Vals []string
}

func (t *testSnippet) New() *testSnippetInstance { return &testSnippetInstance{} }

type testSnippetInstance struct {
	Vals []string
}

func (t *testSnippetInstance) Prepare()            {}
func (t *testSnippetInstance) Set(key, val string) {}
func (t *testSnippetInstance) Generate() {
	t.Vals = []string{"test instance"}
}

func TestSC(t *testing.T) {
	sn := &testSnippet{}
	s := NewSC()
	s.AddSnippetTo("test", sn)
	assert.Equal(t, 1, len(s.Snippets("test")))

	ts := sn.New()
	si := SnippetInstance(ts)
	si.Prepare()
	si.Generate()
	assert.Equal(t, []string{"test instance"}, ts.Vals)
}
