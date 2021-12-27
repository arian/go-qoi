package go_qoi

import (
	"bufio"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
)

const (
	QOI_OP_INDEX = 0x00
	QOI_OP_DIFF  = 0x40
	QOI_OP_LUMA  = 0x80
	QOI_OP_RUN   = 0xC0
	QOI_OP_RGB   = 0xFE
	QOI_OP_RGBA  = 0xFF
	QOI_OP_MASK2 = 0xC0
	QOI_MAGIC    = "qoif"
)

var QOI_PADDING = [8]byte{0, 0, 0, 0, 0, 0, 0, 1}

func readUint32(reader *bufio.Reader) (uint32, error) {
	bytes := make([]byte, 4)
	_, err := reader.Read(bytes)
	if err != nil {
		return 0, err
	}
	value := binary.BigEndian.Uint32(bytes)
	return value, nil
}

func readUint8(reader *bufio.Reader) (uint8, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint8(b), nil
}

func writeUint32(writer io.Writer, value uint32) error {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, value)
	_, err := writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func writeUint8(writer io.Writer, value uint8) error {
	_, err := writer.Write([]byte{value})
	if err != nil {
		return err
	}
	return nil
}

func colorToNRGBA(c color.Color) color.NRGBA {
	return color.NRGBAModel.Convert(c).(color.NRGBA)
}

func colorEquals(c, o color.Color) bool {
	c1 := colorToNRGBA(c)
	o1 := colorToNRGBA(o)
	return c1 == o1
}

func indexPositionHash(px color.NRGBA) uint8 {
	return (px.R*3 + px.G*5 + px.B*7 + px.A*11) % 64
}

func ReadPngFile(filename string) (image.Image, error) {
	i, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer i.Close()

	img, err := png.Decode(i)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ReadQoiFile(filename string) (image.Image, error) {
	i, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer i.Close()

	img, err := Decode(i)
	if err != nil {
		return nil, err
	}

	return img, nil
}

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

func ReadPngAndSaveImageToQoi(input, output string) error {
	o, err := os.Create(output)
	if err != nil {
		return err
	}
	defer o.Close()

	img, err := ReadPngFile(input)

	if err := Encode(o, img); err != nil {
		return err
	}

	return nil
}
