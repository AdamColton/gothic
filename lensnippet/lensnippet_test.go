package lensnippet

import (
	"github.com/adamcolton/gothic"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testErrSnippet string

func (es testErrSnippet) AddContext(key, value string) gothic.Snippet {
	if key == "Error" {
		es = testErrSnippet(value)
	}
	return es
}
func (es testErrSnippet) String() string {
	return "err = fmt.Errorf(\"" + string(es) + "\")"
}

func TestLenSnippet(t *testing.T) {
	m := Min(10, testErrSnippet(""))
	assert.NotNil(t, m)

	s := m.AddContext("On", "test").String()
	assert.Contains(t, s, "if len(test) < 10 {")
	assert.Contains(t, s, "err = fmt.Errorf(\"test must be at least 10 long\")")
}
