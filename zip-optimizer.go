package main

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/pkg/errors"
)

type tZipOptimizer struct {
	tBaseOptimizer
	extractDir string
}

func newZipOptimizer(file string) *tZipOptimizer {
	return &tZipOptimizer{
		tBaseOptimizer: tBaseOptimizer{
			file: file,
		},
	}
}

func (opt *tZipOptimizer) optimize() error {
	defer func() {
		r := recover()
		if r != nil {
			if err, ok := r.(error); ok {
				err = errors.Wrapf(err, "failed to optimize %s", opt.file)
			} else {
				err = errors.Errorf("failed to optimize %s, detail: %#v", opt.file, r)
			}
		}
	}()

	inputReader := opt.getInputFileReader()

	inputFileContent, err := ioutil.ReadAll(inputReader)
	panicIf(err)

	zipReader, err := zip.NewReader(newReaderAtFromBuf(inputFileContent), int64(len(inputFileContent)))
	panicIf(err)

	opt.extractDir = filepath.Join(dirname(opt.file), basenameWithoutExt(opt.file))
	mkdirIfNotExists(opt.extractDir)

	log.Printf("there are %d files in %s", len(zipReader.File), opt.file)
	for i, file := range zipReader.File {
		log.Printf("processing %d/%d: %s in %s", i + 1, len(zipReader.File), file.Name, opt.file)
		opt.dealFileInArchive(file)
	}

	return nil
}

func (opt *tZipOptimizer) dealFileInArchive(file *zip.File) {
	if cfg.isInBlackList(file.Name) {
		return
	}

	if cfg.isOptimizedFile(file.Name) {
		log.Printf("ignore optimized %s in zip file %s", file.Name, opt.file)
		return
	}

	r, err := file.Open()
	panicIf(err)
	defer r.Close()

	log.Printf("extracting %s from %s", file.Name, opt.file)
	fileContent, err := ioutil.ReadAll(r)
	panicIf(err)

	extractedFilePathName := filepath.Join(opt.extractDir, file.Name)
	extractedFileDirName := filepath.Dir(extractedFilePathName)
	if extractedFileDirName != opt.extractDir {
		mkdirIfNotExists(extractedFileDirName)
	}

	ensure(ioutil.WriteFile(extractedFilePathName, fileContent, 0666))

	fileOptimizer, err := newOptimizer(extractedFilePathName)
	if err != nil {
		log.Printf("ignore %s: %v", extractedFilePathName, err)
		return
	}

	log.Printf("Processing %s...", extractedFilePathName)
	fileOptimizer.setInputFileContent(fileContent)
	ensure(fileOptimizer.optimize())

	log.Printf("Processed %s.", extractedFilePathName)
}
