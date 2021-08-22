package storage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var bucket = struct {
	bucketName  string
	credentials string
}{
	bucketName:  os.Getenv("BUCKET"),
	credentials: os.Getenv("CREDENTIALS"),
}

// Upload uploads a file to Cloud Storage.
func Upload(file io.Reader, fileName string) error {
	object := fileName
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(bucket.credentials)))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket.bucketName).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

// Download downloads a file from Cloud Storage and saves it to /tmp.
func Download(object string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(bucket.credentials)))
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket.bucketName).Object(object).NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	tmpFile, err := os.Create(filepath.Join("tmp", object))
	if err != nil {
		return "", fmt.Errorf("unable to create tmp file %s: %w", object, err)
	}
	defer tmpFile.Close()

	// write to tmpFile
	tmpFile.Write(data)

	return tmpFile.Name(), nil
}

// Delete deletes a file in Cloud Storage.
func Delete(object string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(bucket.credentials)))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket.bucketName).Object(object)
	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}
	return nil
}

// DownloadObject represents the downloaded object from Cloud Storage that will be served in the browser.
type DownloadObject struct {
	ContentType string
	Size        string
	Content     []byte
}

// DownloadFromCloudStorage downloads a file from Cloud Storage ready to be served directly in the browser.
func DownloadFromCloudStorage(object string) (DownloadObject, error) {
	var item DownloadObject

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(bucket.credentials)))
	if err != nil {
		return item, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket.bucketName).Object(object).NewReader(ctx)
	if err != nil {
		return item, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return item, fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	item.ContentType = rc.ContentType()
	item.Size = strconv.Itoa(int(rc.Size()))
	item.Content = data

	return item, nil
}
