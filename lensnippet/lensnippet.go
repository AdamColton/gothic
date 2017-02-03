package lensnippet

import (
	"fmt"
	"github.com/adamcolton/gothic"
	"github.com/adamcolton/gothic/gothicgo"
)

type min struct {
	on NameTypeContainer
	l  int
}

type minInstance struct {
	m    *min
	errs string
	name string
}

func (m *min) New() gothic.SnippetInstance {
	return &minInstance{
		m:    m,
		errs: "errs",
		name: m.on.Name(),
	}
}

func (m *minInstance) Set(key, val string) {
	if key == "name" {
		m.name = val
	} else if key == "errs" {
		m.errs = val
	} else {
		panic("Can only set 'name' and 'errs' on Min Length validator instance")
	}
}

func (m *minInstance) Prepare() {}

var minTmpl = "if len(%s) < %d {\n\t%s.Add(\"%s must be at least %d long\")\n}"

func (m *minInstance) Generate() []string {
	return []string{
		fmt.Sprintf(minTmpl, m.name, m.m.l, m.errs, m.m.on.Name(), m.m.l),
	}
}

type NameTypeContainer interface {
	Name() string
	Type() gothicgo.Type
	gothic.SnippetContainer
}

func Min(on NameTypeContainer, l int) gothic.Snippet {
	m := &min{
		on: on,
		l:  l,
	}
	on.AddSnippetTo("validators", m)
	return m
}
