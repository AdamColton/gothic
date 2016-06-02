package blueprint

import (
	"path/filepath"
)

var AppPath string
var ImportPath string

func SetAppPath(path string) error {
	path, err := resolvePathWithTrailingSlash(path)
	AppPath = path
	return err
}

func SetImportPath(path string) error {
	path, err := resolvePathWithTrailingSlash(path)
	ImportPath = path
	return err
}

func resolvePathWithTrailingSlash(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err == nil {
		if path[len(path)-1:] != "/" {
			path += "/"
		}
	}
	return path, err
}

var GothicIncludePath = "github.com/adamcolton/gothic/"
