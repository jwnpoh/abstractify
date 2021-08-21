package app

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/jwnpoh/abstractify/storage"
)

type colorAt struct {
	x, y int
	rgb
}

func Fudge(inFile string) (string, error) {
	log.Printf("processing %s now...\n", inFile)

	tmpFile, err := storage.Download(inFile)

	srcImg, err := gg.LoadImage(tmpFile)
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

	file, err := os.Open(outputFileName)

	storage.Upload(file, outputFileNameBase)

	return outputFileName, nil
}

func sketchIt(s *sketch) {
	s.radius = float64(s.destWidth) / float64(s.cycleCount)
	rand.Seed(time.Now().UnixNano())

	inc := 2
	if s.destWidth < 1000 || s.destHeight < 1000 {
		inc = 1
	}
	for i := 0; i < s.cycleCount; i++ {
		for x := 0; x < s.destWidth; x += inc {
			y := rand.Intn(s.destHeight)
			colorSlice := getRGBSlice(s.source, x, y, int(s.radius))
			s.colorToSketch = averageRGB(*colorSlice)
			s.update(x, y)
		}
	}
}
