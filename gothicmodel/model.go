package gothicmodel

import (
	"fmt"
)

// Field represents a field on a model
type Field struct {
	name    string
	kind    string
	primary bool
	meta    map[string]string
}

// Name of the field
func (f Field) Name() string { return f.name }

// Type of the field
func (f Field) Type() string { return f.kind }

// Primary returns true if the field is the primary field
func (f Field) Primary() bool { return f.primary }

func (f Field) Meta(key string) (string, bool) {
	v, ok := f.meta[key]
	return v, ok
}

func (f Field) AddMeta(key, val string) Field {
	f.meta[key] = val
	return f
}

// Metas returns a slice of the keys. A new slice is generated each time.
func (f Field) Metas() []string {
	ms := make([]string, len(f.meta))
	i := 0
	for k := range f.meta {
		ms[i] = k
		i++
	}
	return ms
}

// GothicModel is a generalization of a data structure. It has a name and some number
// of fields one of which is the primary field
type GothicModel struct {
	name      string
	fieldsMap map[string]Field
	fields    []Field
	primary   Field
	meta      map[string]string
}

// Fields is used to define the fields on a model.
type Fields [][2]string

// New model
func New(name string, fields Fields) (*GothicModel, error) {
	m := &GothicModel{
		name:      name,
		fields:    make([]Field, 0, len(fields)),
		fieldsMap: make(map[string]Field, len(fields)),
		meta:      make(map[string]string),
	}

	for _, field := range fields {
		if _, err := m.AddField(field[0], field[1]); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// Must creates a new model and panics if there is an error
func Must(name string, fields Fields) *GothicModel {
	m, err := New(name, fields)
	if err != nil {
		panic(err)
	}
	return m
}

// AddFields to a model
func (m *GothicModel) AddFields(fields Fields) error {
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
func (m *GothicModel) AddField(name, kind string) (Field, error) {
	if err := m.validate(name, kind); err != nil {
		return Field{}, err
	}
	return m.addField(name, kind), nil
}

func (m *GothicModel) addField(name, kind string) Field {
	f := Field{
		name: name,
		kind: kind,
		meta: make(map[string]string),
	}

	if len(m.fields) == 0 {
		f.primary = true
		m.primary = f
	}
	m.fields = append(m.fields, f)
	m.fieldsMap[name] = f
	return f
}

func (m *GothicModel) validate(name, kind string) error {
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
func (m *GothicModel) AddPrimary(name, kind string) (Field, error) {
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
func (m *GothicModel) Field(name string) (Field, bool) {
	f, b := m.fieldsMap[name]
	return f, b
}

// MustField gets a field by name and will panic if the field does not exist
func (m *GothicModel) MustField(name string) Field {
	f, b := m.fieldsMap[name]
	if !b {
		panic(fmt.Sprintf(`MustField) "%s" is not defined on model %s`, name, m.name))
	}
	return f
}

// Fields returns the fiels on the model. If no field names are given, all
// fields are returned. If a list of names is given, only those fields are
// returned.
func (m *GothicModel) Fields(names ...string) []Field {
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
func (m *GothicModel) Name() string {
	return m.name
}

// Primary field on the model
func (m *GothicModel) Primary() Field {
	return m.primary
}

// SkipFields returns all fields excluding those defined to skip
func (m *GothicModel) SkipFields(skip ...string) []Field {
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

func (m *GothicModel) Meta(key string) (string, bool) {
	v, ok := m.meta[key]
	return v, ok
}

func (m *GothicModel) AddMeta(key, val string) *GothicModel {
	m.meta[key] = val
	return m
}

// Metas returns a slice of the keys. A new slice is generated each time.
func (m *GothicModel) Metas() []string {
	ms := make([]string, len(m.meta))
	i := 0
	for k := range m.meta {
		ms[i] = k
		i++
	}
	return ms
}
