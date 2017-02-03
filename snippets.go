package gothic

// SnippetContainer is anything that can hold a collection of snippets.
type SnippetContainer interface {
	AddSnippetTo(bucketName string, snippet Snippet)
	Snippets(bucketName string) []Snippet
}

// Snippet is only used to generate the SnippetInstance. In practise, the
// Snippet will normally be configured and each instance will inheirit that
// configuration.
type Snippet interface {
	New() SnippetInstance
}

// SnippetInstance is used to generate a string from a Snippet. By setting the
// token values, the SnippetInstance can generate a snippet that works in any
// arbitrary scope.
type SnippetInstance interface {
	Set(key, val string)
	FragGen
}

// SC implements SnippetContainer.
type SC struct {
	buckets map[string][]Snippet
}

func NewSC() *SC {
	return &SC{
		buckets: make(map[string][]Snippet),
	}
}

func (s *SC) AddSnippetTo(bucketName string, snippet Snippet) {
	s.buckets[bucketName] = append(s.buckets[bucketName], snippet)
}

func (s *SC) Snippets(bucketName string) []Snippet {
	return s.buckets[bucketName]
}
