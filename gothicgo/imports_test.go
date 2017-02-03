package gothicgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImports(t *testing.T) {
	i := NewImports()
	assert.Equal(t, "", i.String(), "Empty Import")

	i.AddPackageImport("test")
	i.AddPathImport("foo/bar")

	r := ManualResolver(map[string]string{})
	r.Add("test", "test/test")

	i.ResolvePackages(r)
	expected := "import (\n\t\"foo/bar\"\n\t\"test/test\"\n)\n"
	assert.Equal(t, expected, i.String(), "Imports")
}
