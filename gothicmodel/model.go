package gothicmodel

type Model struct {
	name       string
	fields     map[string]string
	fieldOrder []string
	primary    string
}

func New(name string) *Model {
	return &Model{
		name:   name,
		fields: make(map[string]string),
	}
}

func (m *Model) AddField(name, kind string) *Model {
	m.fields[name] = kind
	m.fieldOrder = append(m.fieldOrder, name)
	return m
}

func (m *Model) AddPrimary(name, kind string) *Model {
	m.primary = name
	return m.AddField(name, kind)
}

func (m *Model) Field(name string) (string, bool) {
	k, b := m.fields[name]
	return k, b
}

func (m *Model) Fields() []string {
	return m.fieldOrder
}

func (m *Model) Name() string {
	return m.name
}

func (m *Model) Primary() string {
	return m.primary
}
