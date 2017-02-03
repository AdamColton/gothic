package gothic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testFragGen struct {
	prepared  bool
	generated bool
	value     string
}

func (t *testFragGen) Prepare() {
	t.prepared = true
}

func (t *testFragGen) Generate() []string {
	t.generated = true
	return []string{t.value}
}

func TestFragGen(t *testing.T) {
	tg := &testFragGen{
		value: "this is a test",
	}
	f := FG{}

	f.AddFragGen(tg)
	assert.False(t, tg.prepared, "tg.prepared should be false")
	assert.False(t, tg.generated, "tg.prepared should be false")

	f.Prepare()
	assert.True(t, tg.prepared, "tg.prepared should be true")
	assert.False(t, tg.generated, "tg.prepared should be false")

	s := f.Generate()
	assert.True(t, tg.prepared, "tg.prepared should be true")
	assert.True(t, tg.generated, "tg.prepared should be false")
	assert.Equal(t, []string{"this is a test"}, s)
}
