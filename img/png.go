package img

import (
	"image"
	"image/png"
	"io"
)

func GenPNG(out io.Writer, w int, h int) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	GenVector(img, w, h)

	err := png.Encode(out, img)

	return err
}
