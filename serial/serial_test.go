package serial

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

var f = rand.Float32

func TestBasicRoundTrip(t *testing.T) {
	x := uint64(300)
	expected := []byte{2, 1, 44}
	b := MarshalUintL(x, 8)
	assert.Equal(t, expected, b)
	assert.Equal(t, x, UnmarshalUintL(&b, 8))
}

func TestBasicRoundTripFixed(t *testing.T) {
	x := uint64(300)
	b := MarshalFixedUintL(x, 8)
	expected := []byte{0, 0, 0, 0, 0, 0, 1, 44}
	assert.Equal(t, expected, b)
	assert.Equal(t, x, UnmarshalFixedUintL(&b, 8))
}

func TestUint64ToByteRoundTrip(t *testing.T) {
	x := uint64(511)
	b := MarshalUint64(x)
	assert.Equal(t, x, UnmarshalUint64(&b))
}

func TestUint64ZeroToByteRoundTrip(t *testing.T) {
	x := uint64(0)
	b := MarshalUint64(x)
	assert.Equal(t, x, UnmarshalUint64(&b))
}

func TestByteSliceRoundTrip(t *testing.T) {
	b := []byte("Hello, world!")
	be := MarshalByteSlice(b)
	long_be := make([]byte, len(be)*2)
	copy(long_be, be)
	b_out := UnmarshalByteSlice(&long_be)
	s := string(b_out)
	assert.Equal(t, "Hello, world!", s, "Byte slice Round trip failed")
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

func TestIntFixedRoundTrip(t *testing.T) {
	numbers := []int{-3141, 12345, 0, 1}
	for _, i := range numbers {
		b := MarshalIntFixed(i)
		if len(b) != 8 {
			t.Error("Wrong Length")
		}
		i_out := UnmarshalIntFixed(&b)
		if i_out != i {
			t.Error(i)
			t.Error(i_out)
			t.Error("Int Fixed Round trip failed")
		}
	}
}

func TestInt16FixedRoundTrip(t *testing.T) {
	numbers := []int16{-3141, 12345, 0, 1}
	for _, i := range numbers {
		b := MarshalInt16Fixed(i)
		if len(b) != 2 {
			t.Error("Wrong Length")
		}
		i_out := UnmarshalInt16Fixed(&b)
		if i_out != i {
			t.Error(i)
			t.Error(i_out)
			t.Error("Int16 Fixed Round trip failed")
		}
	}
}

func TestInt32FixedRoundTrip(t *testing.T) {
	numbers := []int32{-3141, 12345, 0, 1, 1234567, -1234567}
	for _, i := range numbers {
		b := MarshalInt32Fixed(i)
		if len(b) != 4 {
			t.Error("Wrong Length")
		}
		i_out := UnmarshalInt32Fixed(&b)
		if i_out != i {
			t.Error(i)
			t.Error(i_out)
			t.Error("Int32 Fixed Round trip failed")
		}
	}
}

func TestInt64FixedRoundTrip(t *testing.T) {
	numbers := []int64{-3141, 12345, 0, 1, 1234567, -1234567, 1234567890, -1234567890}
	for _, i := range numbers {
		b := MarshalInt64Fixed(i)
		if len(b) != 8 {
			t.Error("Wrong Length")
		}
		i_out := UnmarshalInt64Fixed(&b)
		if i_out != i {
			t.Error(i)
			t.Error(i_out)
			t.Error("Int64 Fixed Round trip failed")
		}
	}
}

func TestUintFixedRoundTrip(t *testing.T) {
	numbers := []uint{3141, 12345, 0, 1}
	for _, i := range numbers {
		b := MarshalUintFixed(i)
		if len(b) != 8 {
			t.Error("Wrong Length")
		}
		i_out := UnmarshalUintFixed(&b)
		if i_out != i {
			t.Error(i)
			t.Error(i_out)
			t.Error("Uint Fixed Round trip failed")
		}
	}
}

func TestUint16FixedRoundTrip(t *testing.T) {
	numbers := []uint16{3141, 12345, 0, 1}
	for _, i := range numbers {
		b := MarshalUint16Fixed(i)
		if len(b) != 2 {
			t.Error("Wrong Length")
		}
		i_out := UnmarshalUint16Fixed(&b)
		if i_out != i {
			t.Error(i)
			t.Error(i_out)
			t.Error("Uint16 Fixed Round trip failed")
		}
	}
}

func TestUint32FixedRoundTrip(t *testing.T) {
	numbers := []uint32{3141, 12345, 0, 1, 1234567}
	for _, i := range numbers {
		b := MarshalUint32Fixed(i)
		if len(b) != 4 {
			t.Error("Wrong Length")
		}
		i_out := UnmarshalUint32Fixed(&b)
		if i_out != i {
			t.Error(i)
			t.Error(i_out)
			t.Error("Uint32 Fixed Round trip failed")
		}
	}
}

func TestUint64FixedRoundTrip(t *testing.T) {
	numbers := []uint64{3141, 12345, 0, 1, 1234567, 1234567890}
	for _, i := range numbers {
		b := MarshalUint64Fixed(i)
		if len(b) != 8 {
			t.Error("Wrong Length")
		}
		i_out := UnmarshalUint64Fixed(&b)
		if i_out != i {
			t.Error(i)
			t.Error(i_out)
			t.Error("Uint64 Fixed Round trip failed")
		}
	}
}
