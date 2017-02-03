package serial

// MarshalUintL encodes a uint64 as []byte
// maxLen is based on what the underlying type was
//
// The first byte is length, unless, the length is maxLen and the first
// byte is >maxlen
//
// The algorithm works by filling from the back and then only returning as
// much of the slice as was actually used.
func MarshalUintL(i uint64, maxLen byte) []byte {
	if i == 0 {
		return []byte{0}
	}
	b := make([]byte, maxLen+1) //bytes buffer
	p := byte(maxLen)           // position
	for ; i > 0; p-- {
		b[p] = byte(i)
		i >>= 8
	}
	b[p] = maxLen - p
	if p == maxLen && b[p+1] > maxLen {
		return b[p+1:]
	}
	return b[p:]
}

// UnmarshalUintL decodes a []byte to uint64
// maxLen is based on what the underlying type was
//
// the first byte is length, unless, the length is maxLen and the first
// byte is >maxlen
func UnmarshalUintL(b *[]byte, maxLen byte) uint64 {
	if len(*b) < 1 {
		return 0
	}
	l := (*b)[0] //length
	j := byte(1) //loop index
	if l > maxLen {
		l = maxLen
		j = 0
	} else {
		l++
	}
	i := uint64(0) //int to be returned
	for ; j < l; j++ {
		i <<= 8
		i += uint64((*b)[j])
	}
	*b = (*b)[l:]
	return i
}

// MarshalIntL the least-significant bit is the sign
// converts int64 to uint64 and calls MarshalUintL
func MarshalIntL(i int64, maxLen byte) []byte {
	var i64 uint64
	if i < 0 {
		i64 = (^uint64(i) << 1) | 1
	} else {
		i64 = (uint64(i) << 1)
	}
	return MarshalUintL(i64, maxLen)
}

// UnmarshalIntL calls UnmarshalUintL, then extracts the sign from the
// least significant bit
func UnmarshalIntL(b *[]byte, maxLen byte) int64 {
	i := UnmarshalUintL(b, maxLen)
	f := (i & 1)
	i >>= 1
	if f == 1 {
		i = (^i)
	}
	return int64(i)
}

// MarshalByteSlice prepends the length of the slice using
// MarshalUintL. This structure can hold over 18,000 PB
func MarshalByteSlice(b []byte) []byte {
	l_buf := MarshalUintL(uint64(len(b)), 8)
	out := make([]byte, len(l_buf)+len(b))
	copy(out, l_buf)
	copy(out[len(l_buf):], b)
	return out
}

// UnmarshalByteSlice reads the length and returns the correct number of bits
func UnmarshalByteSlice(b *[]byte) []byte {
	l := UnmarshalUintL(b, 8) // slice length
	ret := (*b)[:l]
	*b = (*b)[l:]
	return ret
}

// MarshalFixedUintL encodes a uint64 as []byte
// the slice will always be length ln
func MarshalFixedUintL(i uint64, ln byte) []byte {
	b := make([]byte, ln) //bytes buffer
	for p := byte(ln - 1); p != 255; p-- {
		b[p] = byte(i)
		i >>= 8
	}
	return b
}

// UnmarshalFixedUintL decodes a uint64 from a byte slice. It will always use
// ln bytes.
func UnmarshalFixedUintL(b *[]byte, ln byte) uint64 {
	u := uint64(0)
	for i := byte(0); i < ln; i++ {
		u <<= 8
		u |= uint64((*b)[i])
	}
	*b = (*b)[ln:]
	return u
}
