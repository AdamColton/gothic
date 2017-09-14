package gothicgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImports(t *testing.T) {
	i := NewImports(nil)
	assert.Equal(t, "", i.String(), "Empty Import")

	i.AddNameImports("test")
	i.AddRefImports(MustPackageRef("foo/bar"))

	r := ManualResolver(make(map[string]PackageRef))
	r.Add(MustPackageRef("test/test"))

	i.ResolvePackages(r)
	expected := "import (\n\t\"foo/bar\"\n\t\"test/test\"\n)\n"
	assert.Equal(t, expected, i.String(), "Imports")
}
