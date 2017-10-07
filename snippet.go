package gothic

// Snippet provides an interface for generating strings from a context
type Snippet interface {
	AddContext(key, value string) Snippet
	String() string
}

// SnippetContainer holds snippets and shares context between them
type SnippetContainer struct {
	subContainers []*SnippetContainer
	snippets      []Snippet
}

// NewContainer for snippets
func NewContainer() *SnippetContainer {
	return &SnippetContainer{}
}

// AddContext to a snippet
func (sc *SnippetContainer) AddContext(key, value string) {
	for _, s := range sc.subContainers {
		s.AddContext(key, value)
	}
	for _, s := range sc.snippets {
		s.AddContext(key, value)
	}
}

// AddSnippet to a container
func (sc *SnippetContainer) AddSnippet(snippet Snippet) {
	sc.snippets = append(sc.snippets, snippet)
}

// NewSubContainer added to SnippetContainer
func (sc *SnippetContainer) NewSubContainer() *SnippetContainer {
	sub := NewContainer()
	sc.subContainers = append(sc.subContainers, sub)
	return sub
}
