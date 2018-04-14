package main

import "io"

type tReaderAtDataImpl struct {
	data []byte
}

func newReaderAtFromBuf(buf []byte) io.ReaderAt {
	return &tReaderAtDataImpl{
		data: buf,
	}
}

func (buf *tReaderAtDataImpl) ReadAt(p []byte, offset int64) (n int, err error) {
	n = len(buf.data) - int(offset)
	if n <= 0 {
		n = 0
		return
	}

	if n > len(p) {
		n = len(p)
	}

	copy(p[:n], buf.data[offset:int(offset)+n])
	return
}
