package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jwnpoh/abstractify/storage"
)

// Entry represents a single log entry.
type Entry struct {
	TimeOfEntry string `json:"timeOfEntry"`
	FileName    string `json:"fileName"`
	FileSize    string `json:"fileSize"`
	OutputSize  string `json:"outputSize"`
	ProcessTime string `json:"processTime"`
}

// NewEntry creates a new entry for logging.
func NewEntry() *Entry {
	var e Entry
	return &e
}

// LogTime logs the time of the entry.
func (e *Entry) LogTime(t time.Time) {
	e.TimeOfEntry = t.Format("Mon 2 Jan 2006, 15:04:05")
}

// LogFileName logs the name of the uploaded file.
func (e *Entry) LogFileName(s string) {
	e.FileName = s
}

// LogFileSize logs the size of the uploaded file.
func (e *Entry) LogFileSize(i int) {
	s := strconv.Itoa(i)
	e.FileSize = s + "kb"
}

// LogOutPutSize logs the size of the output file.
func (e *Entry) LogOutPutSize(i int) {
	s := strconv.Itoa(i)
	e.OutputSize = s + "kb"
}

// LogProcessTime logs the time taken to generate the abstractified image.
func (e *Entry) LogProcessTime(d time.Duration) {
	e.ProcessTime = d.String()
}

// Entries represents a slice of multiple entries for marshaling into the json log file.
type Entries []Entry

// SubmitLogs creates a json log file from a slice of Entry structs.
func SubmitLogs(xe *Entries) error {
	jsonData, err := json.MarshalIndent(&xe, "", "\t")
	if err != nil {
		return fmt.Errorf("couldn't submit logs: %w", err)
	}

	data := bytes.NewReader(jsonData)
	storage.Upload(data, "logs.json")

	return nil
}

// LoadLogs downloads the persisted log file containing all previous entries and returns a pointer to Entries, so that we can append any new Entry created in the current instance.
func LoadLogs(logFile string) (*Entries, error) {
	logFileObject, err := storage.DownloadFromCloudStorage(logFile)
	if err != nil {
		if err := createNewLogFile(); err != nil {
			return nil, fmt.Errorf("unable to create new log file: %w", err)
		}
		logFileObject, _ = storage.DownloadFromCloudStorage(logFile)
	}

	logBytes := logFileObject.Content

	xe := make(Entries, 0)
	err = json.Unmarshal(logBytes, &xe)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal json from log file: %w", err)
	}

	return &xe, nil
}

func createNewLogFile() error {
	b := make(Entries, 0, 0)

	e := NewEntry()
	e.LogFileName("dummy")
	e.LogFileSize(0)
	e.LogOutPutSize(0)
	e.LogProcessTime(time.Minute)
	e.LogTime(time.Now())

	b = append(b, *e)

	err := SubmitLogs(&b)
	if err != nil {
		return err
	}
	return nil
}
