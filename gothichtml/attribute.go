package gothichtml

import (
	"sort"
	"strings"
)

type attributes map[string]string

func newAttributes(attrs []string) attributes {
	a := make(attributes)
	a.AddAttributes(attrs...)
	return a
}

func (a attributes) AddAttributes(attrs ...string) {
	for i := 0; i < len(attrs); i += 2 {
		key := strings.ToLower(attrs[i])
		if i+1 == len(attrs) {
			a[key] = ""
		} else {
			a[key] = attrs[i+1]
		}
	}
}

func (a attributes) Attributes() []string {
	keys := make([]string, len(a))
	i := 0
	for k := range a {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func (a attributes) Attribute(key string) (string, bool) {
	v, ok := a[key]
	return v, ok
}

func (a attributes) write(w *writer) {
	for _, k := range a.Attributes() {
		v := a[k]
		w.write(" ")
		w.write(k)
		if v != "" {
			w.write(`="`)
			w.write(v)
			w.write(`"`)
		}
	}
}

func (a attributes) Remove(key string) {
	key = strings.ToLower(key)
	delete(a, key)
}
