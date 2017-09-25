package sqlmodel

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
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
}

// New takes a GoModel and returns a wrapper for generating SQL code on that
// struct.
func New(model *gomodel.GoModel) *SQL {
	return &SQL{
		GoModel:   model,
		Conn:      DefaultConn,
		TableName: model.Model.Name(),
		IDQuote:   "`",
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

// Create will add a migration to create the table
func (s *SQL) Create(migration string, fields ...string) (string, error) {
	s.addImport()
	file := getMigrationFile()

	q := s.QueryBuilder(fields...)
	q.Migration = migration
	str, err := q.ExecuteTemplate("createTable")
	if err == nil {
		file.AddCode(str)
	}
	return str, err
}

// GoSQL is the sql package to use
var GoSQL = gothicgo.MustPackageRef("database/sql")

// Scanner defines the scanner interface for a row
var Scanner = gothicgo.DefInterface(GoSQL, "Scanner")

// Scanner will add a function to scan a SQL row into an instance of the struct
func (s *SQL) Scanner() *gothicgo.Func {
	if s.scanner != nil {
		return s.scanner
	}
	scanner := s.QueryBuilderAll().GenericFunction("scan", false)

	s.File().AddRefImports(Scanner.PackageRef())
	scanner.Args = append(scanner.Args, gothicgo.Arg("rows", Scanner))
	s.scanner = scanner
	return s.scanner
}

// Select returns a function to select a slice of instances from the SQL
// database
func (s *SQL) Select() *gothicgo.Func {
	f := s.QueryBuilderAll().GenericFunction("select", true)
	f.Args = append(f.Args, gothicgo.Arg("where", gothicgo.StringType), gothicgo.Arg("args", gothicgo.EmptyInterfaceType))
	f.Variadic = true
	return f
}

// Upsert creates a method that will either insert or update based on the value
// of the primary field.
func (s *SQL) Upsert() *gothicgo.Method {
	return s.QueryBuilder().GenericMethod("upsert")
}
