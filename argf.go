package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Argf struct {
	args   []string
	reader io.ReadCloser
}

func NewArgf(args []string) (*Argf, error) {
	if args == nil || len(args) < 1 {
		return &Argf{args: nil, reader: ioutil.NopCloser(os.Stdin)}, nil
	}
	reader, err := os.Open(args[0])
	if err != nil {
		return nil, err
	}
	return &Argf{args: args[1:], reader: reader}, nil
}

func (this *Argf) Read(data []byte) (int, error) {
	n, err := this.reader.Read(data)
	for {
		if err == io.EOF {
			if this.reader != nil {
				this.reader.Close()
			}
			this.reader = nil
			if this.args != nil && len(this.args) >= 1 {
				fname := this.args[0]
				this.args = this.args[1:]
				this.reader, err = os.Open(fname)
				if err != nil {
					return 0, fmt.Errorf("%s: %w", fname, err)
				}
			} else {
				return n, io.EOF
			}
		}
		if n >= len(data) {
			break
		}
		var m int
		m, err = this.Read(data[n:])
		n += m
	}
	return n, err
}

func (this *Argf) Close() error {
	var err error
	if this.reader != nil {
		err = this.reader.Close()
		this.reader = nil
	}
	return err
}
