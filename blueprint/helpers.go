package blueprint

import (
	"fmt"
)

func PadStr(str string, ln int) string {
	for len(str) < ln {
		str += " "
	}
	return str
}

type StrBuilder struct {
	str string
}

func StringBuilder(format string, a ...interface{}) *StrBuilder {
	sb := &StrBuilder{}
	sb.Add(format, a...)
	return sb
}

func StringBuilderLn(format string, a ...interface{}) *StrBuilder {
	sb := &StrBuilder{}
	sb.AddLn(format, a...)
	return sb
}

func (sb *StrBuilder) Add(format string, a ...interface{}) {
	sb.str += fmt.Sprintf(format, a...)
}

func (sb *StrBuilder) AddLn(format string, a ...interface{}) {
	sb.str += fmt.Sprintf(format, a...) + "\n"
}

func (sb *StrBuilder) String() string {
	return sb.str
}

type StringGeneratorFacade string

func (s StringGeneratorFacade) Prepare()       {}
func (s StringGeneratorFacade) Export() string { return string(s) }
