package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type any interface{}

func panicIf(x any) {
	if x != nil {
		log.Panic(x)
	}
}

func ensure(x any) {
	if x != nil {
		log.Panic(x)
	}
}

func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func mkdirIfNotExists(dirname string) {
	err := os.MkdirAll(dirname, 0666)
	if err != nil {
		log.Fatal("[ERROR]: mkdir %s failed: %#v", dirname, err)
	}
}

func dirname(file string) string {
	return filepath.Dir(file)
}

// basename(file) returns the basename of file with extname,
// without dirname (of cause.)
func basename(file string) string {
	return filepath.Base(file)
}

// basenameWithoutExt(file) returns the basename of file without extname,
// without dirname (of cause.)
func basenameWithoutExt(file string) string {
	f := filepath.Base(file)
	i := strings.LastIndex(f, ".")
	if i < 0 {
		return f
	}

	return f[0:i]
}

func isDir(file string) bool {
	f, err := os.Stat(file)
	if err != nil {
		return false
	}

	return f.IsDir()
}

func listDirFiles(dir string) []string {
	files, err := filepath.Glob(dir + "/*")
	if err != nil {
		return []string{}
	}

	return files
}

func formatError(r any) string {
	return fmt.Sprintf("%#v", r)
}
