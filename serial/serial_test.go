package serial

import (
	"bytes"
	"math/rand"
	"testing"
)

var f = rand.Float32

func TestBasicRoundTrip(t *testing.T) {
	x := uint64(300)
	b := MarshalUintL(x, 8)
	expect := []byte{2, 1, 44}
	if !bytes.Equal(b, expect) {
		t.Error("Marshal Uint failed")
	}

	y := UnmarshalUintL(&b, 8)
	if y != 300 {
		t.Error(y)
		t.Error("Unmarshal Uint failed")
	}
}

func TestUint64ToByteRoundTrip(t *testing.T) {
	x := uint64(511)
	b := MarshalUint64(x)

	y := UnmarshalUint64(&b)
	if x != y {
		t.Error("should be equal")
	}
}

func TestUint64ZeroToByteRoundTrip(t *testing.T) {
	x := uint64(0)
	b := MarshalUint64(x)

	y := UnmarshalUint64(&b)
	if x != y {
		t.Error("should be equal")
	}
}

func TestByteSliceRoundTrip(t *testing.T) {
	b := []byte("Hello, world!")
	be := MarshalByteSlice(b)
	long_be := make([]byte, len(be)*2)
	copy(long_be, be)
	b_out := UnmarshalByteSlice(&long_be)
	s := string(b_out)
	if s != "Hello, world!" {
		t.Error("Byte slice Round trip failed" + s)
	}
}

func TestInt16RoundTrip(t *testing.T) {
	i := int16(-3141)
	b := MarshalInt16(i)

	long_be := make([]byte, len(b)*2)
	copy(long_be, b)

	i_out := UnmarshalInt16(&long_be)
	if i_out != i {
		t.Error("Int16 Round trip failed")
	}
}

func TestTwoStrings(t *testing.T) {
	s1 := "Hello"
	s2 := "World"
	b := MarshalString(s1)
	b = append(b, MarshalString(s2)...)

	s1_o := UnmarshalString(&b)
	s2_o := UnmarshalString(&b)
	if s1 != s1_o {
		t.Error("Failed to recover string 1")
	}
	if s2 != s2_o {
		t.Error("Failed to recover string 2")
	}
}

func TestFloat32RoundTrip(t *testing.T) {
	f := float32(3.141592653)
	b := MarshalFloat32(f)

	long_be := make([]byte, len(b)*2)
	copy(long_be, b)

	f_out := UnmarshalFloat32(&long_be)
	if f_out != f {
		t.Error("Float32 Round trip failed")
	}
}

func TestCastByte(t *testing.T) {
	i := 771
	b := byte(i)
	if b != 3 {
		t.Error("That does not work the way you think")
	}
}

func TestUint32RoundTripMany(t *testing.T) {
	for j := 0; j < 100; j++ {
		i := rand.Uint32()
		b := MarshalUint32(i)
		i_out := UnmarshalUint32(&b)
		if i_out != i {
			t.Error(i)
			t.Error("Int Round trip failed")
		}
	}
}

func TestInt8RoundTripMany(t *testing.T) {
	numbers := []int8{127, -128, 0, 1, -1, 55, -23}
	for _, i := range numbers {
		b := MarshalInt8(i)
		i_out := UnmarshalInt8(&b)
		if i_out != i {
			t.Error(i)
			t.Error(i_out)
			t.Error("Int Round trip failed")
		}
	}
}

func TestInt16MinRoundTripMany(t *testing.T) {
	i := int16(-32768)
	b := MarshalInt16(i)
	i_out := UnmarshalInt16(&b)
	if i_out != i {
		t.Error(i)
		t.Error(i_out)
		t.Error("Int Round trip failed")
	}
}
