package sqlmodel

import (
	"bytes"
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() *SQL {
	m, err := gothicmodel.New("test", gothicmodel.Fields{
		{"ID", "uint"},
		{"Name", "string"},
		{"Age", "int"},
		{"LastLogin", "datetime"},
	})
	if err != nil {
		panic(err)
	}

	pkg, err := gothicgo.NewPackage("test")
	if err != nil {
		panic(err)
	}
	gm, err := gomodel.New(pkg, m)
	if err != nil {
		panic(err)
	}

	return New(gm)
}

func TestQueryBuilder(t *testing.T) {
	sql := setup()
	q := sql.QueryBuilder()

	assert.Equal(t, "t", q.Receiver())
	assert.Equal(t, "test", q.TableName)
	assert.Equal(t, "`test`", q.TableNameQ())
	assert.Equal(t, "ID", q.Primary())
	assert.Equal(t, "`ID`", q.PrimaryQ())
	assert.Equal(t, "uint", q.PrimaryType())
	assert.Equal(t, gothicgo.UintType.String(), q.PrimaryGoType())
	assert.Equal(t, Types["uint"], q.PrimarySQLType())
	assert.Equal(t, "Name, Age, LastLogin", q.Fields())
	assert.Equal(t, "`Name`, `Age`, `LastLogin`", q.FieldsQ())
	assert.Equal(t, "?, ?, ?", q.QM())
	assert.Equal(t, "`Name`=?, `Age`=?, `LastLogin`=?", q.Set())
	assert.Equal(t, "0", q.PrimaryZeroVal())
}

func TestInsert(t *testing.T) {
	sql := setup()

	insert := sql.Insert()
	insert.SetName("insertTest")
	s := insert.String()
	assert.Contains(t, s, "res, err := gsql.Conn.Exec(\"INSERT INTO `test` (`Name`, `Age`, `LastLogin`) VALUES (?, ?, ?)\", t.Name, t.Age, gsql.TimeToString(t.LastLogin))")
	assert.Contains(t, s, "func (t *test) insertTest() error")

	insert = sql.Insert("Name")
	insert.SetName("insertName")
	s = insert.String()
	assert.Contains(t, s, "(\"INSERT INTO `test` (`Name`) VALUES (?)\", t.Name)")
	assert.Contains(t, s, "func (t *test) insertName() error")
}

func TestUpdate(t *testing.T) {
	sql := setup()

	update := sql.Update()
	update.SetName("updateTest")
	s := update.String()
	assert.Contains(t, s, "_, err := gsql.Conn.Exec(\"UPDATE `test` SET `Name`=?, `Age`=?, `LastLogin`=? WHERE `ID`=?\", t.Name, t.Age, gsql.TimeToString(t.LastLogin), t.ID)")
}

func TestCreate(t *testing.T) {
	sql := setup()
	sql.Indexes = append(sql.Indexes,
		Index{
			Cols: []string{"Name", "Age"},
		},
		Index{
			Type: "UNIQUE",
			Name: "UniqueNames",
			Cols: []string{"Name"},
		},
	)

	var buf bytes.Buffer
	_, err := sql.Create("12345_create_user").WriteTo(&buf)
	s := buf.String()

	assert.NoError(t, err)
	assert.Contains(t, s, "\"ID\" int UNSIGNED DEFAULT 0 NOT NULL,")
	assert.Contains(t, s, "CREATE TABLE IF NOT EXISTS \"test\" (")
	assert.Contains(t, s, `INDEX ("Name", "Age"),`)
	assert.Contains(t, s, `UNIQUE "UniqueNames" ("Name")`)
}

func TestScan(t *testing.T) {
	sql := setup()

	scan := sql.Scanner()
	s := scan.String()
	assert.Contains(t, s, "func scantest(row gsql.RowScanner) (*test, error) {")
	assert.Contains(t, s, "err := row.Scan(&(t.ID), &(t.Name), &(t.Age), &(LastLogin))")
	assert.Contains(t, s, "var LastLogin string")
	assert.Contains(t, s, "t.LastLogin = gsql.StringToTime(LastLogin)")

	s = sql.File().Imports.String()
	assert.Contains(t, s, `"github.com/adamcolton/buttress/gsql"`)
}

func TestSelect(t *testing.T) {
	sql := setup()

	slct := sql.Select()
	s := slct.String()
	assert.Contains(t, s, "rows, err := gsql.Conn.Query(\"SELECT `ID`, `Name`, `Age`, `LastLogin` FROM `test` \"+where, args...)")
}

func TestUpsert(t *testing.T) {
	sql := setup()

	slct := sql.Upsert()
	s := slct.String()
	assert.Contains(t, s, "if t.ID != 0 {")
	assert.Contains(t, s, "_, err := gsql.Conn.Exec(\"UPDATE `test` SET `Name`=?, `Age`=?, `LastLogin`=? WHERE `ID`=?\", t.Name, t.Age, gsql.TimeToString(t.LastLogin), t.ID)")
	assert.Contains(t, s, "res, err := gsql.Conn.Exec(\"INSERT INTO `test` (`Name`, `Age`, `LastLogin`) VALUES (?, ?, ?)\", t.Name, t.Age, gsql.TimeToString(t.LastLogin))")
}

func TestSelectSingle(t *testing.T) {
	sql := setup()

	single := sql.SelectSingle()
	s := single.String()
	assert.Contains(t, s, `t, err := selecttest(where+" LIMIT 1", args...)`)
	assert.Contains(t, s, `func selectsingletest(where string, args ...interface{}) (*test, error) {`)
	assert.Contains(t, s, `return t[0], nil`)
}

func TestDelete(t *testing.T) {
	sql := setup()

	del := sql.Delete("ID", "Name")
	s := del.String()
	assert.Contains(t, s, "res, err := gsql.Conn.Exec(\"DELETE FROM `test` WHERE `ID`=? AND `Name`=?\", t.ID, t.Name)")
	assert.Contains(t, s, `return res.RowsAffected()`)
}

func TestWhereEqual(t *testing.T) {
	sql := setup()

	getByID := sql.WhereEqual("GetByID", false, []FieldArg{
		{"id", "ID"},
	})
	s := getByID.String()
	assert.Contains(t, s, "func GetByID(id uint) (*test, error) {")
	assert.Contains(t, s, "t, err := selecttest(\"WHERE `ID`=? LIMIT 1\", id)")
	assert.Contains(t, s, "return t[0], nil")

	getByAge := sql.WhereEqual("GetByNameAge", true, []FieldArg{
		{"name", "Name"},
		{"age", "Age"},
	})
	s = getByAge.String()
	assert.Contains(t, s, "func GetByNameAge(name string, age int) ([]*test, error) {")
	assert.Contains(t, s, "return selecttest(\"WHERE `Name`=? AND `Age`=?\", name, age)")

	mustGetByID := sql.MustWhereEqual("MustGetByID", false, []FieldArg{
		{"id", "ID"},
	})
	s = mustGetByID.String()
	assert.Contains(t, s, "func MustGetByID(id uint) *test {")
	assert.Contains(t, s, "t, err := selecttest(\"WHERE `ID`=? LIMIT 1\", id)")
	assert.Contains(t, s, "return t[0]")

	mustGetByAge := sql.MustWhereEqual("MustGetByNameAge", true, []FieldArg{
		{"name", "Name"},
		{"age", "Age"},
	})
	s = mustGetByAge.String()
	assert.Contains(t, s, "func MustGetByNameAge(name string, age int) []*test {")
	assert.Contains(t, s, "t, err := selecttest(\"WHERE `Name`=? AND `Age`=?\", name, age)")
}

func TestConstVars(t *testing.T) {
	sql := setup()

	wt := sql.ConstTableNameString("testTable")
	buf := &bytes.Buffer{}
	wt.WriteTo(buf)
	expected := "const testTable = \"`test`\""
	assert.Equal(t, expected, buf.String())

	buf.Reset()
	wt = sql.ConstFieldsString("testFields")
	wt.WriteTo(buf)
	expected = "const testFields = \"`test`.`ID`, `test`.`Name`, `test`.`Age`, `test`.`LastLogin`\""
	assert.Equal(t, expected, buf.String())
}
