package go_qoi

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
)

type decoder struct {
	reader bufio.Reader
	header qoiHeader
	img    image.Image
	rgba   *image.NRGBA
}

type qoiHeader struct {
	magic      []byte
	width      uint32
	height     uint32
	channels   uint8
	colorspace uint8
}

const (
	QOI_OP_INDEX = 0x00
	QOI_OP_DIFF  = 0x40
	QOI_OP_LUMA  = 0x80
	QOI_OP_RUN   = 0xC0
	QOI_OP_RGB   = 0xFE
	QOI_OP_RGBA  = 0xFF
	QOI_OP_MASK2 = 0xC0
)

func Decode(r io.Reader) (image.Image, error) {
	reader := bufio.NewReader(r)

	header, err := parseHeader(reader)
	if err != nil {
		return nil, err
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(header.width), int(header.height)))

	d := &decoder{
		reader: *reader,
		header: *header,
		img:    img,
		rgba:   img,
	}

	err = d.decode()
	if err != nil {
		return nil, err
	}

	return img, nil
}

func parseHeader(reader *bufio.Reader) (*qoiHeader, error) {
	// magic
	magic := make([]byte, 4)
	_, err := reader.Read(magic)
	if err != nil {
		return nil, err
	}
	if string(magic) != "qoif" {
		return nil, fmt.Errorf("File is not a qoi image")
	}

	// width
	width, err := readUint32(reader)
	if err != nil {
		return nil, err
	}

	// height
	height, err := readUint32(reader)
	if err != nil {
		return nil, err
	}

	// channels
	channels, err := readUint8(reader)
	if err != nil {
		return nil, err
	}
	if channels != 3 && channels != 4 {
		return nil, fmt.Errorf("invalid channels value")
	}

	// colorspace
	colorspace, err := readUint8(reader)
	if err != nil {
		return nil, err
	}
	if colorspace != 0 && colorspace != 1 {
		return nil, fmt.Errorf("invalid colorspace value")
	}

	return &qoiHeader{
		magic:      magic,
		width:      width,
		height:     height,
		channels:   channels,
		colorspace: colorspace,
	}, nil
}

func (d *decoder) decode() error {
	header := d.header
	reader := &d.reader

	w := int(header.width)
	h := int(header.height)
	l := w * h

	seen_pixels := make([]color.NRGBA, 64)
	px := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	run := 0

	for i := 0; i < l; i++ {
		if run > 0 {
			run--
		} else {

			b1, err := reader.ReadByte()
			if err != nil {
				return err
			}

			if b1 == QOI_OP_RGB {
				// QOI_OP_RGB

				r, err := readUint8(reader)
				if err != nil {
					return err
				}
				px.R = r
				g, err := readUint8(reader)
				if err != nil {
					return err
				}
				px.G = g
				b, err := readUint8(reader)
				if err != nil {
					return err
				}
				px.B = b

			} else if b1 == QOI_OP_RGBA {
				// QOI_OP_RGBA
				r, err := readUint8(reader)
				if err != nil {
					return err
				}
				px.R = r
				g, err := readUint8(reader)
				if err != nil {
					return err
				}
				px.G = g
				b, err := readUint8(reader)
				if err != nil {
					return err
				}
				px.B = b
				a, err := readUint8(reader)
				if err != nil {
					return err
				}
				px.A = a

			} else if (b1 & QOI_OP_MASK2) == QOI_OP_INDEX {
				// QOI_OP_INDEX
				px = seen_pixels[b1]
			} else if (b1 & QOI_OP_MASK2) == QOI_OP_DIFF {
				px.R = px.R + ((b1 >> 4) & 0b11) - 2
				px.G = px.G + ((b1 >> 2) & 0b11) - 2
				px.B = px.B + ((b1 >> 0) & 0b11) - 2

			} else if (b1 & QOI_OP_MASK2) == QOI_OP_LUMA {
				b2, err := readUint8(reader)
				if err != nil {
					return err
				}
				vg := (b1 & 0x3f) - 32
				px.R = px.R + vg - 8 + ((b2 >> 4) & 0x0f)
				px.G = px.G + vg
				px.B = px.B + vg - 8 + (b2 & 0x0f)

			} else if (b1 & QOI_OP_MASK2) == QOI_OP_RUN {
				run = int(b1 & 0b00111111)
			}
		}

		p := px
		hash := indexPositionHash(p)
		seen_pixels[hash] = p

		y := i / w
		x := i % w
		d.rgba.Set(x, y, p)
	}

	return nil
}

func indexPositionHash(px color.NRGBA) uint8 {
	return (px.R*3 + px.G*5 + px.B*7 + px.A*11) % 64
}

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