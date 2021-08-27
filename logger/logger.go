package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jwnpoh/abstractify/storage"
)

type Entry struct {
	TimeOfEntry time.Time     `json:"timeOfEntry"`
	FileName    string        `json:"fileName"`
	FileSize    string        `json:"fileSize"`
	OutputSize  string        `json:"outputSize"`
	ProcessTime time.Duration `json:"processTime"`
}

func NewEntry() *Entry {
	var e Entry
	return &e
}

func (e *Entry) LogTime(t time.Time) {
	e.TimeOfEntry = t
}

func (e *Entry) LogFileName(s string) {
	e.FileName = s
}

func (e *Entry) LogFileSize(s string) {
	e.FileSize = s
}

func (e *Entry) LogOutPutSize(s string) {
	e.OutputSize = s
}

func (e *Entry) LogProcessTime(d time.Duration) {
	e.ProcessTime = d
}

type Entries []Entry

func SubmitLogs(xe *Entries) error {
	jsonData, err := json.MarshalIndent(&xe, "", "\t")
	if err != nil {
		return fmt.Errorf("couldn't submit logs: %w", err)
	}
  log.Printf("SubmitLogs: %v", jsonData)

	data := bytes.NewReader(jsonData)

	storage.Upload(data, "logs.json")

	return nil
}

func LoadLogs(logFile string) (*Entries, error) {
	logFileObject, err := storage.DownloadFromCloudStorage(logFile)
	if err != nil {
		if err := createNewLogFile(); err != nil {
			return nil, fmt.Errorf("unable to create new log file: %w", err)
		}
	}

	logBytes := logFileObject.Content

	xe := make(Entries, 0)
	err = json.Unmarshal(logBytes, &xe)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal json from log file: %w", err)
	}

  log.Printf("LoadLogs: %v", xe)
	return &xe, nil
}

func createNewLogFile() error {
	b := make(Entries, 0, 0)

	e := NewEntry()
	e.LogFileName("dummy")
	e.LogFileSize("dummy")
	e.LogOutPutSize("dummy")
	e.LogProcessTime(0)
	e.LogTime(time.Now())

	b = append(b, *e)

	err := SubmitLogs(&b)
	if err != nil {
		return err
	}
	return nil
}
