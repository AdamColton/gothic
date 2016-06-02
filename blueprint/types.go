package blueprint

var strType Type

func String() Type {
	if strType == nil {
		strType = TypeString("string")
	}
	return strType
}

var intType Type

func Int() Type {
	if intType == nil {
		intType = TypeString("int")
	}
	return intType
}

var uintType Type

func Uint() Type {
	if uintType == nil {
		uintType = TypeString("uint")
	}
	return uintType
}

var byteType Type

func Byte() Type {
	if byteType == nil {
		byteType = TypeString("byte")
	}
	return byteType
}

var uint8Type Type

func Uint8() Type {
	if uint8Type == nil {
		uint8Type = TypeString("uint8")
	}
	return uint8Type
}

var uint16Type Type

func Uint16() Type {
	if uint16Type == nil {
		uint16Type = TypeString("uint16")
	}
	return uint16Type
}

var uint32Type Type

func Uint32() Type {
	if uint32Type == nil {
		uint32Type = TypeString("uint32")
	}
	return uint32Type
}

var uint64Type Type

func Uint64() Type {
	if uint64Type == nil {
		uint64Type = TypeString("uint64")
	}
	return uint64Type
}

var int8Type Type

func Int8() Type {
	if int8Type == nil {
		int8Type = TypeString("int8")
	}
	return int8Type
}

var int16Type Type

func Intint16() Type {
	if int16Type == nil {
		int16Type = TypeString("int16")
	}
	return int16Type
}

var int32Type Type

func Intint32() Type {
	if int32Type == nil {
		int32Type = TypeString("int32")
	}
	return int32Type
}

var int64Type Type

func Intint64() Type {
	if int64Type == nil {
		int64Type = TypeString("int64")
	}
	return int64Type
}
