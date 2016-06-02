package blueprint

import (
	"fmt"
	"testing"
)

func TestTypeString(t *testing.T) {
	expect := "map[foo.bar]*person.Person"
	tp, _ := typeString(expect)
	got := tp.String()
	if tp.String() != expect {
		t.Error("Expected: " + expect + " Got: " + got)
	}
}

func TestPadTest(t *testing.T) {
	if fmt.Sprintf("%-8s-", "test") != "test    -" {
		t.Error("You do not understand padding")
	}
}
