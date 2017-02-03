package gothicgo

import (
	"fmt"
	"github.com/adamcolton/gothic"
	"go/format"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Represents a Go file. Writer is intended as a hook for testing. If it is nil,
// the code will be written to the file normally, if it's set to an
// io.WriteCloser, it will write to that instead
type File struct {
	*Imports
	fg      *gothic.FG
	name    string
	code    []string
	pkg     *Package
	Writer  io.WriteCloser
	Comment string
}

func (f *File) Prepare()                     { f.fg.Prepare() }
func (f *File) AddFragGen(fg gothic.FragGen) { f.fg.AddFragGen(fg) }
func (f *File) AddCode(code ...string)       { f.fg.AddFragGen(gothic.SliceFG(code)) }

func (f *File) Generate() {
	f.Imports.ResolvePackages(f.Package().ImportResolver())
	f.Imports.RemovePath(f.pkg.ImportPath)
	s := []string{
		BuildComment(f.Comment, CommentWidth),
		"package " + f.pkg.Name + "\n",
		f.Imports.String(),
	}

	s = append(s, f.fg.Generate()...)

	code := []byte(strings.Join(s, "\n"))
	fmtCode, err := format.Source(code)

	wc := f.Writer
	if wc == nil {
		wc = f.open()
	}
	if err == nil {
		wc.Write(fmtCode)
	} else {
		wc.Write(code)
		fmt.Println("Failed to format", f.pkg.Name+"/"+f.name+".go ", err)
	}
	wc.Close()
}

// Takes the file name without ".go"
func (p *Package) File(name string) *File {
	if file, exists := p.files[name]; exists {
		return file
	}
	f := &File{
		Imports: NewImports(),
		fg:      &gothic.FG{},
		name:    name,
		pkg:     p,
		Comment: p.Comment,
	}
	p.files[name] = f
	return f
}

func (f *File) open() io.WriteCloser {
	pth := path.Join(f.pkg.OutputPath, f.name+".go")
	pth, e := filepath.Abs(pth)
	if e != nil {
		panic(e)
	}
	wc, e := os.Create(pth)
	if e != nil {
		panic(e)
	}
	return wc
}

// Returns the package
func (f *File) Package() *Package { return f.pkg }
