package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/jwnpoh/abstractify/storage"
	"github.com/nfnt/resize"
)

type colorAt struct {
	x, y int
	rgb
}

// Fudge is the main entry point to the image processing function of the app.
func Fudge(inFile string) (string, error) {
	log.Printf("processing %s now...\n", inFile)
	srcImg, err := gg.LoadImage(inFile)
	if err != nil {
		return "", fmt.Errorf("oops...something went wrong. image file was not successfully decoded: %w", err)
	}

	resizedImg := resize.Resize(1080, 0, srcImg, resize.Bilinear)

	s := newSketch(resizedImg)

	sketchIt(s)

	base := filepath.Base(inFile)
	fileName := strings.TrimSuffix(base, filepath.Ext(base))
	outputFileNameBase := fileName + "-" + "abstractified.png"
	outputFileName := filepath.Join("/tmp", outputFileNameBase)

	err = gg.SavePNG(outputFileName, s.output())
	if err != nil {
		log.Printf("problem saving generated image: %v", err)
		return "", fmt.Errorf("something went wrong with the image generation. Please try again: %w", err)
	}
	log.Printf("successfully generated %s\n", outputFileName)

	file, err := os.Open(outputFileName)
	if err != nil {
		log.Printf("app.Fudge - unable to open saved file for uploading: %v", err)
		return "", fmt.Errorf("unable to access saved file: %w", err)
	}
	defer file.Close()

	log.Printf("uploading %s now...\n", outputFileNameBase)
	err = storage.Upload(file, outputFileNameBase)
	if err != nil {
		log.Printf("app.Fudge - unable to upload file to cloud storage: %v", err)
		return "", fmt.Errorf("unable to upload to cloud storage: %w", err)
	}

	return outputFileNameBase, nil
}
