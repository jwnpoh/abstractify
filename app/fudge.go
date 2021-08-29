package app

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/jwnpoh/abstractify/storage"
	"github.com/nfnt/resize"
)

type Opts struct {
  Shape string
  Size int
  RandomSize bool
}

type colorAt struct {
	x, y int
	rgb
}

// Fudge is the main entry point to the image processing function of the app.
func Fudge(opts *Opts, inFile string) (string, error) {
  log.Printf("app.Fudge: opts received: %v", *opts)

	log.Printf("processing %s now...\n", inFile)
	srcImg, err := gg.LoadImage(inFile)
	if err != nil {
		return "", fmt.Errorf("oops...something went wrong. image file was not successfully decoded: %w", err)
	}

  resizedImg := resizeImage(srcImg)

	s := newSketch(resizedImg, opts)

	sketchIt(s)

  outputFileName, outputFileNameBase, err := processFileNames(inFile)
  if err != nil {
		log.Printf("problem saving generated image: %v", err)
		return "", fmt.Errorf("something went wrong with the image generation. Please try again: %w", err)
  }

	err = gg.SavePNG(outputFileName, s.output())
	if err != nil {
		log.Printf("problem saving generated image: %v", err)
		return "", fmt.Errorf("something went wrong with the image generation. Please try again: %w", err)
	}
	log.Printf("successfully generated %s\n", outputFileName)

  err = uploadProcessedImage(outputFileName, outputFileNameBase)

	return outputFileNameBase, nil
}

func resizeImage(src image.Image) image.Image {
	var resizedImg image.Image
	switch {
	case src.Bounds().Max.X > src.Bounds().Max.Y:
		resizedImg = resize.Resize(1920, 0, src, resize.Bilinear)
	case src.Bounds().Max.X < src.Bounds().Max.Y:
		resizedImg = resize.Resize(0, 1350, src, resize.Bilinear)
	default:
		resizedImg = resize.Resize(0, 1350, src, resize.Bilinear)
	}
  return resizedImg
}

func uploadProcessedImage(outputFileName, outputFileNameBase string) error {
	file, err := os.Open(outputFileName)
	if err != nil {
		log.Printf("app.Fudge - unable to open saved file for uploading: %v", err)
		return fmt.Errorf("unable to access saved file: %w", err)
	}
	defer file.Close()

	log.Printf("uploading %s now...\n", outputFileNameBase)
	err = storage.Upload(file, outputFileNameBase)
	if err != nil {
		log.Printf("app.Fudge - unable to upload file to cloud storage: %v", err)
		return fmt.Errorf("unable to upload to cloud storage: %w", err)
	}
  return nil
}

func processFileNames(inFile string) (string, string, error) {
	base := filepath.Base(inFile)
	fileName := strings.TrimSuffix(base, filepath.Ext(base))
	outputFileNameBase := fileName + "-" + "abstractified.png"
	outputFileName := filepath.Join("/tmp", outputFileNameBase)

  return outputFileName, outputFileNameBase, nil
}
