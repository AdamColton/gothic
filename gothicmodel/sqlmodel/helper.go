package sqlmodel

import (
	"strings"
)

// helper is used to map the model to the templates.
// useFields is changed per call and defines which fields to use in a given
// template
type helper struct {
	fields      []string
	fieldsMap   map[string]bool
	args        map[string]string
	Name        string
	Primary     string
	PrimaryType string
	Receiver    string
	Conn        string
	useFields   []string
}

func createHelper(s *SQL) *helper {
	h := &helper{
		Name:      s.TableName,
		Primary:   s.model.Model.Primary(),
		Receiver:  s.model.Struct.ReceiverName,
		Conn:      s.Conn,
		fieldsMap: make(map[string]bool),
		args:      make(map[string]string),
	}
	for _, f := range s.model.Model.Fields() {
		ts, _ := s.model.Model.Field(f)
		_, sqlOk := Types[ts]
		gt, goOk := s.model.Struct.Field(f)
		if !sqlOk || !goOk {
			continue
		}
		if f == h.Primary {
			h.PrimaryType = gt.Type().RelStr(s.model.Struct.PackageName())
			continue
		}

		arg := h.Receiver + "." + f
		if c, ok := Converters[ts]; ok {
			arg = c.toDB + "(" + arg + ")"
		}
		h.append(f, arg)
	}
	return h
}

func (h *helper) Fields() string {
	return "`" + strings.Join(h.useFields, "`, `") + "`"
}
func (h *helper) QM() string {
	qm := make([]string, len(h.useFields))
	for i := range qm {
		qm[i] = "?"
	}
	return strings.Join(qm, ", ")
}

func (h *helper) Args() string {
	args := make([]string, len(h.useFields))
	for i, field := range h.useFields {
		args[i] = h.args[field]
	}
	return strings.Join(args, ", ")
}

func (h *helper) append(field, arg string) {
	h.fields = append(h.fields, field)
	h.args[field] = arg
	h.fieldsMap[field] = true
}

func (h *helper) QName() string {
	return "`" + h.Name + "`"
}

func (h *helper) Set() string {
	set := make([]string, len(h.useFields))
	for i, field := range h.useFields {
		set[i] = field + "=?"
	}
	return strings.Join(set, ", ")
}

func (h *helper) QPrimary() string {
	return "`" + h.Primary + "`"
}

func (h *helper) PrimaryArg() string {
	return h.Receiver + "." + h.Primary
}
