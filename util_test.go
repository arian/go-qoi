package go_qoi

import (
	"bufio"
	"bytes"
	"image/color"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadUint8(t *testing.T) {

	bs := []byte{23, 100, 255, 0}

	r := bytes.NewReader(bs)

	reader := bufio.NewReader(r)

	assertNext := func(expected uint8) {
		value, err := readUint8(reader)
		assert.Nil(t, err)
		assert.Equal(t, value, expected)
	}

	assertNext(23)
	assertNext(100)
	assertNext(255)
	assertNext(0)
	_, err := readUint8(reader)
	assert.ErrorIs(t, err, io.EOF)
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
	assert.Equal(t, value, expected)
}

func TestWriteUint32(t *testing.T) {
	r := new(bytes.Buffer)

	err := writeUint32(r, 392494848)
	assert.Nil(t, err)

	expected := []byte{23, 100, 255, 0}
	assert.Equal(t, expected, r.Bytes())
}

func TestWriteUint8(t *testing.T) {
	r := new(bytes.Buffer)

	assertNext := func(i uint8) {
		err := writeUint8(r, i)
		assert.Nil(t, err)
	}

	assertNext(23)
	assertNext(100)
	assertNext(255)
	assertNext(0)

	expected := []byte{23, 100, 255, 0}
	assert.Equal(t, expected, r.Bytes())
}

func TestColorEquals(t *testing.T) {
	assert.Equal(
		t,
		color.Black,
		color.Black,
	)
	assert.NotEqual(
		t,
		color.Black,
		color.White,
	)
}
