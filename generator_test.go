package gothic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testGen struct {
	prepared  bool
	generated bool
}

func (t *testGen) Prepare() {
	t.prepared = true
}

func (t *testGen) Generate() {
	t.generated = true

}

func TestGenerator(t *testing.T) {
	tg := &testGen{}
	p := New()

	p.AddGenerator(tg)
	assert.False(t, tg.prepared, "tg.prepared should be false")
	assert.False(t, tg.generated, "tg.prepared should be false")

	p.Prepare()
	assert.True(t, tg.prepared, "tg.prepared should be true")
	assert.False(t, tg.generated, "tg.prepared should be false")

	p.Generate()
	assert.True(t, tg.prepared, "tg.prepared should be true")
	assert.True(t, tg.generated, "tg.prepared should be false")
}
