package gothic

type Snippet interface {
	AddContext(key, value string) Snippet
	String() string
}

type SnippetContainer struct {
	subContainers []*SnippetContainer
	snippets      []Snippet
}

func NewContainer() *SnippetContainer {
	return &SnippetContainer{}
}

func (sc *SnippetContainer) AddContext(key, value string) {
	for _, s := range sc.subContainers {
		s.AddContext(key, value)
	}
	for _, s := range sc.snippets {
		s.AddContext(key, value)
	}
}

func (sc *SnippetContainer) AddSnippet(snippet Snippet) {
	sc.snippets = append(sc.snippets, snippet)
}

func (sc *SnippetContainer) NewSubContainer() *SnippetContainer {
	sub := NewContainer()
	sc.subContainers = append(sc.subContainers, sub)
	return sub
}
