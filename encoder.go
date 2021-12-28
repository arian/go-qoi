package go_qoi

import (
	"image"
	"image/color"
	"io"
)

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

	seen_pixels := make([]color.NRGBA, 64)
	prev_px := color.NRGBA{A: 255}
	run := 0

	for i := 0; i < l; i++ {
		y := i / width
		x := i % width

		px := colorToNRGBA(img.At(x, y))

		if px == prev_px {
			run++
			if run == 62 || i == (l-1) {
				_, err = w.Write([]byte{uint8(QOI_OP_RUN | run - 1)})
				if err != nil {
					return err
				}
				run = 0
			}
		} else {

			if run > 0 {
				_, err = w.Write([]byte{uint8(QOI_OP_RUN | (run - 1))})
				if err != nil {
					return err
				}
				run = 0
			}

			index_pos := indexPositionHash(px)

			if seen_pixels[index_pos] == px {
				_, err = w.Write([]byte{QOI_OP_INDEX | index_pos})
				if err != nil {
					return err
				}
			} else {
				seen_pixels[index_pos] = px

				if px.A == prev_px.A {
					vr := int8(px.R) - int8(prev_px.R)
					vg := int8(px.G) - int8(prev_px.G)
					vb := int8(px.B) - int8(prev_px.B)

					vg_r := vr - vg
					vg_b := vb - vg

					if vr >= -2 && vr < 2 && vg >= -2 && vg < 2 && vb >= -2 && vb < 2 {
						_, err = w.Write([]byte{
							uint8(QOI_OP_DIFF | (vr+2)<<4 | (vg+2)<<2 | (vb + 2)),
						})
						if err != nil {
							return err
						}

					} else if vg_r > -9 && vg_r < 8 &&
						vg > -33 && vg < 32 &&
						vg_b > -9 && vg_b < 8 {

						_, err = w.Write([]byte{
							uint8(QOI_OP_LUMA | (int(vg) + 32)),
							uint8((vg_r+8)<<4 | (vg_b + 8)),
						})
						if err != nil {
							return err
						}

					} else {
						_, err := w.Write([]byte{byte(QOI_OP_RGB), px.R, px.G, px.B})
						if err != nil {
							return err
						}
					}

				} else {
					_, err := w.Write([]byte{QOI_OP_RGBA, px.R, px.G, px.B, px.A})
					if err != nil {
						return err
					}
				}
			}
		}

		prev_px = px
	}

	w.Write(QOI_END_PADDING[:])

	return nil
}
