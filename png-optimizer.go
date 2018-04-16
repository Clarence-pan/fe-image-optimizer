package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

type tPngOptimizer struct {
	tBaseOptimizer
}

func newPngOptimizer(file string) *tPngOptimizer {
	return &tPngOptimizer{
		tBaseOptimizer: tBaseOptimizer{
			file: file,
		},
	}
}

func (opt *tPngOptimizer) optimize() (err error) {
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

	inputReader := opt.getInputFileReader()

	inputFileContent, err := ioutil.ReadAll(inputReader)
	panicIf(err)

	pngOutputFileWriter := opt.getOutputFileWriter(".png")
	defer pngOutputFileWriter.Close()
	optimizePngFile(bytes.NewReader(inputFileContent), pngOutputFileWriter)

	jpgBuf := pngToJpg(bytes.NewReader(inputFileContent))
	jpgOutputFileWriter := opt.getOutputFileWriter(".jpg")
	defer jpgOutputFileWriter.Close()
	optimizeJpgFile(jpgBuf, jpgOutputFileWriter)

	return nil
}

func optimizePngFile(inputReader io.Reader, outputWriter io.WriteCloser) {
	exe := cfg.Pngquant
	args := []string{
		"--force",
		"--verbose",
		"--speed", "1",
		"--quality", fmt.Sprintf("%d-%d", cfg.Png.MinQuality, cfg.Png.MaxQuality),
		"--strip",
		"-",
	}

	execCommand(exe, args, inputReader, outputWriter)
}
