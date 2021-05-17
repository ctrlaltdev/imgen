package img

import (
	"image"
	"image/jpeg"
	"io"
)

func GenJPG(out io.Writer, w int, h int) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	GenVector(img, w, h)

	opts := &jpeg.Options{Quality: 100}
	err := jpeg.Encode(out, img, opts)

	return err
}
