package fop

import (
	"os"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func FileNotExists(filename string) bool {
	_, err := os.Stat(filename)
	return err != nil && os.IsNotExist(err)
}

func DirExists(dirname string) bool {
	stat, err := os.Stat(dirname)
	return err == nil && stat.IsDir()
}

// func DirNotExists(dirname string) bool {
// 	stat, err := os.Stat(dirname)
// 	return err == nil && stat.IsDir()
// }

func CreateFileIfNotExists(filename string) error {
	os.OpenFile(name, flag, perm)
}

func CreateDirIfNotExists(dirname string) error {

}
