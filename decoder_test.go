package go_qoi

import (
	"bufio"
	"bytes"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path"
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

func TestDecoder(t *testing.T) {
	dir := "qoi_test_images"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		name := f.Name()
		ext := name[len(name)-4:]
		if ext != ".qoi" {
			continue
		}

		t.Run(name, func(t *testing.T) {
			qoiPath := path.Join(dir, name)
			pngPath := path.Join(dir, name[:len(name)-3]+"png")

			qoiF, err := os.Open(qoiPath)
			if err != nil {
				t.Fatal(err)
			}
			defer qoiF.Close()
			qoiImg, err := Decode(qoiF)
			if err != nil {
				t.Fatal(err)
			}

			pngF, err := os.Open(pngPath)
			if err != nil {
				t.Fatal(err)
			}
			defer pngF.Close()
			pngImg, err := png.Decode(pngF)
			if err != nil {
				t.Fatal(err)
			}

			if pngImg.Bounds() != qoiImg.Bounds() {
				t.Errorf("Bounds are not equal, expected %v, got %v\n", pngImg.Bounds(), qoiImg.Bounds())
			}

			b := qoiImg.Bounds()
			for x := b.Min.X; x < b.Max.X; x++ {
				for y := b.Min.Y; y < b.Max.Y; y++ {
					pngPx := pngImg.At(x, y)
					qoiPx := qoiImg.At(x, y)
					if !colorEquals(pngPx, qoiPx) {
						t.Errorf("pixels at (%d,%d) are not equal; expected %v, got %v\n", x, y, pngPx, qoiPx)
					}
				}
			}
		})
	}
}

func colorEquals(c, o color.Color) bool {
	cr, cg, cb, ca := c.RGBA()
	or, og, ob, oa := o.RGBA()
	return cr == or && cg == og && cb == ob && ca == oa
}
