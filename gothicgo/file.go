package gothicgo

import (
	"github.com/adamcolton/gothic"
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicio"
	"go/format"
	"io"
	"os"
	"path"
	"path/filepath"
)

// File represents a Go file. Writer is intended as a hook for testing. If it is
// nil, the code will be written to the file normally, if it's set to an
// io.WriteCloser, it will write to that instead
type File struct {
	*Imports
	generators *gothic.Project
	name       string
	code       []io.WriterTo
	pkg        *Package
	Writer     io.WriteCloser
	Comment    string
}

// Prepare runs prepare on all the generators in the file
func (f *File) Prepare() error {
	err := f.generators.Prepare()
	return errCtx(err, "Prepare file %s:", f.name)
}

// AddGenerators to the file
func (f *File) AddGenerators(generators ...gothic.Generator) {
	f.generators.AddGenerators(generators...)
}

// AddWriterTo adds a WriterTo that will be invoked when the file is written,
// but after Prepare and Generate have run.
func (f *File) AddWriterTo(writerTo io.WriterTo) {
	f.code = append(f.code, writerTo)
}

// Generate the file
func (f *File) Generate() error {
	f.Imports.ResolvePackages(f.Package().ImportResolver())
	f.Imports.RemoveRef(f.pkg)

	err := f.generators.Generate()
	if err != nil {
		return errCtx(err, "Generate file %s/%s:", f.pkg.name, f.name)
	}

	buf := bufpool.Get()
	sw := gothicio.NewSumWriter(buf)

	NewComment(f.Comment).WriteTo(sw)
	sw.WriteRune('\n')
	sw.WriteString("package ")
	sw.WriteString(f.pkg.name)
	sw.WriteString("\n\n")
	f.Imports.WriteTo(sw)
	gothicio.MultiWrite(sw, f.code, "\n")
	if sw.Err != nil {
		errCtx(sw.Err, "Generate file %s/%s:", f.pkg.name, f.name)
	}

	code := buf.Bytes()
	fmtCode, fmtErr := format.Source(code)

	wc := f.Writer
	if wc == nil {
		wc, err = f.open()
		if err != nil {
			return errCtx(err, "Generate file %s/%s:", f.pkg.name, f.name)
		}
	}

	if fmtErr == nil {
		_, err = wc.Write(fmtCode)
	} else {
		_, err = wc.Write(code)
		if err == nil {
			err = errCtx(fmtErr, "Failed to format %s/%s:", f.pkg.name, f.name)
		}
	}
	if err != nil {
		wc.Close()
		return errCtx(err, "Generate file %s/%s:", f.pkg.name, f.name)
	}
	return wc.Close()
}

// File creates a file within the package. The name should not include ".go"
// which will be automatically appended.
func (p *Package) File(name string) *File {
	if file, exists := p.files[name]; exists {
		return file
	}
	f := &File{
		Imports:    NewImports(p),
		generators: gothic.New(),
		name:       name,
		pkg:        p,
		Comment:    p.Comment,
	}
	p.files[name] = f
	return f
}

func (f *File) open() (io.WriteCloser, error) {
	pth := path.Join(f.pkg.OutputPath, f.name+".go")
	pth, err := filepath.Abs(pth)
	if err != nil {
		return nil, err
	}
	return os.Create(pth)
}

// Package the file is in
func (f *File) Package() *Package { return f.pkg }

// Name returns the name of the file.
func (f *File) Name() string { return f.name }

// String returns the file as Go code. This is intended for testing and
// debugging, not code generation
func (f *File) String() string {
	f.Prepare()
	buf := gothicio.BufferCloser{bufpool.Get()}
	f.Writer = buf
	f.Generate()
	return bufpool.PutStr(buf.Buffer)
}
