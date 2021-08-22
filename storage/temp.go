package storage

import (
	"fmt"
	"io/ioutil"
	"log"
)

// MakeTempFile creates a temp file in /tmp for processing.
func MakeTempFile(fileBytes []byte, fileName string) (string, error) {
	tmp, err := ioutil.TempFile("", fileName+"")
	if err != nil {
		return "", fmt.Errorf("unable to create temp file for image processing: %w", err)
	}
	defer tmp.Close()

	tmp.Write(fileBytes)

	log.Printf("created temp file %s\n", tmp.Name())

	return tmp.Name(), nil
}
