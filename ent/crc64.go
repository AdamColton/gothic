package ent

import (
	"github.com/adamcolton/gothic/serial"
	"hash/crc64"
)

var tablePolynomial = crc64.MakeTable(crc64.ISO)

// CRC64 returns the a uint64 value of a byte slice using the CRC64 ISO poly.
// This was chosen as the weakmap key for it's speed and size.
func CRC64(b []byte) uint64 {
	hash := crc64.New(tablePolynomial)
	hash.Write(b)
	d := hash.Sum(nil)[:]
	return serial.UnmarshalUint64Fixed(&d)
}
