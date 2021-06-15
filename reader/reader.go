package reader

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	Bit          = 1
	Byte         = 8 * Bit
	ZeroToken    = '0'
	OneToken     = '1'
	SpaceToken   = ' '
	CommentToken = "//"
)

type reader struct {
	src  io.ReadCloser
	scan *bufio.Scanner
	clen int
}

func NewReadCloser(src io.ReadCloser, commandLength int) io.ReadCloser {
	return &reader{
		src:  src,
		scan: bufio.NewScanner(src),
		clen: commandLength,
	}
}

func (r *reader) Read(dst []byte) (int, error) {
	if len(dst) < r.clen {
		return 0, fmt.Errorf("expected buffer length more than %d, got %d", r.clen, len(dst))
	}

	var row string
	for r.scan.Scan() {
		row = r.scan.Text()
		if len(row) == 0 {
			continue
		}
		if strings.HasPrefix(row, CommentToken) {
			continue
		}
		break
	}
	if err := r.scan.Err(); err != nil {
		return 0, err
	}
	if len(row) == 0 {
		return 0, io.EOF
	}

	offset := 0
	for _, sym := range []rune(row) {
		if offset >= r.clen * Byte {
			return 0, fmt.Errorf("line is too long: '%s'", row)
		}
		dst[offset / Byte] <<= 1
		switch sym {
		case SpaceToken:
			continue
		case ZeroToken:
		case OneToken:
			dst[offset / Byte] += 1
		default:
			return 0, fmt.Errorf("unexpected symbol: %#U", sym)
		}
		offset++
	}
	if offset < r.clen * Byte {
		return 0, fmt.Errorf("line is too short: '%s'", row)
	}
	return r.clen, nil
}

func (r *reader) Close() error {
	if err := r.src.Close(); err != nil {
		return err
	} else {
		return nil
	}
}
