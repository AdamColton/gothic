package sqlmodel

import (
	"fmt"
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
	"strings"
)

// QueryBuilder is used to map the model to the templates.
// fields is changed per call and defines which fields to use in a given
// template
type QueryBuilder struct {
	fields []string
	//TODO: make args a string and make it lazy
	args            []string
	fieldConverters []FieldConverter
	scanFields      string
	primaryArg      string
	Migration       string
	//err             error
	*SQL
}

func (s *SQL) queryBuilder(fields []string, allFields bool) *QueryBuilder {
	q := &QueryBuilder{
		SQL: s,
	}

	if primary, ok := q.GoModel.Field(q.Primary()); ok {
		q.primaryArg = s.arg(primary)
	}

	if len(fields) == 0 {
		gfs := s.Fields()
		q.fields = make([]string, 0, len(gfs))
		for _, field := range gfs {
			if _, sqlOk := Types[field.Type()]; !sqlOk || (field.Primary() && !allFields) {
				continue
			}
			q.append(field)
		}
	} else {
		q.fields = make([]string, 0, len(fields))
		for _, f := range fields {
			field, ok := q.GoModel.Field(f)
			if !ok {
				continue
			}
			if _, sqlOk := Types[field.Type()]; !sqlOk {
				continue
			}
			q.append(field)
		}
	}
	return q
}

func (s *SQL) arg(field gomodel.Field) string {
	arg := s.Receiver() + "." + field.Name()
	if c, ok := Converters[field.Type()]; ok {
		arg = c.toDB.Call(s.File(), arg)
	}
	return arg
}

func (q *QueryBuilder) append(field gomodel.Field) {
	q.fields = append(q.fields, field.Name())
	q.args = append(q.args, q.arg(field))
}

func (q *QueryBuilder) Fields() string {
	return strings.Join(q.fields, ", ")
}

func (q *QueryBuilder) FieldsQ() string {
	return q.IDQuote + strings.Join(q.fields, q.IDQuote+", "+q.IDQuote) + q.IDQuote
}

func (q *QueryBuilder) allFields() []string {
	fs := make([]string, len(q.fields)+1)
	copy(fs[1:], q.fields)
	fs[0] = q.Primary()
	return fs
}

func (q *QueryBuilder) QM() string {
	qm := make([]string, len(q.fields))
	for i := range qm {
		qm[i] = "?"
	}
	return strings.Join(qm, ", ")
}

func (q *QueryBuilder) Args() string {
	return strings.Join(q.args, ", ")
}

func (q *QueryBuilder) PrimaryArg() string {
	return q.primaryArg
}

func (q *QueryBuilder) Set() string {
	set := make([]string, len(q.fields))
	for i, field := range q.fields {
		set[i] = q.IDQuote + field + q.IDQuote + "=?"
	}
	return strings.Join(set, ", ")
}

func (q *QueryBuilder) PrimaryZeroVal() string {
	return ZeroVals[q.PrimaryType()]
}

func (q *QueryBuilder) Conn() string {
	return q.SQL.Conn.RelStr(q.File())
}

func (q *QueryBuilder) BackTick() string { return "`" }

func (q *QueryBuilder) ExecuteTemplate(name string) (string, error) {
	return bufpool.ExecuteTemplate(Templates, name, q)
}

func (q *QueryBuilder) GenericMethod(name string) *gothicgo.Method {
	q.addImport()
	m := q.Struct.NewMethod(name)
	m.Returns(gothicgo.Ret(gothicgo.ErrorType))
	m.Body = func() (string, error) {
		return q.ExecuteTemplate(name)
	}
	return m
}

func (q *QueryBuilder) GenericFunction(name string, slice bool) *gothicgo.Func {
	q.addImport()
	fn := q.File().NewFunc(name + q.Name())
	t := q.MethodType()
	if slice {
		t = gothicgo.SliceOf(t)
	}
	fn.Returns(gothicgo.Ret(t), gothicgo.Ret(gothicgo.ErrorType))
	fn.Body = func() (string, error) {
		return q.ExecuteTemplate(name)
	}
	return fn
}

func (q *QueryBuilder) DefineTable() string {
	var rows []string

	if q.Primary() != "" {
		rows = append(rows,
			fmt.Sprintf("\"%s\" %s", q.Primary(), q.PrimarySqlType()),
			fmt.Sprintf("PRIMARY KEY(\"%s\")", q.Primary()),
		)
	}

	for _, field := range q.fields {
		f, _ := q.Field(field)
		rows = append(rows, fmt.Sprintf("\"%s\" %s", field, f.SqlType))
	}
	return strings.Join(rows, ",\n\t\t\t")
}

type QueryBuilderField struct {
	gomodel.Field
	SqlType string
}

func (q *QueryBuilder) Field(name string) (*QueryBuilderField, bool) {
	gf, ok := q.GoModel.Field(name)
	if !ok {
		return nil, false
	}
	st, ok := Types[gf.Type()]
	if !ok {
		return nil, false
	}
	return &QueryBuilderField{
		Field:   gf,
		SqlType: st,
	}, true
}

type FieldConverter struct {
	q *QueryBuilder
	*QueryBuilderField
	*Converter
}

func (f *FieldConverter) FromDB() string {
	return f.fromDB.Call(f.q.File(), f.Name())
}

func (f *FieldConverter) ToDB() string {
	return f.toDB.Call(f.q.File(), f.q.Receiver()+"."+f.Name())
}

func (f *FieldConverter) Receiver() string {
	return f.q.Receiver()
}

func (q *QueryBuilder) populateConverters() {
	if q.fieldConverters == nil {
		r := q.Receiver()
		scanFields := make([]string, len(q.fields))
		for i, f := range q.fields {
			field, _ := q.Field(f)
			if c, ok := Converters[field.Type()]; ok {
				q.fieldConverters = append(q.fieldConverters, FieldConverter{
					q:                 q,
					QueryBuilderField: field,
					Converter:         c,
				})
				scanFields[i] = fmt.Sprintf("&(%s)", field.Name())
			} else {
				scanFields[i] = fmt.Sprintf("&(%s.%s)", r, field.Name())
			}
		}
		q.scanFields = strings.Join(scanFields, ", ")
	}
}

func (q *QueryBuilder) FieldConverters() []FieldConverter {
	q.populateConverters()
	return q.fieldConverters
}

func (q *QueryBuilder) ScanFields() string {
	q.populateConverters()
	return q.scanFields
}

func (q *QueryBuilder) Scanner() string {
	scanner := q.SQL.Scanner()
	return scanner.Name()
}
