package go_qoi

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoderHeader(t *testing.T) {

	img, err := ReadQoiFile("qoi_test_images/testcard.qoi")
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	err = Encode(&b, img)
	if err != nil {
		t.Fatal(err)
	}

	h := b.Next(4)
	assert.Equal(t, string(h), QOI_MAGIC)
}

func TestOnePxImg(t *testing.T) {

	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.NRGBA{R: 100})

	var b bytes.Buffer
	err := Encode(&b, img)
	if err != nil {
		t.Fatal(err)
	}

	bs := b.Bytes()

	// header
	assert.Equal(t, QOI_MAGIC, string(bs[:4]))
	assert.Equal(t, byte(1), bs[7])
	assert.Equal(t, byte(1), bs[11])
	assert.Equal(t, byte(4), bs[12])
	assert.Equal(t, byte(0), bs[13])
	// pixels
	assert.Equal(t, byte(QOI_OP_RGBA), bs[14])
	assert.Equal(t, byte(100), bs[15])
	assert.Equal(t, byte(0), bs[16])
	assert.Equal(t, byte(0), bs[17])
	assert.Equal(t, byte(0), bs[18])
	assert.Equal(t, byte(0), bs[19])
	// padding
	assert.Equal(t, QOI_END_PADDING[:], bs[19:])
}

func TestEncodeDiff(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.NRGBA{R: 100})
	img.Set(1, 0, color.NRGBA{R: 100, G: 1})

	var b bytes.Buffer
	err := Encode(&b, img)
	assert.Nil(t, err)
	bs := b.Bytes()

	// pixels
	assert.Equal(t,
		[]byte{QOI_OP_RGBA, 100, 0, 0, 0},
		bs[14:19],
	)
	assert.Equal(t,
		[]byte{QOI_OP_DIFF | 2<<4 + 3<<2 + 2},
		bs[19:20],
	)

	// padding
	assert.Equal(t, QOI_END_PADDING[:], bs[20:])
}

func TestEncodeLuma(t *testing.T) {

	img := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.NRGBA{R: 100, G: 2})
	img.Set(1, 0, color.NRGBA{R: 102, G: 0})

	var b bytes.Buffer
	var w io.Writer
	w = &b

	err := Encode(w, img)
	assert.Nil(t, err)

	bs := b.Bytes()

	// header
	// pixels
	assert.Equal(t,
		[]byte{QOI_OP_RGBA, 100, 2, 0, 0},
		bs[14:19],
	)

	assert.Equal(t,
		[]byte{
			QOI_OP_LUMA | 30,
			12<<4 | 10,
		},
		bs[19:21],
	)

	// padding
	assert.Equal(t, QOI_END_PADDING[:], bs[21:])
}

func TestEncoderBytes(t *testing.T) {
	bs, err := os.ReadFile("qoi_test_images/testcard.qoi")
	if err != nil {
		t.Fatal(err)
	}

	img, err := Decode(bytes.NewBuffer(bs))
	if err != nil {
		t.Fatal(err)
	}

	b := new(bytes.Buffer)
	err = Encode(b, img)

	assert.Nil(t, err)
	assert.Equal(t, bs, b.Bytes())
}

func TestEncoderAndDecode(t *testing.T) {
	dir := "qoi_test_images"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		name := f.Name()
		ext := name[len(name)-4:]
		if ext != ".png" {
			continue
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			img, err := ReadPngFile(path.Join(dir, name))
			if err != nil {
				t.Fatal(err)
			}

			var b bytes.Buffer

			err = Encode(&b, img)
			assert.Nil(t, err)

			imgQoi, err := Decode(&b)
			if !assert.Nil(t, err) {
				return
			}

			assertEqualImages(t, img, imgQoi)
		})
	}
}
