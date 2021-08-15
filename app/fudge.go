package app

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
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
	const cycleCount = 150

	fileDL, err := storage.Download(inFile)
  if err != nil {
    return "", fmt.Errorf("oops...something went wrong with the file upload...please try again: %w", err)
  }

	log.Printf("received %s and processing now\n", fileDL)
	srcImg, err := gg.LoadImage(fileDL)
	if err != nil {
		return "", fmt.Errorf("oops...something went wrong. image file was not successfully decoded: %w", err)
	}

	s := newSketch(srcImg)
	s.cycleCount = cycleCount

	sketchIt(s)

	base := filepath.Base(fileDL)
	fileName := strings.TrimSuffix(base, filepath.Ext(base))
	outputFileNameBase := fileName + "-" + "abstractified.png"
	outputFileName := filepath.Join("tmp", outputFileNameBase)

	err = gg.SavePNG(outputFileName, s.output())

	file, err := os.Open(outputFileName)
	if err != nil {
		return "", fmt.Errorf("unable to create tmp file for processing: %w", err)
	}
	defer file.Close()

	var w http.ResponseWriter
	err = storage.Upload(w, file, outputFileNameBase)
	if err != nil {
		return "", fmt.Errorf("oops...unable to create abstractified image...%w", err)
	}

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
