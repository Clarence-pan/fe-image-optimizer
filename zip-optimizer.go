package main

import (
	"archive/zip"
	"context"
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

func (opt *tZipOptimizer) optimize(ctx context.Context) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			if e, ok := r.(error); ok {
				err = errors.Wrapf(e, "failed to optimize %s", opt.file)
			} else {
				err = errors.Errorf("failed to optimize %s, detail: %#v", opt.file, r)
			}
		}
	}()

	ensureContextNotDone(ctx)
	inputReader := opt.getInputFileReader()

	ensureContextNotDone(ctx)
	inputFileContent, err := ioutil.ReadAll(inputReader)
	panicIf(err)

	ensureContextNotDone(ctx)
	zipReader, err := zip.NewReader(newReaderAtFromBuf(inputFileContent), int64(len(inputFileContent)))
	panicIf(err)

	ensureContextNotDone(ctx)
	opt.extractDir = filepath.Join(dirname(opt.file), basenameWithoutExt(opt.file))
	mkdirIfNotExists(opt.extractDir)

	ensureContextNotDone(ctx)
	log.Printf("there are %d files in %s", len(zipReader.File), opt.file)
	for i, file := range zipReader.File {
		log.Printf("processing %d/%d: %s in %s", i+1, len(zipReader.File), file.Name, opt.file)
		ensureContextNotDone(ctx)
		opt.dealFileInArchive(file, ctx)
	}

	return nil
}

func (opt *tZipOptimizer) dealFileInArchive(file *zip.File, ctx context.Context) {
	if cfg.isInBlackList(file.Name) {
		log.Printf("ignore black listed %s in zip file %s", file.Name, opt.file)
		return
	}

	if cfg.isOptimizedFile(file.Name) {
		log.Printf("ignore optimized %s in zip file %s", file.Name, opt.file)
		return
	}

	fileInfo := file.FileInfo()
	if fileInfo.IsDir() {
		mkdirIfNotExists(filepath.Join(opt.extractDir, file.Name))
		return
	}

	ensureContextNotDone(ctx)
	r, err := file.Open()
	panicIf(err)
	defer r.Close()

	ensureContextNotDone(ctx)
	log.Printf("extracting %s from %s", file.Name, opt.file)
	fileContent, err := ioutil.ReadAll(r)
	panicIf(err)

	extractedFilePathName := filepath.Join(opt.extractDir, file.Name)
	extractedFileDirName := filepath.Dir(extractedFilePathName)
	if extractedFileDirName != opt.extractDir {
		mkdirIfNotExists(extractedFileDirName)
	}

	ensureContextNotDone(ctx)
	ensure(ioutil.WriteFile(extractedFilePathName, fileContent, 0666))

	ensureContextNotDone(ctx)
	fileOptimizer, err := newOptimizer(extractedFilePathName)
	if err != nil {
		log.Printf("ignore %s: %v", extractedFilePathName, err)
		return
	}

	ensureContextNotDone(ctx)
	log.Printf("Processing %s...", extractedFilePathName)
	fileOptimizer.setInputFileContent(fileContent)
	ensure(fileOptimizer.optimize(ctx))

	log.Printf("Processed %s.", extractedFilePathName)
}
