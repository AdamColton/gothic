package sqlmodel

import (
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
)

// Types maps Model types to SQL types. This should be extended as new types are
// added
var Types = map[string]string{
	"bool":  "int(1) UNSIGNED DEFAULT 0 NOT NULL",
	"byte":  "int(8) UNSIGNED DEFAULT 0 NOT NULL",
	"int":   "int DEFAULT 0 NOT NULL",
	"int8":  "int(8) DEFAULT 0 NOT NULL",
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
	"uint8":  "int(8) UNSIGNED DEFAULT 0 NOT NULL",
	"uint16": "int(16) UNSIGNED DEFAULT 0 NOT NULL",
	"uint32": "int(32) UNSIGNED DEFAULT 0 NOT NULL",
	"uint64": "int(64) UNSIGNED DEFAULT 0 NOT NULL",
	// "uintptr":    gothicgo.UintptrType,
	"datetime": "DATETIME",
}

type converter struct {
	toDB   gothicgo.FuncCaller
	fromDB gothicgo.FuncCaller
}

var sqlPkg = gothicgo.MustPackageRef("github.com/adamcolotn/gsql")
var timeTime = gothicgo.DefStruct(gothicgo.MustPackageRef("time"), "Time")

var converters = map[string]*converter{
	"datetime": &converter{
		toDB:   gothicgo.FuncCall(sqlPkg, "TimeToString", gothicgo.Rets(timeTime), gothicgo.Rets(gothicgo.StringType)),
		fromDB: gothicgo.FuncCall(sqlPkg, "StringToTime", gothicgo.Rets(gothicgo.StringType), gothicgo.Rets(timeTime)),
	},
}

// AddConverter to handle conversions to and from the database
func AddConverter(name string, toDB, fromDB gothicgo.FuncCaller) error {
	toDBargs := toDB.Args()
	if len(toDBargs) != 1 {
		return fmt.Errorf("toDB FuncCaller must take 1 arg")
	}
	toDBrets := toDB.Rets()
	if len(toDBrets) != 1 {
		return fmt.Errorf("toDB FuncCaller must return 1 value")
	}
	fromDBargs := fromDB.Args()
	if len(fromDBargs) != 1 {
		return fmt.Errorf("fromDB FuncCaller must take 1 arg")
	}
	fromDBrets := fromDB.Rets()
	if len(fromDBrets) != 1 {
		return fmt.Errorf("fromDB FuncCaller must return 1 value")
	}

	if toDBargs[0].T.String() != fromDBrets[0].T.String() {
		return fmt.Errorf("toDB argument must be the same as fromDB return")
	}
	if toDBrets[0].T.String() != fromDBargs[0].T.String() {
		return fmt.Errorf("toDB return must be the same as fromDB argument")
	}

	converters[name] = &converter{
		toDB:   toDB,
		fromDB: fromDB,
	}
	return nil
}

// ZeroVals defines the zero value in SQL for Model types
var ZeroVals = map[string]string{
	"bool":  "0",
	"byte":  "0",
	"int":   "0",
	"int8":  "0",
	"int16": "0",
	"int32": "0",
	"int64": "0",
	// "complex128": gothicgo.Complex128Type,
	// "complex64":  gothicgo.Complex64Type,
	//"float32":    gothicgo.Float32Type,
	//"float64":    gothicgo.Float64Type,
	//"rune":       gothicgo.RuneType,
	"string": `""`,
	"uint":   "0",
	// "uint8":      gothicgo.Uint8Type,
	// "uint16":     gothicgo.Uint16Type,
	// "uint32":     gothicgo.Uint32Type,
	// "uint64":     gothicgo.Uint64Type,
	// "uintptr":    gothicgo.UintptrType,
	//"datetime": "DATETIME",
}
