package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/jwnpoh/abstractify/app"
	"github.com/jwnpoh/abstractify/storage"
)

const kilobyte = 1024

// Entry represents a single log entry.
type Entry struct {
	Index       int      `json:"index"`
	TimeOfEntry string   `json:"timeOfEntry"`
	FileName    string   `json:"fileName"`
	FileSize    string   `json:"fileSize"`
	OutputSize  string   `json:"outputSize"`
	ProcessTime string   `json:"processTime"`
	Opts        app.Opts `json:"opts"`
}

// newEntry creates a new entry for logging.
func newEntry(index int) *Entry {
	var e Entry
	e.Index = index
	return &e
}

// LogTime logs the time of the entry.
func (e *Entry) logTime() {
	loc := time.FixedZone("UTC+8", 8*60*60)

	t := time.Now().In(loc)
	e.TimeOfEntry = t.Format("Mon 2 Jan 2006, 15:04:05")
}

// LogProcessTime logs the time taken to generate the abstractified image.
func (e *Entry) logProcessTime(d time.Duration) {
	e.ProcessTime = d.String()
}

func (e *Entry) logOpts(opts *app.Opts) {
	e.Opts = *opts
}

// Entries represents a slice of multiple entries for marshaling into the json log file.
type Entries []Entry

// submitLogs creates a json log file from a slice of Entry structs.
func submitLogs(xe *Entries) error {
	jsonData, err := json.MarshalIndent(&xe, "", "\t")
	if err != nil {
		return fmt.Errorf("couldn't submit logs: %w", err)
	}

	data := bytes.NewReader(jsonData)
	storage.Upload(data, "logs.json")

	return nil
}

// loadLogs downloads the persisted log file containing all previous entries and returns a pointer to Entries, so that we can append any new Entry created in the current instance.
func loadLogs(logFile string) (*Entries, int, error) {
	logFileObject, err := storage.DownloadFromCloudStorage(logFile)
	if err != nil {
		if err := createNewLogFile(); err != nil {
			return nil, 0, fmt.Errorf("unable to create new log file: %w", err)
		}
		logFileObject, _ = storage.DownloadFromCloudStorage(logFile)
	}

	logBytes := logFileObject.Content

	xe := make(Entries, 0)
	err = json.Unmarshal(logBytes, &xe)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to unmarshal json from log file: %w", err)
	}

	return &xe, len(xe), nil
}

func createNewLogFile() error {
	b := make(Entries, 0, 0)

	e := newEntry(0)
	e.logFileInfo("dummy", 0)
	e.logOutput("dummy")
	e.logProcessTime(time.Minute)
	e.logTime()
	e.logOpts(&app.Opts{Shape: "Random", Size: 0, RandomSize: false})

	b = append(b, *e)

	err := submitLogs(&b)
	if err != nil {
		return err
	}
	return nil
}

// LogInstance logs an entry from the current instance.
func LogInstance(header *multipart.FileHeader, fileName string, timeSince time.Duration, opts *app.Opts) error {
	entries, nextIndex, err := loadLogs("logs.json")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	entry := newEntry(nextIndex)

	entry.logFileInfo(header.Filename, int(header.Size))

	err = entry.logOutput(fileName)
	if err != nil {
		return err
	}

	entry.logProcessTime(timeSince)

	entry.logTime()

	entry.logOpts(opts)

	*entries = append(*entries, *entry)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = submitLogs(entries)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	log.Printf("log updated on %v\n", entry.TimeOfEntry)
	log.Print(*entry)

	return nil
}

func (e *Entry) logFileInfo(fileName string, fileSize int) {
	e.FileName = fileName
	s := strconv.Itoa(fileSize / kilobyte)
	e.FileSize = s + "kb"
}

func (e *Entry) logOutput(fileName string) error {
	if fileName == "dummy" {
		e.OutputSize = "0kb"
		return nil
	}

	filePath := "/tmp/" + fileName
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("unable to access output file %s info for logging: %w", filePath, err)
	}

	s := strconv.Itoa(int(fileInfo.Size()) / kilobyte)
	e.OutputSize = s + "kb"
	return nil
}
