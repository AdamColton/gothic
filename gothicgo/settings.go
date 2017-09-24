package gothicgo

import (
	"github.com/adamcolton/gothic"
	"regexp"
	"strings"
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
func Prepare() { packages.Prepare() }

// Generate all packages
func Generate() { packages.Generate() }

func init() {
	gothic.AddGenerators(packages)
}

// BuildComment creates a comment wrapped to a specific width.
func BuildComment(comment string, width int) string {
	targetWidth := width - 3
	lines := []string{}
	buf, line := " ", ""
	for _, c := range comment {
		switch c {
		case ' ':
			fallthrough
		case '\t':
			fallthrough
		case '\n':
			lb, ll := len(buf), len(line)
			if lb > targetWidth {
				if ll > 0 {
					lines = append(lines, line)
				}
				lines = append(lines, buf)
				buf, line = " ", ""
			} else if lb+ll > targetWidth {
				lines = append(lines, line)
				buf, line = " ", buf
			} else {
				line += buf
				buf = " "
			}
		default:
			buf += string(c)
		}
	}
	lb, ll := len(buf), len(line)
	if lb > targetWidth || lb+ll > targetWidth {
		if len(line) > 0 {
			lines = append(lines, line)
		}
		lines = append(lines, buf)
	} else {
		lines = append(lines, line+buf)
	}

	return "//" + strings.Join(lines, "\n//") + "\n\n"
}
