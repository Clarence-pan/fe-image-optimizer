package main

import (
	"bufio"
	"io"
	"log"
)

// NTS buffered reader, data is from queue
type tBufReader struct {
	queue chan []byte
	buf   []byte
}

// TS system logger
type tSysLogger struct {
	queue chan []byte
}

var sysLogger *tSysLogger = newSysLogger()
var _ io.Writer = &tSysLogger{}
var _ io.WriteCloser = &tSysLogger{}
var _ io.Reader = &tBufReader{}
var _ io.ReadCloser = &tBufReader{}

func newSysLogger() *tSysLogger {
	t := &tSysLogger{
		queue: make(chan []byte, 4096),
	}

	go func() {
		r := bufio.NewReader(&tBufReader{queue: t.queue})

		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}

			log.Print(line)
		}
	}()

	return t
}

// Write() writes a block of data. it will be blocked if the queue is full.
func (t *tSysLogger) Write(p []byte) (n int, err error) {
	n = len(p)
	q := make([]byte, n)
	copy(q, p)

	t.queue <- q

	return
}

func (t *tSysLogger) Close() (err error) {
	close(t.queue)
	return nil
}

// Read() reads a block. it will be blocked if the queue is empty
func (t *tBufReader) Read(p []byte) (n int, err error) {
	if t.buf == nil {
		t.buf = <-t.queue
	}

	if t.buf == nil {
		return 0, io.EOF
	}

	bufLen := len(t.buf)
	pLen := len(p)
	if pLen < bufLen {
		n = pLen
		copy(p, t.buf[:n])
		t.buf = t.buf[n:]
	} else {
		n = bufLen
		copy(p, t.buf)
		t.buf = nil
	}

	return
}

func (t *tBufReader) Close() (err error) {
	close(t.queue)
	return nil
}
