package img

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"time"

	"golang.org/x/image/vector"
)

type Vertex struct {
	x float32
	y float32
}

func (v Vertex) Translate(x float32, y float32) Vertex {
	v.x = v.x + x
	v.y = v.y + y
	return v
}

func (v Vertex) Rotate(x float32, y float32, a float64) Vertex {
	vertex := Vertex{}
	vertex.x = (v.x-x)*float32(math.Cos(a)) - (v.y-y)*float32(math.Sin(a)) + x
	vertex.y = (v.x-x)*float32(math.Sin(a)) + (v.y-y)*float32(math.Cos(a)) + y
	v = vertex
	return v
}

type Vertices struct {
	v []Vertex
}

func (v Vertices) Translate(x float32, y float32) Vertices {
	for i, vertex := range v.v {
		v.v[i] = vertex.Translate(x, y)
	}
	return v
}

func (v Vertices) Rotate(x float32, y float32, a float64) Vertices {
	for i, vertex := range v.v {
		v.v[i] = vertex.Rotate(x, y, a)
	}
	return v
}

func CalcVertices(w int, h int, rw int, rh int, g int) Vertices {
	vertices := []Vertex{
		{x: float32(w/2 - rw/2 + g), y: float32(h/2 - rh/2 + g)},
		{x: float32(w/2 - rw/2 + rw - g), y: float32(h/2 - rh/2 + g)},
		{x: float32(w/2 - rw/2 + rw - g), y: float32(h/2 - rh/2 + rh - g)},
		{x: float32(w/2 - rw/2 + g), y: float32(h/2 - rh/2 + rh - g)},
	}
	return Vertices{vertices}
}

func DrawRectOutline(r *vector.Rasterizer, dst *image.RGBA, w int, h int, c color.RGBA, outV []Vertex, inV []Vertex) {
	r.Reset(w, h)
	r.MoveTo(outV[0].x, outV[0].y)
	r.LineTo(outV[1].x, outV[1].y)
	r.LineTo(outV[2].x, outV[2].y)
	r.LineTo(outV[3].x, outV[3].y)
	r.LineTo(outV[0].x, outV[0].y)

	r.LineTo(inV[0].x, inV[0].y)
	r.LineTo(inV[3].x, inV[3].y)
	r.LineTo(inV[2].x, inV[2].y)
	r.LineTo(inV[1].x, inV[1].y)
	r.LineTo(inV[0].x, inV[0].y)
	r.ClosePath()

	r.Draw(dst, dst.Bounds(), &image.Uniform{c}, image.Point{})
}

func GenVector(img *image.RGBA, w int, h int) error {
	r := vector.NewRasterizer(w, h)
	r.DrawOp = draw.Src

	var (
		rw     int
		rh     int
		angle  float64
		stroke = w / 50
		outV   Vertices
		inV    Vertices
		o      Vertex
		t      Vertex
		c      color.RGBA
	)

	rand.Seed(time.Now().UnixNano())

	// CYAN
	rw = rand.Intn(w/2-w/6) + w/6
	rh = rand.Intn(h/2-h/6) + h/6
	angle = float64(rand.Intn(90) - 45)

	o = Vertex{float32(w/2 + rand.Intn(rw) - rw/2), float32(h/2 + rand.Intn(rh) - rh/2)}
	t = Vertex{float32(rand.Intn(rw) - rw/2), float32(rand.Intn(rh) - rh/2)}

	outV = CalcVertices(w, h, rw, rh, 0).Rotate(o.x, o.y, angle).Translate(t.x, t.y)
	inV = CalcVertices(w, h, rw, rh, stroke).Rotate(o.x, o.y, angle).Translate(t.x, t.y)

	c = color.RGBA{0, 255, 255, 255}
	DrawRectOutline(r, img, w, h, c, outV.v, inV.v)

	// MAGENTA
	rw = rand.Intn(w/2-w/6) + w/6
	rh = rand.Intn(h/2-h/6) + h/6
	angle = float64(rand.Intn(90) - 45)

	o = Vertex{float32(w/2 + rand.Intn(rw) - rw/2), float32(h/2 + rand.Intn(rh) - rh/2)}
	t = Vertex{float32(rand.Intn(rw) - rw/2), float32(rand.Intn(rh) - rh/2)}

	outV = CalcVertices(w, h, rw, rh, 0).Rotate(o.x, o.y, angle).Translate(t.x, t.y)
	inV = CalcVertices(w, h, rw, rh, stroke).Rotate(o.x, o.y, angle).Translate(t.x, t.y)

	c = color.RGBA{255, 0, 255, 255}
	DrawRectOutline(r, img, w, h, c, outV.v, inV.v)

	// YELLOW
	rw = rand.Intn(w/2-w/6) + w/6
	rh = rand.Intn(h/2-h/6) + h/6
	angle = float64(rand.Intn(90) - 45)

	o = Vertex{float32(w/2 + rand.Intn(rw) - rw/2), float32(h/2 + rand.Intn(rh) - rh/2)}
	t = Vertex{float32(rand.Intn(rw) - rw/2), float32(rand.Intn(rh) - rh/2)}

	outV = CalcVertices(w, h, rw, rh, 0).Rotate(o.x, o.y, angle).Translate(t.x, t.y)
	inV = CalcVertices(w, h, rw, rh, stroke).Rotate(o.x, o.y, angle).Translate(t.x, t.y)

	c = color.RGBA{255, 255, 0, 255}
	DrawRectOutline(r, img, w, h, c, outV.v, inV.v)

	return nil
}
