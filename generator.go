package gothic

// A Generator is expected to write any relevant data to a persistant medium
// when Generate is called
type Generator interface {
	Prepare()
	Generate()
}

// A Project is a collection of generators. It also implements the Generator
// interface, so projects can contain sub-projects.
type Project struct {
	generators []Generator
}

// Returns a new Project
func New() *Project {
	return &Project{}
}

// Calls prepare on all Generators
func (p *Project) Prepare() {
	for _, g := range p.generators {
		g.Prepare()
	}
}

// Calls Generate on all Generators
func (p *Project) Generate() {
	for _, g := range p.generators {
		g.Generate()
	}
}

// Wraps a call to Prepare() and Generate()
func (p *Project) Export() {
	p.Prepare()
	p.Generate()
}

// Adds a generator
func (p *Project) AddGenerator(g Generator) {
	p.generators = append(p.generators, g)
}

var allGenerators = New()

func AddGenerator(g Generator) { allGenerators.AddGenerator(g) }
func Prepare()                 { allGenerators.Prepare() }
func Generate()                { allGenerators.Generate() }
func Export()                  { allGenerators.Export() }
