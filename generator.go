package gothic

import "fmt"

/*
Todo: Clean up state and error handling
*/

// Generator is expected to write any relevant data to a persistant medium when
// Generate is called
type Generator interface {
	Prepare() error
	Generate() error
}

const (
	StateError = int8(iota - 1)
	StateReady
	StatePrepared
	StateGenerated
)

// Project is a collection of generators. It also implements the Generator
// interface, so projects can contain sub-projects.
type Project struct {
	state      int8
	generators []Generator
}

// New Project
func New() *Project {
	return &Project{}
}

// Prepare calls Prepare() on all Generators
func (p *Project) Prepare() error {
	if p.state != StateReady {
		return fmt.Errorf("Bad State")
	}
	for _, g := range p.generators {
		err := g.Prepare()
		if err != nil {
			p.state = StateError
			return err
		}
	}
	p.state = StatePrepared
	return nil
}

// Generate calls Generate() on all Generators
func (p *Project) Generate() error {
	if p.state != StatePrepared {
		return fmt.Errorf("Bad State")
	}
	for _, g := range p.generators {
		err := g.Generate()
		if err != nil {
			p.state = StateError
			return err
		}
	}
	return nil
}

// Export wraps a call to Prepare() and Generate()
func (p *Project) Export() error {
	err := p.Prepare()
	if err != nil {
		return err
	}
	return p.Generate()
}

// AddGenerator adds a generator to the project
func (p *Project) AddGenerators(g ...Generator) error {
	if p.state != StateReady {
		return fmt.Errorf("Bad State")
	}
	p.generators = append(p.generators, g...)
	return nil
}

var allGenerators = New()

// AddGenerator adds a generator to the global Project
func AddGenerators(g ...Generator) error { return allGenerators.AddGenerators(g...) }

// Prepare calls Prepare() on the global Project
func Prepare() error { return allGenerators.Prepare() }

// Generate calls Generate() on the global project
func Generate() error { return allGenerators.Generate() }

// Export calls Export() on the global project
func Export() error { return allGenerators.Export() }
