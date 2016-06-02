package entity

import (
	"crypto/rand"
	"github.com/adamcolton/gothic/serial"
)

type Ent struct {
	id uint64
}

func (ent Ent) ID() uint64 { return ent.id }
func idToBytes(id uint64) []byte {
	s := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b := byte(id & 255)
		id >>= 8
		s[i] = b
	}
	return s
}

func New() Ent {
	b := make([]byte, 8)
	rand.Read(b)
	id := uint64(b[0]) + 256*uint64(b[1]) + 65536*uint64(b[2]) + 16777216*uint64(b[3]) + 4294967296*uint64(b[4]) + 1099511627776*uint64(b[5]) + 281474976710656*uint64(b[6]) + 72057594037927936*uint64(b[7])
	return Ent{
		id: id,
	}
}

func Def(id uint64) Ent {
	return Ent{
		id: id,
	}
}

func (ent Ent) Marshal() []byte { return serial.MarshalUint64(ent.id) }
func Unmarshal(b *[]byte) Ent   { return Def(serial.UnmarshalUint64(b)) }
