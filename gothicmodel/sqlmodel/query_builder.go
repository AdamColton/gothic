package sqlmodel

import (
	"fmt"
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
	"io"
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

// Fields returns all the fields used in the query builder as a comma delimited
// list
func (q *QueryBuilder) Fields() string {
	return strings.Join(q.fields, ", ")
}

// FieldsQ returns all the fields used in the query builder as a comma delimited
// list with each field wrapped in an IDQuote
func (q *QueryBuilder) FieldsQ() string {
	return q.IDQuote + strings.Join(q.fields, q.IDQuote+", "+q.IDQuote) + q.IDQuote
}

func (q *QueryBuilder) allFields() []string {
	fs := make([]string, len(q.fields)+1)
	copy(fs[1:], q.fields)
	fs[0] = q.Primary()
	return fs
}

// QM returns a list of question marks the length of Fields for constructing
// a query.
func (q *QueryBuilder) QM() string {
	qm := make([]string, len(q.fields))
	for i := range qm {
		qm[i] = "?"
	}
	return strings.Join(qm, ", ")
}

// Args returns the fields as args to be passed into a call to Exec. Fields that
// need to pass through a converter will be replaced by a local variable.
func (q *QueryBuilder) Args() string {
	return strings.Join(q.args, ", ")
}

// PrimaryArg returns the arg for the primary fields
func (q *QueryBuilder) PrimaryArg() string {
	return q.primaryArg
}

// Set returns each field followed by "=?" for a call to set.
func (q *QueryBuilder) Set() string {
	set := make([]string, len(q.fields))
	for i, field := range q.fields {
		set[i] = q.IDQuote + field + q.IDQuote + "=?"
	}
	return strings.Join(set, ", ")
}

// PrimaryZeroVal gets the zero value of the primary field
func (q *QueryBuilder) PrimaryZeroVal() string {
	return ZeroVals[q.PrimaryType()]
}

// Conn gets the SQL connection as a string
func (q *QueryBuilder) Conn() string {
	return q.SQL.Conn.RelStr(q.File())
}

// BackTick allows a backtick to be injected into a template
func (q *QueryBuilder) BackTick() string { return "`" }

// ExecuteTemplate by name and pass the QueryBuilder in as the data.
func (q *QueryBuilder) ExecuteTemplate(name string) (string, error) {
	return bufpool.ExecuteTemplate(Templates, name, q)
}

type QueryBuilderTemplateWriteTo struct {
	*QueryBuilder
	Name string
}

func (qbtw *QueryBuilderTemplateWriteTo) WriteTo(w io.Writer) (int64, error) {
	buf := bufpool.Get()
	err := Templates.ExecuteTemplate(buf, qbtw.Name, qbtw.QueryBuilder)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(buf.Bytes())
	return int64(n), err
}

// TemplateWriteTo returns an object that fulfils WriterTo and will write the
// template to the writer
func (q *QueryBuilder) TemplateWriteTo(name string) io.WriterTo {
	return &QueryBuilderTemplateWriteTo{
		QueryBuilder: q,
		Name:         name,
	}
}

// GenericMethod takes a template name and wraps it in a method on the struct.
func (q *QueryBuilder) GenericMethod(name string) *gothicgo.Method {
	q.addImport()
	m := q.Struct.NewMethod(name)
	m.Returns(gothicgo.Ret(gothicgo.ErrorType))
	m.Body = func() (string, error) {
		return q.ExecuteTemplate(name)
	}
	return m
}

// GenericFunction takes a template and wraps it in a function that returns
// either a single instance or a slice of the struct.
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

// DefineTable returns the create statement for the table
func (q *QueryBuilder) DefineTable() string {
	var rows []string

	if q.Primary() != "" {
		rows = append(rows,
			fmt.Sprintf("\"%s\" %s", q.Primary(), q.PrimarySQLType()),
			fmt.Sprintf("PRIMARY KEY(\"%s\")", q.Primary()),
		)
	}

	for _, field := range q.fields {
		f, _ := q.Field(field)
		rows = append(rows, fmt.Sprintf("\"%s\" %s", field, f.SQLType))
	}
	return strings.Join(rows, ",\n\t\t\t")
}

// QueryBuilderField wraps the gomodel Field and includes the SQL type
type QueryBuilderField struct {
	gomodel.Field
	SQLType string
}

// Field gets a field by name
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
		SQLType: st,
	}, true
}

// FieldConverter ... I don't remember how this works.
type FieldConverter struct {
	q *QueryBuilder
	*QueryBuilderField
	*Converter
}

// FromDB converts a value from the database and returns a value of the type
// used in the struct
func (f *FieldConverter) FromDB() string {
	return f.fromDB.Call(f.q.File(), f.Name())
}

// ToDB converts a value from the type used in the struct to a value that can be
// passed to the database.
func (f *FieldConverter) ToDB() string {
	return f.toDB.Call(f.q.File(), f.q.Receiver()+"."+f.Name())
}

// Receiver returns the variable used as the reciver on the struct
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

// FieldConverters returns FieldConverters
func (q *QueryBuilder) FieldConverters() []FieldConverter {
	q.populateConverters()
	return q.fieldConverters
}

// ScanFields knows what you think
func (q *QueryBuilder) ScanFields() string {
	q.populateConverters()
	return q.scanFields
}

// Scanner returns the function that scans a SQL row into an instance of the
// struct
func (q *QueryBuilder) Scanner() string {
	scanner := q.SQL.Scanner()
	return scanner.Name()
}
