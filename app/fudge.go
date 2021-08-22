package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/jwnpoh/abstractify/storage"
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
	outputFileName := filepath.Join("/tmp", outputFileNameBase)

	err = gg.SavePNG(outputFileName, s.output())
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
