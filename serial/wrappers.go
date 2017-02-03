package serial

import (
	"unsafe"
)

func MarshalUint(i uint) []byte     { return MarshalUintL(uint64(i), 8) }
func MarshalUint16(i uint16) []byte { return MarshalUintL(uint64(i), 2) }
func MarshalUint32(i uint32) []byte { return MarshalUintL(uint64(i), 4) }
func MarshalUint64(i uint64) []byte { return MarshalUintL(i, 8) }

func MarshalInt(i int) []byte     { return MarshalIntL(int64(i), 8) }
func MarshalInt16(i int16) []byte { return MarshalIntL(int64(i), 2) }
func MarshalInt32(i int32) []byte { return MarshalIntL(int64(i), 4) }
func MarshalInt64(i int64) []byte { return MarshalIntL(i, 8) }

func UnmarshalUint(b *[]byte) uint     { return uint(UnmarshalUintL(b, 8)) }
func UnmarshalUint16(b *[]byte) uint16 { return uint16(UnmarshalUintL(b, 2)) }
func UnmarshalUint32(b *[]byte) uint32 { return uint32(UnmarshalUintL(b, 4)) }
func UnmarshalUint64(b *[]byte) uint64 { return UnmarshalUintL(b, 8) }

func UnmarshalInt(b *[]byte) int     { return int(UnmarshalIntL(b, 8)) }
func UnmarshalInt16(b *[]byte) int16 { return int16(UnmarshalIntL(b, 2)) }
func UnmarshalInt32(b *[]byte) int32 { return int32(UnmarshalIntL(b, 4)) }
func UnmarshalInt64(b *[]byte) int64 { return UnmarshalIntL(b, 8) }

func MarshalString(s string) []byte    { return MarshalByteSlice([]byte(s)) }
func UnmarshalString(b *[]byte) string { return string(UnmarshalByteSlice(b)) }

func MarshalFloat32(f float32) []byte { return MarshalUint32(*(*uint32)(unsafe.Pointer(&f))) }
func MarshalFloat64(f float64) []byte { return MarshalUint64(*(*uint64)(unsafe.Pointer(&f))) }

func MarshalIntFixed(i int) []byte     { return MarshalFixedUintL(uint64(uint(i)), 8) }
func MarshalInt16Fixed(i int16) []byte { return MarshalFixedUintL(uint64(uint16(i)), 2) }
func MarshalInt32Fixed(i int32) []byte { return MarshalFixedUintL(uint64(uint32(i)), 4) }
func MarshalInt64Fixed(i int64) []byte { return MarshalFixedUintL(uint64(i), 8) }

func UnmarshalIntFixed(b *[]byte) int     { return int(UnmarshalFixedUintL(b, 8)) }
func UnmarshalInt16Fixed(b *[]byte) int16 { return int16(UnmarshalFixedUintL(b, 2)) }
func UnmarshalInt32Fixed(b *[]byte) int32 { return int32(UnmarshalFixedUintL(b, 4)) }
func UnmarshalInt64Fixed(b *[]byte) int64 { return int64(UnmarshalFixedUintL(b, 8)) }

func MarshalUintFixed(i uint) []byte     { return MarshalFixedUintL(uint64(i), 8) }
func MarshalUint16Fixed(i uint16) []byte { return MarshalFixedUintL(uint64(i), 2) }
func MarshalUint32Fixed(i uint32) []byte { return MarshalFixedUintL(uint64(i), 4) }
func MarshalUint64Fixed(i uint64) []byte { return MarshalFixedUintL(uint64(i), 8) }

func UnmarshalUintFixed(b *[]byte) uint     { return uint(UnmarshalFixedUintL(b, 8)) }
func UnmarshalUint16Fixed(b *[]byte) uint16 { return uint16(UnmarshalFixedUintL(b, 2)) }
func UnmarshalUint32Fixed(b *[]byte) uint32 { return uint32(UnmarshalFixedUintL(b, 4)) }
func UnmarshalUint64Fixed(b *[]byte) uint64 { return UnmarshalFixedUintL(b, 8) }

func UnmarshalFloat32(b *[]byte) float32 {
	i := UnmarshalUint32(b)
	return *(*float32)(unsafe.Pointer(&i))
}
func UnmarshalFloat64(b *[]byte) float64 {
	i := UnmarshalUint64(b)
	return *(*float64)(unsafe.Pointer(&i))
}

func MarshalByte(i byte) []byte {
	return []byte{i}
}
func UnmarshalByte(b *[]byte) byte {
	i := (*b)[0]
	*b = (*b)[1:]
	return i
}

func MarshalUint8(i uint8) []byte {
	return []byte{byte(i)}
}
func UnmarshalUint8(b *[]byte) uint8 {
	i := (*b)[0]
	*b = (*b)[1:]
	return uint8(i)
}

func MarshalInt8(i int8) []byte {
	var b byte
	if i < 0 {
		b = (^byte(i) << 1) | 1
	} else {
		b = (byte(i) << 1)
	}
	return []byte{b}
}
func UnmarshalInt8(b *[]byte) int8 {
	i := (*b)[0]
	*b = (*b)[1:]
	f := (i & 1)
	i >>= 1
	if f == 1 {
		i = (^i)
	}
	return int8(i)
}
