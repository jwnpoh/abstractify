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
	colorSlice := make([]rgb, 0, radius*4)

	radius /= 2

	var color rgb

	for i := 0; i <= radius; i++ {
		if x+i >= srcImg.Bounds().Dx() {
			break
		}
		r, g, b := rgb255(srcImg.At(x+i, y))
		color.r, color.g, color.b = r, g, b
		colorSlice = append(colorSlice, color)
	}

	for i := 0; i <= radius; i++ {
		if y+i >= srcImg.Bounds().Dy() {
			break
		}
		r, g, b := rgb255(srcImg.At(x, y+i))
		color.r, color.g, color.b = r, g, b
		colorSlice = append(colorSlice, color)
	}

	for i := 0; i <= radius; i++ {
		if x-i <= srcImg.Bounds().Min.X {
			break
		}
		r, g, b := rgb255(srcImg.At(x-i, y))
		color.r, color.g, color.b = r, g, b
		colorSlice = append(colorSlice, color)
	}

	for i := 0; i <= radius; i++ {
		if y-i <= srcImg.Bounds().Min.Y {
			break
		}
		r, g, b := rgb255(srcImg.At(x, y-i))
		color.r, color.g, color.b = r, g, b
		colorSlice = append(colorSlice, color)
	}

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
	source                image.Image
	dc                    *gg.Context
	radius                float64
	colorToSketch         rgb
	cycleCount            int
}

func newSketch(src image.Image) *sketch {
	var s sketch

	bounds := src.Bounds()
	s.destWidth, s.destHeight = bounds.Max.X, bounds.Max.Y

	s.cycleCount = 150

	canvas := gg.NewContext(s.destWidth, s.destHeight)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(s.destWidth), float64(s.destHeight))
	canvas.FillPreserve()

	s.source = src
	s.dc = canvas

	return &s
}

func sketchIt(s *sketch) {
	s.radius = float64(s.destWidth) / float64(s.cycleCount)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < s.cycleCount; i++ {
		for x := 0; x < s.destWidth; x++ {
			y := rand.Intn(s.destHeight)
			colorSlice := getRGBSlice(s.source, x, y, int(s.radius))
			s.colorToSketch = averageRGB(*colorSlice)
			s.update(x, y)
		}
	}
}

func (s *sketch) update(x, y int) {
	r, g, b := s.colorToSketch.r, s.colorToSketch.g, s.colorToSketch.b

	rand.Seed(time.Now().UnixNano())

	a := rand.Intn(80)

	radius := rand.Float64() * float64(rand.Intn(3)) * s.radius
	switch {
	case s.destWidth < 1000:
		radius *= 2
	case s.destWidth < 2000:
		radius *= 1.5
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
