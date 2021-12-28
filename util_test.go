package go_qoi

import (
	"bufio"
	"bytes"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
