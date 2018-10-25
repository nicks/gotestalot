package server

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type ResultReader struct {
	result *Result
	pos    int
}

func (r *ResultReader) Read(p []byte) (int, error) {
	c := r.result.cond
	c.L.Lock()
	defer c.L.Unlock()

	noData := r.pos >= r.result.buf.Len()
	endOfData := r.result.eof
	if noData && !endOfData {
		c.Wait()
	}

	noData = r.pos >= r.result.buf.Len()
	endOfData = r.result.eof
	if noData && endOfData {
		return 0, io.EOF
	}
	bs := r.result.buf.Bytes()
	n := copy(p, bs[r.pos:])
	r.pos += n
	return n, nil
}

type Result struct {
	buf  *bytes.Buffer
	mu   *sync.Mutex
	eof  bool
	cond *sync.Cond
}

func NewResult() *Result {
	mu := &sync.Mutex{}
	return &Result{
		buf:  bytes.NewBuffer(nil),
		mu:   mu,
		cond: sync.NewCond(mu),
	}
}

func (r *Result) NewReader() *ResultReader {
	return &ResultReader{
		result: r,
	}
}

func (r *Result) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	defer r.cond.Broadcast()
	if r.eof {
		return 0, fmt.Errorf("Error writing to closed result")
	}
	return r.buf.Write(p)
}

func (r *Result) End() {
	r.mu.Lock()
	defer r.mu.Unlock()
	defer r.cond.Broadcast()
	r.eof = true
}
