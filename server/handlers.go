package server

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/jwnpoh/abstractify/app"
	"github.com/jwnpoh/abstractify/storage"
)

const (
	megabyte = 1024 * 1024
	kilobyte = 1024
)

func index(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, "oops...something went wrong with the upload"+fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = validateUpload(header)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "oops...something went wrong with the upload"+fmt.Sprint(err), http.StatusBadRequest)
		return
	}

	fileName, err := storage.MakeTempFile(fileBytes, header.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	log.Printf("file size is %vkb\n", header.Size/kilobyte)

	now := time.Now()

	outFileName, err := app.Fudge(fileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	since := time.Since(now)
	log.Printf("took %v\n", since)
	log.Println(strings.Repeat("-", 20))

	data := struct {
		FileName string
	}{
		FileName: outFileName,
	}
	tpl.ExecuteTemplate(w, "download.html", data)
}

func validateUpload(header *multipart.FileHeader) error {
	mimetype := header.Header.Get("Content-Type")
	if !strings.Contains(mimetype, "image") {
		return fmt.Errorf("please upload only a JPEG or PNG image")
	}

	if header.Size > 4*megabyte {
		return fmt.Errorf("please upload files no larger than 4mb")
	}

	return nil
}

func download(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "oops...something went wrong with the file download...please try again", http.StatusBadRequest)
	}

	filename := r.Form.Get("download")
	filenamebase := filepath.Base(filename)
	w.Header().Set("Content-Disposition", "attachment; filename="+filenamebase)
	http.ServeFile(w, r, filename)
	log.Printf("delivered %s\n", filename)
	log.Println(strings.Repeat("-", 20))
}
