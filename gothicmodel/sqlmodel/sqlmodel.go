package sqlmodel

import (
	"fmt"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicio"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
	"io"
	"strings"
)

// DefaultConn is used when creating a new SQL instance
var DefaultConn = gothicgo.NewPackageVarRef(gothicgo.MustPackageRef("github.com/adamcolton/buttress/gsql"), "Conn", nil)

// SQL represents a GoModel with database methods
type SQL struct {
	*gomodel.GoModel
	Conn        gothicgo.PackageVarRef
	TableName   string
	Migration   string
	IDQuote     string
	LongIDQuote string
	scanner     *gothicgo.Func
	slct        *gothicgo.Func
	Indexes     []Index
}

type Index struct {
	Type string
	Cols []string
	Name string
}

func (i Index) String() string {
	strs := make([]string, 0, 10)
	if i.Type == "" {
		strs = append(strs, "INDEX")
	} else {
		strs = append(strs, i.Type)
	}
	if i.Name != "" {
		strs = append(strs, fmt.Sprintf(`"%s"`, i.Name))
	}
	strs = append(strs, fmt.Sprintf(`("%s")`, strings.Join(i.Cols, `", "`)))
	return strings.Join(strs, " ")
}

// New takes a GoModel and returns a wrapper for generating SQL code on that
// struct.
func New(goModel *gomodel.GoModel) *SQL {
	return &SQL{
		GoModel:     goModel,
		Conn:        DefaultConn,
		TableName:   goModel.GothicModel.Name(),
		IDQuote:     "`",
		LongIDQuote: `"`,
	}
}

// QueryBuilder returns a QueryBuilder that uses the specified fields
func (s *SQL) QueryBuilder(fields ...string) *QueryBuilder {
	return s.queryBuilder(fields, false)
}

// QueryBuilderAll returns a QueryBuilder that uses all fields, including the
// primary field
func (s *SQL) QueryBuilderAll() *QueryBuilder {
	return s.queryBuilder(nil, true)
}

// Receiver returns the ReceiverName from the struct
func (s *SQL) Receiver() string {
	return s.Struct.ReceiverName
}

func (s *SQL) quote(str string) string {
	return s.IDQuote + str + s.IDQuote
}

func (s *SQL) longquote(str string) string {
	return s.LongIDQuote + str + s.LongIDQuote
}

// TableNameQ returns the Table Name surrounded by ID Quotes
func (s *SQL) TableNameQ() string {
	return s.quote(s.TableName)
}

// TableNameQ returns the Table Name surrounded by ID Quotes
func (s *SQL) TableNameLQ() string {
	return s.longquote(s.TableName)
}

func (s *SQL) addImport() {
	if pkg := s.Conn.PackageRef(); pkg != nil {
		s.File().AddRefImports(pkg)
	}
}

// Primary returns the name of the primary field
func (s *SQL) Primary() string {
	return s.GothicModel.Primary().Name()
}

// PrimaryQ returns the name of the primary field surrounded by ID quotes
func (s *SQL) PrimaryQ() string {
	return s.quote(s.GothicModel.Primary().Name())
}

// PrimaryType returns the model type string
func (s *SQL) PrimaryType() string {
	return s.GothicModel.Primary().Type()
}

// PrimaryGoType returns the Go Type as a string
func (s *SQL) PrimaryGoType() string {
	return gomodel.Types[s.PrimaryType()].RelStr(s)
}

// PrimarySQLType returns the SQL type as a string
func (s *SQL) PrimarySQLType() string {
	return Types[s.PrimaryType()]
}

// Insert will generate a SQL Insert method
func (s *SQL) Insert(fields ...string) *gothicgo.Method {
	return s.QueryBuilder(fields...).GenericMethod("insert")
}

// Update will generate a SQL Update method
func (s *SQL) Update(fields ...string) *gothicgo.Method {
	return s.QueryBuilder(fields...).GenericMethod("update")
}

// Delete will generate a SQL Delete method
func (s *SQL) Delete(fields ...string) *gothicgo.Method {
	m := s.QueryBuilder(fields...).GenericMethod("delete")
	m.UnnamedReturns(gothicgo.Int64Type, gothicgo.ErrorType)
	return m
}

// MigrationFile is the file to write the migraions to.
var MigrationFile *gothicgo.File

func getMigrationFile() *gothicgo.File {
	if MigrationFile == nil {
		pkg, _ := gothicgo.NewPackage("db")
		MigrationFile = pkg.File("migrations.gen")
	}
	MigrationFile.AddRefImports(gothicgo.MustPackageRef("github.com/adamcolton/buttress/gsql"))
	return MigrationFile
}

var initfn *gothicgo.Func

func getInit() *gothicgo.Func {
	if initfn != nil {
		return initfn
	}
	file := getMigrationFile()
	initfn = file.NewFunc("init")
	return initfn
}

