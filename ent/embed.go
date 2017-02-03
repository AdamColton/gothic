package ent

import (
	"crypto/rand"
)

// Ent can be embedded to provide an ID
type Ent struct {
	id []byte
}

// ID returns the id of the Ent
func (ent Ent) ID() []byte { return ent.id }

// New returns a new Ent with a random ID.
func New() Ent {
	b := make([]byte, 8)
	rand.Read(b)
	return Ent{b}
}

// Def creates an Ent with a given ID.
func Def(id []byte) Ent { return Ent{id} }

// Def64 creates an Ent from a uint64 ID value.
func Def64(id64 uint64) Ent {
	id := make([]byte, 8)
	for i := 0; id64 > 0; i++ {
		id[i] = byte(id64)
		id64 >>= 8
	}
	return Ent{id}
}

// Marshal returns the ID, matches Marshaling interfaces.
func (ent Ent) Marshal() []byte { return ent.id }

// Unmarshal follows the unmarshalling standard.
func Unmarshal(b *[]byte) Ent {
	e := Ent{(*b)[:8]}
	*b = (*b)[8:]
	return e
}
