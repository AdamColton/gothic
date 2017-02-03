package gothicserial

import (
	"github.com/adamcolton/gothic/gothicgo"
)

// Defines the Marshal and Unmarshal call
type SerializeDef interface {
	Marshal(v, pkg string) string
	Unmarshal(v, pkg string) string
	PackageName() string
	File() *gothicgo.File
}
