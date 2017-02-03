package ent

import (
	"crypto/rand"
)

// Ent can be embedded to provide an ID
type Ent struct {
	id []byte
}

func (ent Ent) ID() []byte { return ent.id }

func New() Ent {
	b := make([]byte, 8)
	rand.Read(b)
	return Ent{b}
}

func Def(id []byte) Ent { return Ent{id} }
func Def64(id64 uint64) Ent {
	id := make([]byte, 8)
	for i := 0; id64 > 0; i++ {
		id[i] = byte(id64)
		id64 >>= 8
	}
	return Ent{id}
}

func (ent Ent) Marshal() []byte { return ent.id }
func Unmarshal(b *[]byte) Ent {
	e := Ent{(*b)[:8]}
	*b = (*b)[8:]
	return e
}
