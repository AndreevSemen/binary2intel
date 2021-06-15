package writer

import (
	"fmt"
	"io"
)

type CheckSummer interface {
	Sum(line uint32, data []byte) byte
}

type writer struct {
	dst    io.WriteCloser
	sum    CheckSummer
	fmt    string
	eof    string
	wrdLen int
	num    uint32
}

func NewWriteCloser(dst io.WriteCloser, sum CheckSummer, cmdFmt, cmdEof string, wordLength int) io.WriteCloser {
	return &writer{
		dst:    dst,
		sum:    sum,
		fmt:    fmt.Sprintf("%s\n", cmdFmt),
		eof:    fmt.Sprintf("%s\n", cmdEof),
		wrdLen: wordLength,
		num:    0,
	}
}

func (w *writer) Write(src []byte) (int, error) {
	if len(src) < w.wrdLen {
		return 0, fmt.Errorf("src is too short to write it as intel command")
	}
	sum := w.sum.Sum(w.num, src[:w.wrdLen])
	if n, err := fmt.Fprintf(w.dst, w.fmt, w.num, src[:w.wrdLen], sum); err != nil {
		return 0, err
	} else {
		w.num++
		return n, nil
	}
}

func (w *writer) Close() error {
	if _, err := fmt.Fprint(w.dst, w.eof); err != nil {
		return err
	}
	if err := w.dst.Close(); err != nil {
		return err
	}
	return nil
}
