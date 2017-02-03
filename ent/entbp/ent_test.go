package entbp

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicserial"
	"github.com/adamcolton/sai"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnt(t *testing.T) {
	p := NewPackage("test")
	e := p.Ent("Person")
	e.AddField("Name", gothicgo.StringType)
	e.AddField("Age", gothicgo.IntType)
	assert.Equal(t, "*test.PersonRef", e.Ref().String())

	strct := e.GoStruct()
	assert.Equal(t, "*test.PersonRef", strct.Ref().String())
	strct.File().Package().ImportResolver().Add("serial", "github.com/adamcolton/gothic/serial")

	sd, ok := gothicserial.GetSerializeDef("*test.Person")
	if !ok {
		t.Error("Did not find a serizalize definition")
	}
	expected := "adam.Marshal()"
	got := sd.Marshal("adam", "test")
	assert.Equal(t, expected, got, "Ent Marshal")

	expected = "test.UnmarshalPerson(buf)"
	got = sd.Unmarshal("buf", "foo")
	assert.Equal(t, expected, got, "Ent Unmarshal")

	wc := sai.New()
	f := strct.File()
	f.Writer = wc
	f.Prepare()
	f.Generate()
	expected = "type Person struct {\n\tent.Ent\n\tName string\n\tAge  int\n}"
	assert.Equal(t, expected, strct.String(), "Ent Struct")

	expectedStrs := []string{
		"type Person struct",
		"func (o *Person) Marshal() []byte",
		"func UnmarshalPerson(b *[]byte) *Person",
		"type PersonRef struct",
		"func GetPerson(id []byte) *Person",
		"func UnmarshalPersonEntity(b *[]byte) ent.Entity",
		"func (r *PersonRef) Get() *Person",
		"func UnmarshalPersonRef(b *[]byte) *PersonRef",
		"func (o *Person) Save()",
		"func (o *Person) Ref() *PersonRef",
	}
	got = wc.String()
	for _, str := range expectedStrs {
		assert.Contains(t, got, str)
	}
}
