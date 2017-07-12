package sqlmodel

import (
	"bytes"
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
)

var DefaultConn = "conn"
var ConnPackage = ""

type SQL struct {
	model       *gomodel.GoModel
	helper      *helper
	Conn        string
	ConnPackage string
	TableName   string
	Migration   string
	scanner     *gothicgo.Func
}

func New(model *gomodel.GoModel) *SQL {
	return &SQL{
		model:       model,
		Conn:        DefaultConn,
		ConnPackage: ConnPackage,
		TableName:   model.Model.Name(),
	}
}

func (s *SQL) getHelper(fields []string, allFields bool) *helper {
	if s.helper == nil {
		s.helper = createHelper(s)
	}

	if len(fields) == 0 {
		if allFields {
			s.helper.useFields = s.helper.allFields()
		} else {
			s.helper.useFields = s.helper.fields
		}
	} else {
		s.helper.useFields = make([]string, 0, len(fields))
		for _, field := range fields {
			if s.helper.fieldsMap[field] {
				s.helper.useFields = append(s.helper.useFields, field)
			}
		}
	}

	return s.helper
}

func (s *SQL) addImport() {
	if s.ConnPackage != "" {
		s.model.File().AddPackageImport(s.ConnPackage)
	}
}

func (s *SQL) genericMethod(name string, fields []string, allFields bool) *gothicgo.Method {
	s.addImport()
	m := s.model.Struct.NewMethod(name)
	m.Returns(gothicgo.Ret(gothicgo.ErrorType))
	buf := &bytes.Buffer{}
	err := templates.ExecuteTemplate(buf, name, s.getHelper(fields, allFields))
	if err != nil {
		fmt.Println(err)
	}
	m.Body = buf.String()
	return m
}

func (s *SQL) genericFunction(name string, slice bool, fields []string, allFields bool) *gothicgo.Func {
	s.addImport()
	f := s.model.Struct.File().NewFunc(name + s.model.Name())
	t := gothicgo.PointerTo(s.model.Type())
	if slice {
		t = gothicgo.SliceOf(t)
	}
	f.Returns(gothicgo.Ret(t), gothicgo.Ret(gothicgo.ErrorType))
	buf := &bytes.Buffer{}
	err := templates.ExecuteTemplate(buf, name, s.getHelper(fields, allFields))
	if err != nil {
		fmt.Println(err)
	}
	f.Body = buf.String()
	return f
}

func (s *SQL) Insert(fields ...string) *gothicgo.Method {
	return s.genericMethod("insert", fields, false)
}

func (s *SQL) Update(fields ...string) *gothicgo.Method {
	return s.genericMethod("update", fields, false)
}

var MigrationFile *gothicgo.File

func getMigrationFile() *gothicgo.File {
	if MigrationFile == nil {
		MigrationFile = gothicgo.NewPackage("db").File("migrations.gen")
	}
	MigrationFile.AddPackageImport("gsql")
	return MigrationFile
}

func (s *SQL) Create(migration string, fields ...string) string {
	s.addImport()
	file := getMigrationFile()

	buf := &bytes.Buffer{}
	h := s.getHelper(fields, false)
	h.Migration = migration
	templates.ExecuteTemplate(buf, "createTable", h)
	str := buf.String()
	file.AddCode(str)
	return str
}

func (s *SQL) Scan() *gothicgo.Func {
	scannerItfc, err := s.model.File().NewInterface("scanner")
	if err == nil {
		args := []gothicgo.Type{gothicgo.EmptyInterfaceType}
		rets := []gothicgo.Type{gothicgo.ErrorType}
		scannerItfc.AddMethod("Scan", args, rets, true)
	}
	s.scanner = s.genericFunction("scan", false, nil, true)
	s.scanner.Args = append(s.scanner.Args, gothicgo.Arg("rows", scannerItfc))
	return s.scanner
}

func (s *SQL) Select() *gothicgo.Func {
	if s.scanner == nil {
		s.Scan()
	}
	m := s.genericFunction("select", true, nil, true)
	m.Args = append(m.Args, gothicgo.Arg("where", gothicgo.StringType), gothicgo.Arg("args", gothicgo.EmptyInterfaceType))
	m.Variadic = true
	return m
}

func (s *SQL) Upsert() *gothicgo.Method {
	return s.genericMethod("upsert", nil, false)
}
