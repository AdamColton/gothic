package gothicgo

import (
	"testing"
)

func TestAutoResolver(t *testing.T) {
	a := AutoResolver()
	if a.Resolve("fmt") != "fmt" {
		t.Error("Did not correctly resolve fmt")
	}
}
