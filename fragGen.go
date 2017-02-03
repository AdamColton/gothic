package gothic

// FragGen is a fragment generator is a useful interface for sub-generators. They
// will return a slice of strings when Generate is called.
type FragGen interface {
	Prepare()
	Generate() []string
}

// FG is a fragment generator container which itself implements the FragGen
// interface. Useful to embed.
type FG struct {
	generators []FragGen
}

// Prepare calls Prepare() on all Fragment Generators
func (fg *FG) Prepare() {
	for _, g := range fg.generators {
		g.Prepare()
	}
}

// Generate calls Generate() on all Fragment Generators
func (fg *FG) Generate() []string {
	s := []string{}
	for _, g := range fg.generators {
		s = append(s, g.Generate()...)
	}
	return s
}

// AddFragGen adds a Fragment Generator
func (fg *FG) AddFragGen(g FragGen) {
	fg.generators = append(fg.generators, g)
}

// SliceFG is a helper, it impements the FragGen interface on a string slice
type SliceFG []string

// Prepare just exists to fulfill the interface, doesn't do anything
func (s SliceFG) Prepare() {}

// Generate returns SliceFG as []string
func (s SliceFG) Generate() []string { return []string(s) }

// SliceFGFromString takes one or more strings and returns a fragment generator
func SliceFGFromString(strings ...string) SliceFG {
	return SliceFG(strings)
}
