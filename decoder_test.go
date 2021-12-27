package go_qoi

import (
	"image"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			t.Parallel()

			qoiImg, err := ReadQoiFile(path.Join(dir, name))
			assert.Nil(t, err)

			pngImg, err := ReadPngFile(path.Join(dir, name[:len(name)-3]+"png"))
			assert.Nil(t, err)

			assertEqualImages(t, pngImg, qoiImg)
		})
	}
}

func assertEqualImages(t *testing.T, expected, actual image.Image) {
	cond := func() bool {
		b := actual.Bounds()
		assert.Equal(t, expected.Bounds(), b)

		for x := b.Min.X; x < b.Max.X; x++ {
			for y := b.Min.Y; y < b.Max.Y; y++ {
				exPx := colorToNRGBA(expected.At(x, y))
				acPx := colorToNRGBA(actual.At(x, y))
				ok := assert.Equalf(t, exPx, acPx, "pixel at (%d,%d) should match", x, y)
				if !ok {
					return false
				}
			}
		}

		return true
	}
	assert.Condition(t, cond)
}
