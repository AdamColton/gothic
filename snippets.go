package gothic

// SnippetContainer is anything that can hold a collection of snippets.
type SnippetContainer interface {
	AddSnippetTo(bucketName string, snippet Snippet)
	Snippets(bucketName string) []Snippet
}

// Snippet is the concept of a small string generator with a few configurable
// values. Kind of like a named collection of Printf call.
//
// The Snippet interface is only used to generate the SnippetInstance. In
// practice, the Snippet will normally be configured and each instance will
// inherit that configuration.
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

// NewSC returns a new SC which implements the SnippetContainer interface
func NewSC() *SC {
	return &SC{
		buckets: make(map[string][]Snippet),
	}
}

// AddSnippetTo will add a snippet to a bucket in the container
func (s *SC) AddSnippetTo(bucketName string, snippet Snippet) {
	s.buckets[bucketName] = append(s.buckets[bucketName], snippet)
}

// Snippets returns the snippets associated with a bucket
func (s *SC) Snippets(bucketName string) []Snippet {
	return s.buckets[bucketName]
}
