package storage

import (
	"context"
	"io/ioutil"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func CreateClient(ctx context.Context) (*storage.Client, error) {
	return storage.NewClient(ctx, option.WithCredentialsJSON(Key))
}

type StorageOp struct {
	ctx context.Context
	rc  *storage.BucketHandle
}

// Create Google Cloud Storage operation handler.
func NewStorageOp(ctx context.Context, client storage.Client, bucketName string) *StorageOp {
	rc := client.Bucket(bucketName)

	return &StorageOp{
		ctx: ctx,
		rc:  rc,
	}
}

// Create Storage object.
// Returns an ObjectHandle that looks like a directory and file to manipulate.
func (s *StorageOp) Object(dirs []string, fileName string) *storage.ObjectHandle {
	path := append(dirs, fileName)
	return s.rc.Object(strings.Join(path, "/"))
}

// Check if file exists.
// Exist if true, false not.
func (s *StorageOp) FileExist(dirs []string, fileName string) (bool, error) {
	object := s.Object(dirs, fileName)
	_, err := object.Attrs(s.ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Read file.
func (s *StorageOp) ReadFile(dirs []string, fileName string) ([]byte, error) {
	object := s.Object(dirs, fileName)
	reader, err := object.NewReader(s.ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Write file
func (s *StorageOp) WriteFile(dirs []string, fileName string, body []byte) error {
	object := s.Object(dirs, fileName)
	writer := object.NewWriter(s.ctx)

	_, err := writer.Write(body)
	if err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

// Delete files
func (s *StorageOp) Delete(prefix string) error {
	objects := s.rc.Objects(s.ctx, &storage.Query{
		Prefix: prefix,
	})

	for {
		attrs, err := objects.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if err := s.rc.Object(attrs.Name).Delete(s.ctx); err != nil {
			return err
		}
	}
	return nil
}

// disable to versioning.
func (s *StorageOp) DisableVersioning() error {
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		VersioningEnabled: false,
	}

	if _, err := s.rc.Update(s.ctx, bucketAttrsToUpdate); err != nil {
		return err
	}
	return nil
}
