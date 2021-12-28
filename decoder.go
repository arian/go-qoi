package go_qoi

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
)

type decoder struct {
	reader      bufio.Reader
	header      qoiHeader
	img         image.Image
	rgba        *image.NRGBA
	px          color.NRGBA
	run         uint8
	seen_pixels [64]color.NRGBA
}

type qoiHeader struct {
	magic      []byte
	width      uint32
	height     uint32
	channels   uint8
	colorspace uint8
}

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
		px:     color.NRGBA{R: 0, G: 0, B: 0, A: 255},
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
	if string(magic) != QOI_MAGIC {
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
	channels, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if channels != 3 && channels != 4 {
		return nil, fmt.Errorf("invalid channels value")
	}

	// colorspace
	colorspace, err := reader.ReadByte()
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
	w := int(d.header.width)
	h := int(d.header.height)
	l := w * h

	for i := 0; i < l; i++ {
		err := d.runAtPix(i, i%w, i/w)
		if err != nil {
			return err
		}
	}

	return nil
}

// copy the px to the image and store it in the seen_pixels cache
func (d *decoder) copyPx(x, y int) {
	p := d.px
	hash := indexPositionHash(p)
	d.seen_pixels[hash] = p
	d.rgba.Set(x, y, p)
}

// only set the px to the image, when we already know for sure the px is in the seen_pixels cache
func (d *decoder) setPx(x, y int) {
	d.rgba.Set(x, y, d.px)
}

func (d *decoder) runAtPix(i, x, y int) error {
	if d.run > 0 {
		d.run--
		d.setPx(x, y)
		return nil
	}

	b1, err := d.reader.ReadByte()
	if err != nil {
		return err
	}

	if b1 == QOI_OP_RGB {
		r, err := d.reader.ReadByte()
		if err != nil {
			return err
		}
		d.px.R = r
		g, err := d.reader.ReadByte()
		if err != nil {
			return err
		}
		d.px.G = g
		b, err := d.reader.ReadByte()
		if err != nil {
			return err
		}
		d.px.B = b

		d.copyPx(x, y)
		return nil
	}

	if b1 == QOI_OP_RGBA {
		r, err := d.reader.ReadByte()
		if err != nil {
			return err
		}
		d.px.R = r
		g, err := d.reader.ReadByte()
		if err != nil {
			return err
		}
		d.px.G = g
		b, err := d.reader.ReadByte()
		if err != nil {
			return err
		}
		d.px.B = b
		a, err := d.reader.ReadByte()
		if err != nil {
			return err
		}
		d.px.A = a

		d.copyPx(x, y)
		return nil
	}

	if (b1 & QOI_OP_MASK2) == QOI_OP_INDEX {
		d.px = d.seen_pixels[b1]
		d.setPx(x, y)
		return nil
	}

	if (b1 & QOI_OP_MASK2) == QOI_OP_DIFF {
		d.px.R = d.px.R + ((b1 >> 4) & 0b11) - 2
		d.px.G = d.px.G + ((b1 >> 2) & 0b11) - 2
		d.px.B = d.px.B + ((b1 >> 0) & 0b11) - 2
		d.copyPx(x, y)
		return nil
	}

	if (b1 & QOI_OP_MASK2) == QOI_OP_LUMA {
		b2, err := d.reader.ReadByte()
		if err != nil {
			return err
		}
		vg := (b1 & 0x3f) - 32
		d.px.R = d.px.R + vg - 8 + ((b2 >> 4) & 0x0f)
		d.px.G = d.px.G + vg
		d.px.B = d.px.B + vg - 8 + (b2 & 0x0f)

		d.copyPx(x, y)
		return nil
	}

	if (b1 & QOI_OP_MASK2) == QOI_OP_RUN {
		d.run = uint8(b1 & 0b00111111)
		d.setPx(x, y)
		return nil
	}

	return nil
}
