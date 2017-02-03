package gothic

// Fragment Generator is a useful interface for sub-generators. They will return
// a slice of strings when Generate is called.
type FragGen interface {
	Prepare()
	Generate() []string
}

// Fragment Generator container. Useful to embed.
type FG struct {
	generators []FragGen
}

// Calls prepare on all Fragment Generators
func (fg *FG) Prepare() {
	for _, g := range fg.generators {
		g.Prepare()
	}
}

// Calls prepare on all Fragment Generators
func (fg *FG) Generate() []string {
	s := []string{}
	for _, g := range fg.generators {
		s = append(s, g.Generate()...)
	}
	return s
}

// Adds a Fragment Generator
func (fg *FG) AddFragGen(g FragGen) {
	fg.generators = append(fg.generators, g)
}

// SliceFG is a helper, it impementsthe FragGen interface on a string slice
type SliceFG []string

// Just exists to fulfill the interface, doesn't do anything
func (s SliceFG) Prepare() {}

// Generate returns SliceFG as []string
func (s SliceFG) Generate() []string { return []string(s) }

// Often, only a single string is required
func SliceFGFromString(str string) SliceFG {
	return SliceFG([]string{str})
}
