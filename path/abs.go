package path

import (
	"os"
	"path/filepath"
)

var absroot, _ = filepath.Abs(filepath.Dir(os.Args[0]))

func RootPath() string { return absroot }

func AbsPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(RootPath(), path)
}
