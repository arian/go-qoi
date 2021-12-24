package go_qoi

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestReadUint8(t *testing.T) {

	bs := []byte{23, 100, 255, 0}

	r := bytes.NewReader(bs)

	reader := bufio.NewReader(r)

	assertNext := func(expected uint8) {
		value, err := readUint8(reader)
		if err != nil {
			t.Error(err)
		}
		if value != expected {
			t.Errorf("got %d, expected %d", value, expected)
		}
	}

	assertNext(23)
	assertNext(100)
	assertNext(255)
	assertNext(0)
	_, err := readUint8(reader)
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestReadUint32(t *testing.T) {

	bs := []byte{23, 100, 255, 0}

	r := bytes.NewReader(bs)

	reader := bufio.NewReader(r)

	value, err := readUint32(reader)
	if err != nil {
		t.Error(err)
	}
	var expected uint32 = 392494848
	if value != expected {
		t.Errorf("got %d, expected %d", value, expected)
	}
}
