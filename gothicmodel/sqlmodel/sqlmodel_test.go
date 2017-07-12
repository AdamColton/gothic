package sqlmodel

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel"
	"github.com/adamcolton/gothic/gothicmodel/gomodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() *SQL {
	m := gothicmodel.New("test").
		AddPrimary("ID", "uint").
		AddField("Name", "string").
		AddField("Age", "int").
		AddField("LastLogin", "datetime")

	pkg := gothicgo.NewPackage("test")
	gm := gomodel.Struct(pkg, m)

	return New(gm)
}

func TestInsert(t *testing.T) {
	sql := setup()

	insert := sql.Insert()
	insert.SetName("insertTest")
	s := insert.String()
	assert.Contains(t, s, "res, err := conn.Exec(\"INSERT INTO `test` (`Name`, `Age`, `LastLogin`) VALUES (?, ?, ?)\", t.Name, t.Age, TimeToString(t.LastLogin))")
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
	assert.Contains(t, s, "_, err := conn.Exec(\"UPDATE `test` SET (Name=?, Age=?, LastLogin=?) WHERE `ID`=?\", t.Name, t.Age, TimeToString(t.LastLogin), t.ID)")
}

func TestCreate(t *testing.T) {
	sql := setup()

	s := sql.Create("12345_create_user")
	assert.Contains(t, s, "\"ID\" int UNSIGNED DEFAULT 0 NOT NULL,")
	assert.Contains(t, s, "CREATE TABLE IF NOT EXISTS \"test\" (")
}

func TestScan(t *testing.T) {
	sql := setup()

	scan := sql.Scan()
	s := scan.String()
	assert.Contains(t, s, "err := rows.Scan(&(t.ID), &(t.Name), &(t.Age), &(LastLogin))")
	assert.Contains(t, s, "t.LastLogin = StringToTime(LastLogin)")
}

func TestSelect(t *testing.T) {
	sql := setup()

	slct := sql.Select()
	s := slct.String()
	assert.Contains(t, s, "rows, err := conn.Query(\"SELECT `ID`, `Name`, `Age`, `LastLogin` FROM `test` \"+where, args...)")
}

func TestUpsert(t *testing.T) {
	sql := setup()

	slct := sql.Upsert()
	s := slct.String()
	assert.Contains(t, s, "if t.ID != 0 {")
	assert.Contains(t, s, "_, err := conn.Exec(\"UPDATE `test` SET (Name=?, Age=?, LastLogin=?) WHERE `ID`=?\", t.Name, t.Age, TimeToString(t.LastLogin), t.ID)")
	assert.Contains(t, s, "res, err := conn.Exec(\"INSERT INTO `test` (`Name`, `Age`, `LastLogin`) VALUES (?, ?, ?)\", t.Name, t.Age, TimeToString(t.LastLogin))")
}
