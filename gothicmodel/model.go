package gothicmodel

import (
	"fmt"
)

// Field represents a field on a model
type Field struct {
	name    string
	kind    string
	primary bool
}

// Name of the field
func (f Field) Name() string { return f.name }

// Type of the field
func (f Field) Type() string { return f.kind }

// Primary returns true if the field is the primary field
func (f Field) Primary() bool { return f.primary }

// Model is a generalization of a data structure. It has a name and some number
// of fields one of which is the primary field
type Model struct {
	name      string
	fieldsMap map[string]Field
	fields    []Field
	primary   Field
}

// Fields is used to define the fields on a model.
type Fields [][2]string

// New model
func New(name string, fields Fields) (*Model, error) {
	m := &Model{
		name:      name,
		fields:    make([]Field, 0, len(fields)),
		fieldsMap: make(map[string]Field, len(fields)),
	}

	for _, field := range fields {
		if _, err := m.AddField(field[0], field[1]); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// Must creates a new model and panics if there is an error
func Must(name string, fields Fields) *Model {
	m, err := New(name, fields)
	if err != nil {
		panic(err)
	}
	return m
}

// AddFields to a model
func (m *Model) AddFields(fields Fields) error {
	names := make(map[string]struct{})
	for _, field := range fields {
		if err := m.validate(field[0], field[1]); err != nil {
			return err
		}
		if _, defined := names[field[0]]; defined {
			return fmt.Errorf("%s) field name repeated", field[0])
		}
		names[field[0]] = struct{}{}
	}
	for _, field := range fields {
		m.addField(field[0], field[1])
	}
	return nil
}

// AddField by name and kind
func (m *Model) AddField(name, kind string) (Field, error) {
	if err := m.validate(name, kind); err != nil {
		return Field{}, err
	}
	return m.addField(name, kind), nil
}

func (m *Model) addField(name, kind string) Field {
	f := Field{
		name: name,
		kind: kind,
	}

	if len(m.fields) == 0 {
		f.primary = true
		m.primary = f
	}
	m.fields = append(m.fields, f)
	m.fieldsMap[name] = f
	return f
}

func (m *Model) validate(name, kind string) error {
	if name == "" {
		return fmt.Errorf("Name cannot be empty")
	}
	if kind == "" {
		return fmt.Errorf("%s) Type cannot be empty", name)
	}
	if _, defined := m.fieldsMap[name]; defined {
		return fmt.Errorf("%s) field name repeated", name)
	}
	return nil
}

// AddPrimary changes which field is primary
func (m *Model) AddPrimary(name, kind string) (Field, error) {
	if len(m.fields) == 0 {
		return m.AddField(name, kind)
	}

	var f Field
	if err := m.validate(name, kind); err != nil {
		return f, err
	}

	f.name = name
	f.kind = kind
	f.primary = true

	old := m.primary
	old.primary = false
	m.fieldsMap[old.name] = old
	m.fields[0] = old

	m.primary = f
	m.fieldsMap[f.name] = f
	m.fields = append([]Field{f}, m.fields...)
	return f, nil
}

// Field gets a field by name
func (m *Model) Field(name string) (Field, bool) {
	k, b := m.fieldsMap[name]
	return k, b
}

// Fields returns the fiels on the model. If no field names are given, all
// fields are returned. If a list of names is given, only those fields are
// returned.
func (m *Model) Fields(names ...string) []Field {
	if names == nil {
		cp := make([]Field, len(m.fields))
		copy(cp, m.fields)
		return cp
	}
	var fields []Field
	for _, name := range names {
		if field, ok := m.fieldsMap[name]; ok {
			fields = append(fields, field)
		}
	}
	return fields
}

// Name of the model
func (m *Model) Name() string {
	return m.name
}

// Primary field on the model
func (m *Model) Primary() Field {
	return m.primary
}

// SkipFields returns all fields excluding those defined to skip
func (m *Model) SkipFields(skip ...string) []Field {
	sm := make(map[string]bool, len(skip))
	for _, s := range skip {
		sm[s] = true
	}
	var fields []Field
	for _, f := range m.fields {
		if !sm[f.name] {
			fields = append(fields, f)
		}
	}
	return fields
}
