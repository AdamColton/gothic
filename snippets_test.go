package gothic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testSnippet struct{}

func (t *testSnippet) New() SnippetInstance { return testSnippetInstance{} }

type testSnippetInstance struct{}

func (t testSnippetInstance) Prepare()            {}
func (t testSnippetInstance) Set(key, val string) {}
func (t testSnippetInstance) Generate() []string {
	return []string{"test instance"}
}

func TestSC(t *testing.T) {
	sn := &testSnippet{}
	s := NewSC()
	s.AddSnippetTo("test", sn)
	assert.Equal(t, 1, len(s.Snippets("test")))

	si := sn.New()
	si.Prepare()
	assert.Equal(t, []string{"test instance"}, si.Generate())
}
