package main

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type tBaseOptimizer struct {
	file         string
	inputReader  io.Reader
	outputWriter io.WriteCloser
}

func (opt *tBaseOptimizer) setInputFileReader(reader io.Reader) {
	opt.inputReader = reader
}

func (opt *tBaseOptimizer) getInputFileReader() (reader io.Reader) {
	if opt.inputReader != nil {
		return opt.inputReader
	}

	reader, err := os.Open(opt.file)
	panicIf(err)

	return
}

func (opt *tBaseOptimizer) getOutputFileName(fileExtName string) string {
	i := strings.LastIndex(opt.file, ".")
	if i < 0 {
		return opt.file + ".optimized" + fileExtName
	}

	return opt.file[0:i] + ".optimized" + fileExtName
}

func (opt *tBaseOptimizer) setOutputFileWriter(writer io.WriteCloser) {
	opt.outputWriter = writer
}

func (opt *tBaseOptimizer) getOutputFileWriter(fileExtName string) (writer io.WriteCloser) {
	if opt.outputWriter != nil {
		return opt.outputWriter
	}

	outputFile := opt.getOutputFileName(fileExtName)

	file, err := os.OpenFile(outputFile, os.O_CREATE, 0666)
	panicIf(err)

	return file
}

func jpgToPng(inputBuf io.Reader) (outputBuf *bytes.Buffer) {
	img, err := jpeg.Decode(inputBuf)
	panicIf(err)

	outputBuf = bytes.NewBuffer(make([]byte, 0))
	png.Encode(outputBuf, img)

	return
}

func pngToJpg(inputBuf io.Reader) (outputBuf *bytes.Buffer) {
	img, err := png.Decode(inputBuf)
	panicIf(err)

	outputBuf = bytes.NewBuffer(make([]byte, 0))
	jpeg.Encode(outputBuf, img, &jpeg.Options{
		Quality: cfg.Jpeg.MaxQuality,
	})

	return
}

func execCommand(exe string, args []string, stdin io.Reader, stdout io.WriteCloser) {
	log.Printf("[EXEC]: %s %v", exe, args)

	cmd := exec.Command(exe, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	panicIf(err)
}
