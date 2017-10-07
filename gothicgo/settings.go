package gothicgo

import (
	"github.com/adamcolton/gothic"
	"io"
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

type Comment struct {
	Comment string
	Width   int
}

func NewComment(comment string) Comment {
	return Comment{
		Comment: comment,
		Width:   CommentWidth,
	}
}

const wsRunes = "\t "

type SumWriter struct {
	W   io.Writer
	Sum int64
	Err error
}

func (s *SumWriter) WriteString(str string) { s.Write([]byte(str)) }
func (s *SumWriter) WriteRune(r rune)       { s.Write([]byte(string(r))) }

func (s *SumWriter) Write(b []byte) {
	if s.Err != nil {
		return
	}
	var n int
	n, s.Err = s.W.Write(b)
	s.Sum += int64(n)
}

var commentStart = []byte("// ")
var nl = []byte("\n")

func (c Comment) WriteTo(w io.Writer) (int64, error) {
	sum := SumWriter{W: w}
	targetWidth := c.Width - 3
	until := len(c.Comment) - targetWidth

	cur := 0
	for cur < until && sum.Err == nil {
		s := c.Comment[cur : cur+targetWidth]
		end := strings.IndexRune(s, '\n')
		if end == -1 {
			end = strings.LastIndexAny(s, " \t")
			if end == -1 {
				s = c.Comment[cur:]
				end = strings.IndexAny(s, " \t\n")
				if end == -1 {
					break
				}
			}
		}
		sum.Write(commentStart)
		sum.Write([]byte(s[:end]))
		sum.Write(nl)
		cur += end + 1
	}

	sum.Write(commentStart)
	sum.Write([]byte(c.Comment[cur:]))
	sum.Write(nl)

	return sum.Sum, sum.Err
}
