package sqlmodel

import (
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
	gm := gomodel.Struct(pkg, m)

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
	assert.Equal(t, Types["uint"], q.PrimarySqlType())
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
	assert.Contains(t, s, "_, err := gsql.Conn.Exec(\"UPDATE `test` SET (`Name`=?, `Age`=?, `LastLogin`=?) WHERE `ID`=?\", t.Name, t.Age, gsql.TimeToString(t.LastLogin), t.ID)")
}

func TestCreate(t *testing.T) {
	sql := setup()

	s, err := sql.Create("12345_create_user")
	assert.NoError(t, err)
	assert.Contains(t, s, "\"ID\" int UNSIGNED DEFAULT 0 NOT NULL,")
	assert.Contains(t, s, "CREATE TABLE IF NOT EXISTS \"test\" (")
}

func TestScan(t *testing.T) {
	sql := setup()

	scan := sql.Scanner()
	s := scan.String()
	assert.Contains(t, s, "func scantest(rows sql.Scanner) (*test, error) {")
	assert.Contains(t, s, "err := rows.Scan(&(t.ID), &(t.Name), &(t.Age), &(LastLogin))")
	assert.Contains(t, s, "t.LastLogin = gsql.StringToTime(LastLogin)")

	s = sql.File().Imports.String()
	assert.Contains(t, s, `"database/sql"`)
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
	assert.Contains(t, s, "_, err := gsql.Conn.Exec(\"UPDATE `test` SET (`Name`=?, `Age`=?, `LastLogin`=?) WHERE `ID`=?\", t.Name, t.Age, gsql.TimeToString(t.LastLogin), t.ID)")
	assert.Contains(t, s, "res, err := gsql.Conn.Exec(\"INSERT INTO `test` (`Name`, `Age`, `LastLogin`) VALUES (?, ?, ?)\", t.Name, t.Age, gsql.TimeToString(t.LastLogin))")
}
