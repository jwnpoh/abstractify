package app

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

type rgb struct {
	r, g, b int
}

func getRGBSlice(srcImg image.Image, x, y, radius int) *[]rgb {
	colorSlice := make([]rgb, 0, radius*8)

	for i := 0; i <= radius; i++ {
		if x+i >= srcImg.Bounds().Dx() {
			break
		}
		var color rgb
		r, g, b := rgb255(srcImg.At(x+i, y))
		color.r, color.g, color.b = r, g, b
		colorSlice = append(colorSlice, color)
	}

	for i := 0; i <= radius; i++ {
		if y+i >= srcImg.Bounds().Dy() {
			break
		}
		var color rgb
		r, g, b := rgb255(srcImg.At(x, y+i+i+i-1))
		color.r, color.g, color.b = r, g, b
		colorSlice = append(colorSlice, color)
	}

	for i := 0; i <= radius; i++ {
		if x-i <= srcImg.Bounds().Min.X {
			break
		}
		var color rgb
		r, g, b := rgb255(srcImg.At(x-i, y))
		color.r, color.g, color.b = r, g, b
		colorSlice = append(colorSlice, color)
	}

	for i := 0; i <= radius; i++ {
		if y-i <= srcImg.Bounds().Min.Y {
			break
		}
		var color rgb
		r, g, b := rgb255(srcImg.At(x, y-i-i-i+1))
		color.r, color.g, color.b = r, g, b
		colorSlice = append(colorSlice, color)
	}

	// 	for i := 0; i <= radius; i++ {
	// 		if x-i <= srcImg.Bounds().Min.X || y-i <= srcImg.Bounds().Min.Y {
	// 			break
	// 		}
	// 		var color rgb
	// 		r, g, b := rgb255(srcImg.At(x-i, y-i))
	// 		color.r, color.g, color.b = r, g, b
	// 		colorSlice = append(colorSlice, color)
	// 	}
	//
	// 	for i := 0; i <= radius; i++ {
	// 		if y-i <= srcImg.Bounds().Min.Y || x+i >= srcImg.Bounds().Dx() {
	// 			break
	// 		}
	// 		var color rgb
	// 		r, g, b := rgb255(srcImg.At(x+i, y-i))
	// 		color.r, color.g, color.b = r, g, b
	// 		colorSlice = append(colorSlice, color)
	// 	}
	//
	// 	for i := 0; i <= radius; i++ {
	// 		if x-i <= srcImg.Bounds().Min.X || y+i >= srcImg.Bounds().Dy() {
	// 			break
	// 		}
	// 		var color rgb
	// 		r, g, b := rgb255(srcImg.At(x-i, y+i))
	// 		color.r, color.g, color.b = r, g, b
	// 		colorSlice = append(colorSlice, color)
	// 	}
	//
	// 	for i := 0; i <= radius; i++ {
	// 		if x+i >= srcImg.Bounds().Dx() || y+i >= srcImg.Bounds().Dy() {
	// 			break
	// 		}
	// 		var color rgb
	// 		r, g, b := rgb255(srcImg.At(x+i, y+i))
	// 		color.r, color.g, color.b = r, g, b
	// 		colorSlice = append(colorSlice, color)
	// 	}

	return &colorSlice
}

func averageRGB(xr []rgb) rgb {
	var avgRGB rgb
	var r, g, b, count int
	count = len(xr)

	for _, j := range xr {
		r += j.r
		g += j.g
		b += j.b
	}

	avgRGB.r = r / count
	avgRGB.g = g / count
	avgRGB.b = b / count

	return avgRGB
}

type sketch struct {
	destWidth, destHeight int
	edgeCount             int
	source                image.Image
	dc                    *gg.Context
	radius                float64
	colorToSketch         rgb
	cycleCount            int
}

// newSketch returns a *sketch to work with.
func newSketch(src image.Image) *sketch {
	var s sketch

	bounds := src.Bounds()
	s.destWidth, s.destHeight = bounds.Max.X, bounds.Max.Y

	s.cycleCount = 150

	canvas := gg.NewContext(s.destWidth, s.destHeight)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.destWidth), float64(s.destHeight))
	canvas.FillPreserve()

	s.source = src
	s.dc = canvas

	return &s
}

func (s *sketch) update(x, y int) {
	r, g, b := s.colorToSketch.r, s.colorToSketch.g, s.colorToSketch.b

	rand.Seed(time.Now().UnixNano())

	a := rand.Intn(100)

	radius := rand.Float64() * float64(rand.Intn(3)) * s.radius
	if s.destWidth < 1000 || s.destHeight < 1000 {
		radius *= 3
	}

	s.dc.SetRGBA255(r, g, b, a)
	s.dc.DrawRegularPolygon(6, float64(x), float64(y), radius, rand.Float64())
	s.dc.FillPreserve()
	s.dc.Stroke()
}

func (s *sketch) output() image.Image {
	return s.dc.Image()
}

func rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 257), int(g0 / 257), int(b0 / 257)
}
