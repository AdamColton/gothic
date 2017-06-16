package sqlmodel

var Types = map[string]string{
	"bool":  "tinyint(1) UNSIGNED DEFAULT 0 NOT NULL",
	"byte":  "tinyint(8) UNSIGNED DEFAULT 0 NOT NULL",
	"int":   "int DEFAULT 0 NOT NULL",
	"int8":  "tinyint(8) DEFAULT 0 NOT NULL",
	"int16": "int(16) DEFAULT 0 NOT NULL",
	"int32": "int(32) DEFAULT 0 NOT NULL",
	"int64": "int(64) DEFAULT 0 NOT NULL",
	// "complex128": gothicgo.Complex128Type,
	// "complex64":  gothicgo.Complex64Type,
	//"float32":    gothicgo.Float32Type,
	//"float64":    gothicgo.Float64Type,
	//"rune":       gothicgo.RuneType,
	"string": "varchar(255) NOT NULL",
	"uint":   "int UNSIGNED DEFAULT 0 NOT NULL",
	// "uint8":      gothicgo.Uint8Type,
	// "uint16":     gothicgo.Uint16Type,
	// "uint32":     gothicgo.Uint32Type,
	// "uint64":     gothicgo.Uint64Type,
	// "uintptr":    gothicgo.UintptrType,
}

type Converter struct {
	toDB   string
	fromDB string
}

var Converters = make(map[string]*Converter)
