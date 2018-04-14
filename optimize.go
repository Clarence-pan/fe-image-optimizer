package main

import (
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type tOptimizer interface {
	optimize() error
	inputFile() string
	setInputFileContent(fileContent []byte)
	setInputFileReader(fileReader io.Reader)
}

func doOptimize(file string) {
	if strings.ContainsAny(file, "*?") {
		files, e := filepath.Glob(file)
		if e != nil {
			log.Printf("failed to match `%s`: %v", file, e)
			return
		}

		doOptimizeFiles(files)

		return
	}

	doOptimizeFile(file)
}

func doOptimizeFiles(files []string) {
	for _, x := range files {
		doOptimizeFile(x)
	}
}

func doOptimizeFile(file string) {
	if cfg.isInBlackList(file) {
		return
	}

	if cfg.isOptimizedFile(file) {
		log.Printf("ignore optimized %s", file)
		return
	}

	if isDir(file) {
		log.Printf("Enter directory %s.", file)

		for _, x := range listDirFiles(file) {
			doOptimizeFile(x)
		}

		log.Printf("Leave directory %s.", file)
		return
	}

	log.Printf("Processing %s...", file)
	optimizer, err := newOptimizer(file)
	if err != nil {
		log.Printf("ignore %s: %v", file, err)
		return
	}

	err = optimizer.optimize()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Processed %s.", file)
	return
}

func newOptimizer(file string) (opt tOptimizer, err error) {
	extName := filepath.Ext(file)

	switch strings.ToLower(extName) {
	case ".png":
		return newPngOptimizer(file), nil
	case ".jpg":
		return newJpgOptimizer(file), nil
	case ".jpeg":
		return newJpgOptimizer(file), nil
	case ".zip":
		return newZipOptimizer(file), nil
	default:
		return nil, errors.Errorf("unknown file type: %s", extName)
	}
}
