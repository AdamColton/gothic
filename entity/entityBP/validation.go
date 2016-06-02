package entityBP

import (
	"fmt"
)

type Validator interface {
	Export(string, *Entity) string
	Import(string, *Entity)
}

type MinLen struct {
	Len int
}

var minLenTemplate = `  if len(%s) < %d {
    valErrs.Add("%s", "Did not meet minimum length")
  }
`

func (m MinLen) Export(fieldName string, s *Entity) string {
	return fmt.Sprintf(minLenTemplate, fieldName, m.Len, fieldName)
}

func (m MinLen) Import(fieldName string, s *Entity) {
	s.ImportGothic("validation")
}
