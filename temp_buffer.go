package goutils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
)

type TempBuffer struct {
	b      bytes.Buffer
	bw     *bufio.Writer
	offset int64
}

func NewTempBuffer() *TempBuffer {
	tb := &TempBuffer{}
	tb.bw = bufio.NewWriter(&tb.b)
	return tb
}

func (o *TempBuffer) WriteAt(p []byte, off int64) (n int, err error) {
	//fmt.Println("WriteAt", p, off)
	nn, err := o.b.Write(p)
	if err != nil {
		return nn, err
	}

	//fmt.Println("WriteAt...result", o.b, o.offset)

	return nn, nil
}
func (o *TempBuffer) Bytes() []byte {
	return o.b.Bytes()
}

func (o *TempBuffer) Read(p []byte) (int, error) {
	//fmt.Println("Read", p, o.offset)

	n, err := bytes.NewBuffer(o.b.Bytes()[o.offset:]).Read(p)

	if err == nil {
		if o.offset+int64(len(p)) < int64(o.b.Len()) {
			o.offset += int64(len(p))
		} else {
			o.offset = int64(o.b.Len())
		}
	}

	//fmt.Println("Read", p, o.offset)

	return n, err
}

func (o *TempBuffer) Write(p []byte) (int, error) {
	//fmt.Println("Write", p, o.offset)

	n, err := o.b.Write(p)

	if err == nil {
		o.offset = int64(o.b.Len())
	}

	//fmt.Println("Write", o.offset)

	return n, err
}

func (o *TempBuffer) Seek(offset int64, whence int) (int64, error) {

	fmt.Println("Seek", offset, whence)
	var err error

	switch whence {
	case 0: // 시작 위치 기준
		if offset >= int64(o.b.Len()) || offset < 0 {
			err = errors.New("Invalid Offset.")
		} else {
			o.offset = offset
		}
	case 1: // 현재 위치 기준
		if o.offset+offset >= int64(o.b.Len()) ||
			o.offset+offset < 0 {
			err = errors.New("Invalid Offset.")
		} else {
			o.offset += offset
		}

	case 2: // 끝 위치 기준
		if int64(o.b.Len())-offset < 0 ||
			offset >= int64(o.b.Len()) {
			err = errors.New("Invalid Offset.")
		} else {
			o.offset = int64(o.b.Len()) - offset
		}
	default:
		err = errors.New("Unsupported Seek Method.")
	}

	fmt.Println("Seek result", o.offset)

	return o.offset, err
}
