package sqlmodel

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
)

var DefaultConn = gothicgo.NewPackageVarRef(gothicgo.MustPackageRef("github.com/adamcolton/buttress/gsql"), "Conn", nil)

type SQL struct {
	*gomodel.GoModel
	Conn      gothicgo.PackageVarRef
	TableName string
	Migration string
	IDQuote   string
	scanner   *gothicgo.Func
}

func New(model *gomodel.GoModel) *SQL {
	return &SQL{
		GoModel:   model,
		Conn:      DefaultConn,
		TableName: model.Model.Name(),
		IDQuote:   "`",
	}
}

func (s *SQL) QueryBuilder(fields ...string) *QueryBuilder {
	return s.queryBuilder(fields, false)
}

func (s *SQL) QueryBuilderAll() *QueryBuilder {
	return s.queryBuilder(nil, true)
}

func (s *SQL) Receiver() string {
	return s.Struct.ReceiverName
}

func (s *SQL) quote(str string) string {
	return s.IDQuote + str + s.IDQuote
}

func (s *SQL) TableNameQ() string {
	return s.quote(s.TableName)
}

func (s *SQL) addImport() {
	if pkg := s.Conn.PackageRef(); pkg != nil {
		s.File().AddRefImports(pkg)
	}
}

func (s *SQL) Primary() string {
	return s.Model.Primary().Name()
}

func (s *SQL) PrimaryQ() string {
	return s.quote(s.Model.Primary().Name())
}

func (s *SQL) PrimaryType() string {
	return s.Model.Primary().Type()
}

func (s *SQL) PrimaryGoType() string {
	return gomodel.Types[s.PrimaryType()].RelStr(s)
}

func (s *SQL) PrimarySqlType() string {
	return Types[s.PrimaryType()]
}

func (s *SQL) Insert(fields ...string) *gothicgo.Method {
	return s.QueryBuilder(fields...).GenericMethod("insert")
}

func (s *SQL) Update(fields ...string) *gothicgo.Method {
	return s.QueryBuilder(fields...).GenericMethod("update")
}

var MigrationFile *gothicgo.File

func getMigrationFile() *gothicgo.File {
	if MigrationFile == nil {
		pkg, _ := gothicgo.NewPackage("db")
		MigrationFile = pkg.File("migrations.gen")
	}
	MigrationFile.AddRefImports(gothicgo.MustPackageRef("github.com/adamcolton/buttress/gsql"))
	return MigrationFile
}

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

var GoSql = gothicgo.MustPackageRef("database/sql")
var Scanner = gothicgo.DefInterface(GoSql, "Scanner")

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

func (s *SQL) Select() *gothicgo.Func {
	f := s.QueryBuilderAll().GenericFunction("select", true)
	f.Args = append(f.Args, gothicgo.Arg("where", gothicgo.StringType), gothicgo.Arg("args", gothicgo.EmptyInterfaceType))
	f.Variadic = true
	return f
}

func (s *SQL) Upsert() *gothicgo.Method {
	return s.QueryBuilder().GenericMethod("upsert")
}
