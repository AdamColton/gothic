package lapizBP

import (
	"github.com/adamcolton/gothic/blueprint"
	"github.com/adamcolton/gothic/entity/entityBP"
)

type LapizClass struct {
	ent *entityBP.Entity
}

func NewLapizClass(ent *entityBP.Entity) *LapizClass {
	lc := &LapizClass{
		ent: ent,
	}
	blueprint.Register(lc)
	return lc
}

func (lc *LapizClass) Prepare() {}

/*
type Person struct {
  entity.Ent
  Name    string
  Age     int
  Numbers []int
}

var Person = Lapiz.Constructor(function(ID, Name, Age, Numbers){
  this.priv.properties({
    ID      : "string",
    Name    : "string",
    Age     : "int",
    Numbers : Lapiz.parse.array("int"),
  }, this.priv.argDict());
});

type Course struct {
  entity.Ent
  Name     string
  Teacher  person.PersonRef
  Students []person.PersonRef
}

var Person = Lapiz.Constructor(function(ID, Name, TeacherID, StudentIDs){
  this.priv.properties({
    ID         : "string",
    Name       : "string",
    TeacherID  : "string",
    Teacher    : Lapiz.parse.relational("string", this, "TeacherID", Person.get),
    StudentIDs : Lapiz.parse.array("string"),
    Students   : Lapiz.parse.array()
  }, this.priv.argDict());
});

*/

func (lc *LapizClass) Export() {

}
