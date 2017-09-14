package gothicgo

import (
	"testing"
)

func TestAutoResolver(t *testing.T) {
	a := AutoResolver()
	if a.Resolve("fmt").String() != "fmt" {
		t.Error("Did not correctly resolve fmt")
	}
}
