package ent

import (
	"github.com/adamcolton/gothic/serial"
	"hash/crc64"
)

var tablePolynomial = crc64.MakeTable(crc64.ISO)

// Returns the CRC64 value of a byte slice using the CRC64 ISO poly
func CRC64(b []byte) uint64 {
	hash := crc64.New(tablePolynomial)
	hash.Write(b)
	d := hash.Sum(nil)[:]
	return serial.UnmarshalUint64Fixed(&d)
}
