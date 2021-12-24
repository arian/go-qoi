package go_qoi

import (
	"image"
	"image/png"
	"os"
)

func SaveImageToPngFile(img *image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	if err := png.Encode(f, *img); err != nil {
		f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
