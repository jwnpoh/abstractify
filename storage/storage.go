package storage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

func Upload(w io.Writer, file io.Reader, fileName string) error {
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
