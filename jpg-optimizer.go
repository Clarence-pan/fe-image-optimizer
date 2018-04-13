package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"
)

type tJpgOptimizer struct {
	tBaseOptimizer
}

func newJpgOptimizer(file string) *tJpgOptimizer {
	return &tJpgOptimizer{
		tBaseOptimizer: tBaseOptimizer{
			file: file,
		},
	}
}

func (opt *tJpgOptimizer) optimize() (err error) {
	defer func() {
		r := recover()
		if r != nil {
			var ok bool
			if err, ok = r.(error); ok {
				err = errors.Wrapf(err, "failed to optimize %s", opt.file)
			} else {
				err = errors.Errorf("failed to optimize %s, detail: %#v", opt.file, r)
			}
		}
	}()

	inputReader := opt.getInputFileReader()

	inputFileContent, err := ioutil.ReadAll(inputReader)
	panicIf(err)

	jpgOutputFileWriter := opt.getOutputFileWriter(".jpg")
	defer jpgOutputFileWriter.Close()
	optimizeJpgFile(bytes.NewReader(inputFileContent), jpgOutputFileWriter)

	pngBuf := jpgToPng(bytes.NewReader(inputFileContent))
	pngOutputFileWriter := opt.getOutputFileWriter(".png")
	defer pngOutputFileWriter.Close()
	optimizePngFile(pngBuf, pngOutputFileWriter)

	return nil
}

func optimizeJpgFile(inputReader io.Reader, outputWriter io.WriteCloser) {
	exe := cfg.Jpegoptim
	args := []string{
		"--force",
		"--verbose",
		"--max=" + strconv.Itoa(cfg.Jpeg.MaxQuality),
		"--stdout",
		"--stdin",
	}

	execCommand(exe, args, inputReader, outputWriter)
}
