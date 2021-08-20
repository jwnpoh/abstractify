package app

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

type colorAt struct {
	x, y int
	rgb
}

func Fudge(inFile string) (string, error) {
	log.Printf("processing %s now...\n", inFile)
	srcImg, err := gg.LoadImage(inFile)
	if err != nil {
		return "", fmt.Errorf("oops...something went wrong. image file was not successfully decoded: %w", err)
	}

	s := newSketch(srcImg)

	sketchIt(s)

	base := filepath.Base(inFile)
	fileName := strings.TrimSuffix(base, filepath.Ext(base))
	outputFileNameBase := fileName + "-" + "abstractified.png"
	outputFileName := filepath.Join("tmp", outputFileNameBase)

	err = gg.SavePNG(outputFileName, s.output())
	log.Printf("successfully generated %s\n", outputFileName)

	return outputFileName, nil
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
