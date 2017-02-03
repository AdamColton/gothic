package gothic

// Generator is expected to write any relevant data to a persistant medium when
// Generate is called
type Generator interface {
	Prepare()
	Generate()
}

// Project is a collection of generators. It also implements the Generator
// interface, so projects can contain sub-projects.
type Project struct {
	generators []Generator
}

// New Project
func New() *Project {
	return &Project{}
}

// Prepare calls Prepare() on all Generators
func (p *Project) Prepare() {
	for _, g := range p.generators {
		g.Prepare()
	}
}

// Generate calls Generate() on all Generators
func (p *Project) Generate() {
	for _, g := range p.generators {
		g.Generate()
	}
}

// Export wraps a call to Prepare() and Generate()
func (p *Project) Export() {
	p.Prepare()
	p.Generate()
}

// AddGenerator adds a generator to the project
func (p *Project) AddGenerator(g Generator) {
	p.generators = append(p.generators, g)
}

var allGenerators = New()

// AddGenerator adds a generator to the global Project
func AddGenerator(g Generator) { allGenerators.AddGenerator(g) }

// Prepare calls Prepare() on the global Project
func Prepare() { allGenerators.Prepare() }

// Generate calls Generate() on the global project
func Generate() { allGenerators.Generate() }

// Export calls Export() on the global project
func Export() { allGenerators.Export() }
