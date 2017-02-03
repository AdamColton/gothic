package lensnippet

import (
	"github.com/adamcolton/gothic"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testContainer struct {
	*gothic.SC
	*gothicgo.NameType
}

func TestLenSnippet(t *testing.T) {
	tc := &testContainer{
		SC: gothic.NewSC(),
		NameType: &gothicgo.NameType{
			N: "test",
			T: gothicgo.StringType,
		},
	}

	m := Min(tc, 5)
	if m == nil {
		t.Error("Bad")
	}

	mi := m.New()
	mi.Prepare()
	expected := []string{"if len(test) < 5 {\n\terrs.Add(\"test must be at least 5 long\")\n}"}

	assert.Equal(t, expected, mi.Generate())
}
