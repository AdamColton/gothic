package structBP

import (
  "github.com/adamcolton/gothic/blueprint"
)

type Field struct {
  name      string
  typ       blueprint.Type
}

func (f *Field) Name() string { return f.name}
func (f *Field) Type() blueprint.Type { return f.typ}

func (f *Field) ID() string {
  if f.name != ""{
    return f.name
  }
  return f.typ.Name()
}