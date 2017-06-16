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
		AddField("Age", "int")

	pkg := gothicgo.NewPackage("test")
	gm := gomodel.Struct(pkg, m)

	return New(gm)
}

func TestInsert(t *testing.T) {
	sql := setup()

	insert := sql.Insert()
	insert.SetName("insertTest")
	s := insert.String()
	assert.Contains(t, s, "INSERT INTO `test` (`Name`, `Age`) VALUES (?, ?)")
	assert.Contains(t, s, "func (t *test) insertTest() error")

	insert = sql.Insert("Name")
	insert.SetName("insertName")
	s = insert.String()
	assert.Contains(t, s, "(\"INSERT INTO `test` (`Name`) VALUES (?)\", t.Name)")
	assert.Contains(t, s, "func (t *test) insertName() error")
}

func TestUpdate(t *testing.T) {
	sql := setup()

	insert := sql.Update()
	insert.SetName("insertTest")
	s := insert.String()
	t.Error(s)
}
