package gothicgo

import (
	"github.com/adamcolton/gothic"
	"regexp"
)

var importPath string //TODO: can we deduce the Import path with importResolver?

// OutputPath is where the project will be written
var OutputPath = "./Project"

var importPathRe = regexp.MustCompile(`^([\w\-\.]+\/)*$`)

// ErrBadImportPath indicates a poorly formatted import path. Path must end with
// /
const ErrBadImportPath = errStr("Bad Import Path")

// SetImportPath for the project. It is safe to change import path during
// generation, anything that uses teh default import path will get a copy at the
// time of it's instanciation.
func SetImportPath(path string) error {
	if !importPathRe.MatchString(path) {
		return ErrBadImportPath
	}
	importPath = path
	return nil
}

// DefaultComment to be included at the top of each generated file.
var DefaultComment = "This code was generated from a Gothic Blueprint, DO NOT MODIFY"

// CommentWidth to wrap comments
var CommentWidth = 80

var packages = gothic.New()

// Prepare all packages
func Prepare() error { return packages.Prepare() }

// Generate all packages
func Generate() error { return packages.Generate() }

// Export exports all the packages
func Export() error {
	return packages.Export()
}

func init() {
	gothic.AddGenerators(packages)
}
