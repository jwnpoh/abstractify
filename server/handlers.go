package server

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jwnpoh/abstractify/app"
	"github.com/jwnpoh/abstractify/logger"
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

	opts := parseOptions(r)

	outFileName, err := app.Fudge(opts, fileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	since := time.Since(now)
	log.Printf("took %v\n", since)

	err = logIt(header, outFileName, since)
	if err != nil {
		log.Printf("something went wrong with the logging: %v", err)
	}
	log.Printf("log updated on %v\n", time.Now())

	data := struct {
		FileName string
	}{
		FileName: outFileName,
	}
	tpl.ExecuteTemplate(w, "download.html", data)
}

func download(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "oops...something went wrong with the file download...please try again", http.StatusBadRequest)
	}

	filename := r.Form.Get("download")
	log.Printf("retrieving %s from cloud storage", filename)

	item, err := storage.DownloadFromCloudStorage(filename)
	if err != nil {
		log.Printf("error retrieving %s from cloud storage: %v", filename, err)
		http.Error(w, "oops...something went wrong with the file download...please try again", http.StatusBadRequest)
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", item.ContentType)
	w.Header().Set("Content-Length", item.Size)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.WriteHeader(http.StatusOK)
	w.Write(item.Content)
	log.Printf("delivered %s", filename)

	log.Printf("deleting %s from cloud storage now...", filename)
	err = storage.Delete(filename)
	if err != nil {
		log.Printf("unable to delete %s from cloud storage:", filename)
	}
	log.Printf("deleted %s from cloud storage.", filename)
	log.Println("Exiting...")
	log.Println(strings.Repeat("-", 20))
}

func validateUpload(header *multipart.FileHeader) error {
	mimetype := header.Header.Get("Content-Type")
	if !strings.Contains(mimetype, "image") {
		return fmt.Errorf("please upload only a JPEG or PNG image")
	}

	if header.Size > 4*megabyte {
		return fmt.Errorf("please upload files no larger than 3mb")
	}

	return nil
}

func parseOptions(r *http.Request) *app.Opts {
	var opts app.Opts

	opts.Shape = r.FormValue("shape")

	size, err := strconv.Atoi(r.FormValue("size"))
	if err != nil {
		size = 1
	}
	opts.Size = size

	sizeOpt := r.FormValue("randomSizeOpt")
	if sizeOpt == "yes" {
		opts.RandomSize = true
	} else {
		opts.RandomSize = false
	}

  log.Println(opts)

	return &opts
}

func logIt(header *multipart.FileHeader, fileName string, timeSince time.Duration) error {
	entry := logger.NewEntry()

	entry.LogFileName(header.Filename)
	entry.LogFileSize(int(header.Size) / kilobyte)

	filePath := "/tmp/" + fileName
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("unable to access output file %s info for logging: %w", filePath, err)
	}

	entry.LogOutPutSize(int(fileInfo.Size()) / kilobyte)
	entry.LogProcessTime(timeSince)
	loc := time.FixedZone("UTC+8", 8*60*60)
	if err != nil {
		return fmt.Errorf("unable to get local time in Singapore: %v", err)
	}

	t := time.Now().In(loc)
	entry.LogTime(t)

	entries, err := logger.LoadLogs("logs.json")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	*entries = append(*entries, *entry)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = logger.SubmitLogs(entries)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
