package entityBP

import (
	"fmt"
)

type unmarshalEntityGenerator struct {
	ent *Entity
}

func (ent *Entity) addUnmarshalEntityGenerator() *Entity {
	return ent.AddGenerator(unmarshalEntityGenerator{
		ent: ent,
	})
}

const unmarshalEntityTemplate = `func Unmarshal%sEntity(b *[]byte) entity.Entity{
  return entity.Entity(Unmarshal%s(b))
}
`

func (ug unmarshalEntityGenerator) Export() string {
	return fmt.Sprintf(unmarshalEntityTemplate, ug.ent.Name(), ug.ent.Name())
}

func (ug unmarshalEntityGenerator) Prepare() {}
