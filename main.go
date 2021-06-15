package main

import (
	"bin2intel/check_sum"
	"bin2intel/reader"
	"bin2intel/writer"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	cmdLen = flag.Int("cmd-length", 4, "source command length")
	cmdFmt = flag.String("fmt", ":01%04X00%02X%02X", "encoded command format")
	cmdEof = flag.String("eof", ":01000FF", "trailing command")

	src = flag.String("f", "code.bin", "file with code to translate")
	dst = flag.String("o", "translated", "file for translated code prefix")
)

func init() {
	flag.Parse()
}

func main() {
	srcFile, err := os.Open(*src)
	if err != nil {
		log.Fatalf("open src file error: %s", err)
	}
	r := reader.NewReadCloser(srcFile, *cmdLen)
	defer r.Close()

	var ws = make([]io.WriteCloser, 0, *cmdLen)
	for i := 0; i < *cmdLen; i++ {
		dstFile, err := os.Create(fmt.Sprintf("%s_%d.hex", *dst, i))
		if err != nil {
			log.Fatalf("open %s_%d.hex file error: %s", *dst, i, err)
		}

		w := writer.NewWriteCloser(dstFile, check_sum.CheckSum{}, *cmdFmt, *cmdEof, 1)
		defer w.Close()

		ws = append(ws, w)
	}

	for {
		var cmd = make([]byte, *cmdLen)
		if _, err = r.Read(cmd); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("read error: %s", err)
		}

		for i := 0; i < *cmdLen; i++ {
			if _, err = ws[i].Write(cmd[i:i+1]); err != nil {
				log.Fatalf("write %s_%d.hex error: %s", *dst, i, err)
			}
		}
	}
}
