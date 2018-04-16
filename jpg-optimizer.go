package main

import (
	"bytes"
	"context"
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

func (opt *tJpgOptimizer) optimize(ctx context.Context) (err error) {
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

	ensureContextNotDone(ctx)
	inputFileContent, err := ioutil.ReadAll(inputReader)
	panicIf(err)

	ensureContextNotDone(ctx)
	jpgOutputFileWriter := opt.getOutputFileWriter(".jpg")
	defer jpgOutputFileWriter.Close()
	optimizeJpgFile(bytes.NewReader(inputFileContent), jpgOutputFileWriter, ctx)

	ensureContextNotDone(ctx)
	pngBuf := jpgToPng(bytes.NewReader(inputFileContent))
	pngOutputFileWriter := opt.getOutputFileWriter(".png")
	defer pngOutputFileWriter.Close()

	optimizePngFile(pngBuf, pngOutputFileWriter, ctx)

	return nil
}

func optimizeJpgFile(inputReader io.Reader, outputWriter io.WriteCloser, ctx context.Context) {
	exe := cfg.Jpegoptim
	args := []string{
		"--force",
		"--verbose",
		"--max=" + strconv.Itoa(cfg.Jpeg.MaxQuality),
		"--stdout",
		"--stdin",
	}

	execCommand(exe, args, inputReader, outputWriter, ctx)
}
