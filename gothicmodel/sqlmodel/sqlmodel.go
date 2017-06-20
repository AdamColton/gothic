package sqlmodel

import (
	"bytes"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
)

var DefaultConn = "conn"

type SQL struct {
	model     *gomodel.GoModel
	helper    *helper
	Conn      string
	TableName string
	Migration string
}

func New(model *gomodel.GoModel) *SQL {
	return &SQL{
		model:     model,
		Conn:      DefaultConn,
		TableName: model.Model.Name(),
	}
}

func (s *SQL) getHelper(fields ...string) *helper {
	if s.helper == nil {
		s.helper = createHelper(s)
	}

	if len(fields) == 0 {
		s.helper.useFields = s.helper.fields
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

func (s *SQL) Insert(fields ...string) *gothicgo.Method {
	m := s.model.Struct.NewMethod("insert")
	m.Returns(gothicgo.Ret(gothicgo.ErrorType))
	buf := &bytes.Buffer{}
	templates.ExecuteTemplate(buf, "insert", s.getHelper(fields...))
	m.Body = buf.String()
	return m
}

func (s *SQL) Update(fields ...string) *gothicgo.Method {
	m := s.model.Struct.NewMethod("update")
	m.Returns(gothicgo.Ret(gothicgo.ErrorType))
	buf := &bytes.Buffer{}
	templates.ExecuteTemplate(buf, "update", s.getHelper(fields...))
	m.Body = buf.String()
	return m
}

var MigrationFile *gothicgo.File

func (s *SQL) Create(migration string, fields ...string) string {
	file := MigrationFile
	if file == nil {
		file = s.model.Struct.File()
	}
	file.AddPackageImport("gsql")
	buf := &bytes.Buffer{}
	h := s.getHelper(fields...)
	h.Migration = migration
	templates.ExecuteTemplate(buf, "createTable", h)
	str := buf.String()
	file.AddCode(str)
	return str
}
