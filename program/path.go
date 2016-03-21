package program

import (
	"os"
	"path/filepath"
)

var absroot, _ = filepath.Abs(filepath.Dir(os.Args[0]))

func Root() string { return absroot }

func AbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(Root(), path)
}