// Create will add a migration to create the table
func (s *SQL) Create(migration string, fields ...string) io.WriterTo {
	s.addImport()
	initfn := getInit()

	q := s.QueryBuilder(fields...)
	q.Migration = migration
	wt := q.TemplateWriterTo("createTable")
	initfn.Body = gothicio.WriterToMerge(initfn.Body, wt)
	return wt
}

// GSQL is the sql package to use
var GSQL = gothicgo.MustPackageRef("github.com/adamcolton/buttress/gsql")

// Scanner defines the scanner interface for a row
var Scanner = gothicgo.DefInterface(GSQL, "RowScanner")

// Scanner will add a function to scan a SQL row into an instance of the struct
func (s *SQL) Scanner() *gothicgo.Func {
	if s.scanner != nil {
		return s.scanner
	}
	scanner := s.QueryBuilderAll().GenericFunction("scan", false, false)

	s.File().AddRefImports(Scanner.PackageRef())
	scanner.Sig.Args = append(scanner.Sig.Args, gothicgo.Arg("row", Scanner))
	s.scanner = scanner
	return s.scanner
}

// Select returns a function to select a slice of instances from the SQL
// database
func (s *SQL) Select() *gothicgo.Func {
	if s.slct != nil {
		return s.slct
	}
	if s.scanner == nil {
		s.Scanner()
	}
	s.slct = s.QueryBuilderAll().GenericFunction("select", true, false)
	s.slct.Sig.Args = append(s.slct.Sig.Args, gothicgo.Arg("where", gothicgo.StringType), gothicgo.Arg("args", gothicgo.EmptyInterfaceType))
	s.slct.Variadic = true
	return s.slct
}

// SelectSingle builds func that selects a single instance
func (s *SQL) SelectSingle() *gothicgo.Func {
	if s.slct == nil {
		s.Select()
	}
	fn := s.QueryBuilderAll().GenericFunction("selectsingle", false, false)
	fn.Sig.Args = append(fn.Sig.Args, gothicgo.Arg("where", gothicgo.StringType), gothicgo.Arg("args", gothicgo.EmptyInterfaceType))
	fn.Variadic = true
	return fn
}

// Upsert creates a method that will either insert or update based on the value
// of the primary field.
func (s *SQL) Upsert() *gothicgo.Method {
	return s.QueryBuilder().GenericMethod("upsert")
}

// WhereEqual will generate a method to select either the first or all rows that
// equal the input values
func (s *SQL) WhereEqual(name string, returnSlice bool, fieldArgs []FieldArg) *gothicgo.Func {
	templateName := "whereEqualSingle"
	if returnSlice {
		templateName = "whereEqual"
	}
	return s.whereEqual(name, templateName, returnSlice, false, fieldArgs)
}

func (s *SQL) MustWhereEqual(name string, returnSlice bool, fieldArgs []FieldArg) *gothicgo.Func {
	templateName := "mustWhereEqualSingle"
	if returnSlice {
		templateName = "mustWhereEqual"
	}
	return s.whereEqual(name, templateName, returnSlice, true, fieldArgs)
}

func (s *SQL) whereEqual(name, templateName string, returnSlice, must bool, fieldArgs []FieldArg) *gothicgo.Func {
	fields := make([]string, len(fieldArgs))
	args := make([]gothicgo.NameType, len(fieldArgs))
	callArgs := make([]string, len(fieldArgs))
	for i, f := range fieldArgs {
		fields[i] = f.Field
		callArgs[i] = f.Arg
		args[i].N = f.Arg
		gf, ok := s.GoModel.Field(f.Field)
		if !ok {
			// TODO: return error
			return nil
		}
		args[i].T = gf.GoType()
	}
	qb := s.QueryBuilder(fields...)
	qb.Custom = strings.Join(callArgs, ", ")
	f := qb.GenericFunction(templateName, returnSlice, must)
	f.Sig.Args = args
	f.Sig.Name = name
	return f
}

// FieldArg provides a convenient way to define a mapping of args to fields.
type FieldArg struct {
	Arg, Field string
}

// ConstTableNameString writes the table name a constant string value
func (s *SQL) ConstTableNameString(varName string) io.WriterTo {
	str := fmt.Sprintf(`const %s = "%s"`, varName, s.TableNameQ())
	swt := gothicio.StringWriterTo(str)
	s.File().AddWriterTo(swt)
	return swt
}

func (s *SQL) ConstFieldsString(varName string, fields ...string) io.WriterTo {
	var qb *QueryBuilder
	if fields != nil {
		qb = s.QueryBuilder(fields...)
	} else {
		qb = s.QueryBuilderAll()
	}
	fields = make([]string, len(qb.fields))
	tnq := qb.TableNameQ()
	for i, f := range qb.fields {
		fields[i] = fmt.Sprintf("%s.%s", tnq, qb.quote(f))
	}
	fieldsStr := strings.Join(fields, ", ")
	fieldsStr = strings.Replace(fieldsStr, "\"", "\\\"", -1)

	str := fmt.Sprintf(`const %s = "%s"`, varName, fieldsStr)
	swt := gothicio.StringWriterTo(str)
	s.File().AddWriterTo(swt)
	return swt
}
