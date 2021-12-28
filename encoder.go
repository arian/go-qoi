package go_qoi

import (
	"image"
	"image/color"
	"io"
)

type encoder struct {
	w           io.Writer
	img         image.Image
	seen_pixels [64]color.NRGBA
	run         uint8
	prev_px     color.NRGBA
	length, i   int
}

func Encode(w io.Writer, img image.Image) error {

	_, err := w.Write([]byte(QOI_MAGIC))
	if err != nil {
		return err
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	err = writeUint32(w, uint32(width))
	if err != nil {
		return err
	}
	err = writeUint32(w, uint32(height))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte{4})
	if err != nil {
		return err
	}

	_, err = w.Write([]byte{0})
	if err != nil {
		return err
	}

	l := width * height
	encoder := &encoder{
		w:       w,
		img:     img,
		prev_px: color.NRGBA{A: 255},
		run:     0,
		length:  l,
	}

	for i := 0; i < l; i++ {
		err := encoder.runAtPix(i, i%width, i/width)
		if err != nil {
			return err
		}
	}

	_, err = w.Write(QOI_END_PADDING[:])
	if err != nil {
		return err
	}

	return nil
}

func (e *encoder) atLastPx() bool {
	return e.i == e.length-1
}

func (e *encoder) runAtPix(i, x, y int) error {
	e.i = i
	var err error

	px := colorToNRGBA(e.img.At(x, y))
	same_pixel := px == e.prev_px

	if same_pixel {
		e.run++
	}

	if (!same_pixel && e.run > 0) || (same_pixel && (e.run == 62 || e.atLastPx())) {
		_, err = e.w.Write([]byte{uint8(QOI_OP_RUN | (e.run - 1))})
		if err != nil {
			return err
		}
		e.run = 0
	}

	if same_pixel {
		return nil
	}

	index_pos := indexPositionHash(px)

	if e.seen_pixels[index_pos] == px {
		_, err = e.w.Write([]byte{QOI_OP_INDEX | index_pos})
		if err != nil {
			return err
		}
		e.prev_px = px
		return nil
	}

	e.seen_pixels[index_pos] = px

	if px.A != e.prev_px.A {
		_, err := e.w.Write([]byte{QOI_OP_RGBA, px.R, px.G, px.B, px.A})
		if err != nil {
			return err
		}
		e.prev_px = px
		return nil
	}

	vr := int8(px.R) - int8(e.prev_px.R)
	vg := int8(px.G) - int8(e.prev_px.G)
	vb := int8(px.B) - int8(e.prev_px.B)

	if vr >= -2 && vr < 2 && vg >= -2 && vg < 2 && vb >= -2 && vb < 2 {
		_, err = e.w.Write([]byte{
			uint8(QOI_OP_DIFF | ((vr + 2) << 4) | ((vg + 2) << 2) | (vb + 2)),
		})
		if err != nil {
			return err
		}
		e.prev_px = px
		return nil

	}

	vg_r := vr - vg
	vg_b := vb - vg

	if vg_r >= -8 && vg_r < 8 && vg >= -32 && vg < 32 && vg_b >= -8 && vg_b < 8 {
		_, err = e.w.Write([]byte{
			uint8(QOI_OP_LUMA | (int(vg) + 32)),
			uint8(((vg_r + 8) << 4) | (vg_b + 8)),
		})
		if err != nil {
			return err
		}
		e.prev_px = px
		return nil
	}

	_, err = e.w.Write([]byte{
		byte(QOI_OP_RGB), px.R, px.G, px.B,
	})
	if err != nil {
		return err
	}

	e.prev_px = px
	return nil
}
