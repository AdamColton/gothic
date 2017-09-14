package gothicgo

import (
	"github.com/adamcolton/gothic"
	"regexp"
	"strings"
)

var importPath string //TODO: can we deduce the Import path with importResolver?
var OutputPath = "./Project"

var importPathRe = regexp.MustCompile(`^([\w\-\.]+\/)*$`)

const ErrBadImportPath = errStr("Bad Import Path")

func SetImportPath(path string) error {
	if !importPathRe.MatchString(path) {
		return ErrBadImportPath
	}
	importPath = path
	return nil
}

var DefaultComment = "This code was generated from a Gothic Blueprint, DO NOT MODIFY"
var CommentWidth = 80

var packages = gothic.New()

func Prepare()  { packages.Prepare() }
func Generate() { packages.Generate() }
func Export()   { packages.Export() }

func init() {
	gothic.AddGenerators(packages)
}

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
