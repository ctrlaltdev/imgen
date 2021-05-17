package img

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	svg "github.com/ajstarks/svgo"
)

func GenSVG(out io.Writer, w int, h int) error {
	canvas := svg.New(out)
	canvas.Start(w, h)

	var (
		rw     int
		rh     int
		angle  int
		stroke = w / 50
	)

	rand.Seed(time.Now().UnixNano())

	rw = rand.Intn(w/2-w/6) + w/6
	rh = rand.Intn(h/2-h/6) + h/6
	angle = rand.Intn(90) - 45

	canvas.RotateTranslate(0, 0, float64(angle))
	canvas.Rect(w/2-rw/2, h/2-rh/2, rw, rh, fmt.Sprintf("stroke:%s;fill:none;stroke-width:%dpx", "cyan", stroke))
	canvas.Gend()

	rw = rand.Intn(w/2-w/6) + w/6
	rh = rand.Intn(h/2-h/6) + h/6
	angle = rand.Intn(90) - 45

	canvas.RotateTranslate(0, 0, float64(angle))
	canvas.Rect(w/2-rw/2, h/2-rh/2, rw, rh, fmt.Sprintf("stroke:%s;fill:none;stroke-width:%dpx", "magenta", stroke))
	canvas.Gend()

	rw = rand.Intn(w/2-w/6) + w/6
	rh = rand.Intn(h/2-h/6) + h/6
	angle = rand.Intn(90) - 45

	canvas.RotateTranslate(0, 0, float64(angle))
	canvas.Rect(w/2-rw/2, h/2-rh/2, rw, rh, fmt.Sprintf("stroke:%s;fill:none;stroke-width:%dpx", "yellow", stroke))
	canvas.Gend()

	canvas.End()

	return nil
}
