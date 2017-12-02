package sqlmodel

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicio"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
	"io"
)

// DefaultConn is used when creating a new SQL instance
var DefaultConn = gothicgo.NewPackageVarRef(gothicgo.MustPackageRef("github.com/adamcolton/buttress/gsql"), "Conn", nil)

// SQL represents a GoModel with database methods
type SQL struct {
	*gomodel.GoModel
	Conn      gothicgo.PackageVarRef
	TableName string
	Migration string
	IDQuote   string
	scanner   *gothicgo.Func
	slct      *gothicgo.Func
	Index     map[string]string
}

// New takes a GoModel and returns a wrapper for generating SQL code on that
// struct.
func New(model *gomodel.GoModel) *SQL {
	return &SQL{
		GoModel:   model,
		Conn:      DefaultConn,
		TableName: model.Model.Name(),
		IDQuote:   "`",
		Index:     make(map[string]string),
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

// TableNameQ returns the Table Name surrounded by ID Quotes
func (s *SQL) TableNameQ() string {
	return s.quote(s.TableName)
}

func (s *SQL) addImport() {
	if pkg := s.Conn.PackageRef(); pkg != nil {
		s.File().AddRefImports(pkg)
	}
}

// Primary returns the name of the primary field
func (s *SQL) Primary() string {
	return s.Model.Primary().Name()
}

// PrimaryQ returns the name of the primary field surrounded by ID quotes
func (s *SQL) PrimaryQ() string {
	return s.quote(s.Model.Primary().Name())
}

// PrimaryType returns the model type string
func (s *SQL) PrimaryType() string {
	return s.Model.Primary().Type()
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
	scanner := s.QueryBuilderAll().GenericFunction("scan", false)

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
	s.slct = s.QueryBuilderAll().GenericFunction("select", true)
	s.slct.Sig.Args = append(s.slct.Sig.Args, gothicgo.Arg("where", gothicgo.StringType), gothicgo.Arg("args", gothicgo.EmptyInterfaceType))
	s.slct.Variadic = true
	return s.slct
}

// SelectSingle builds func that selects a single instance
func (s *SQL) SelectSingle() *gothicgo.Func {
	if s.slct == nil {
		s.Select()
	}
	fn := s.QueryBuilderAll().GenericFunction("selectsingle", false)
	fn.Sig.Args = append(fn.Sig.Args, gothicgo.Arg("where", gothicgo.StringType), gothicgo.Arg("args", gothicgo.EmptyInterfaceType))
	fn.Variadic = true
	return fn
}

// Upsert creates a method that will either insert or update based on the value
// of the primary field.
func (s *SQL) Upsert() *gothicgo.Method {
	return s.QueryBuilder().GenericMethod("upsert")
}
