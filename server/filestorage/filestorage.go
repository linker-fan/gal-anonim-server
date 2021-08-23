package filestorage

import (
	"log"

	"github.com/minio/minio-go"
)

type FileStorage struct {
	client *minio.Client
}

func NewFileStorage(endpoint, accessKeyID, secretAccessKey string, secure bool) (*FileStorage, error) {
	c, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
	if err != nil {
		log.Printf("NewFileStorage: minio.New failed: %v\n", err)
		return nil, err
	}

	return &FileStorage{
		client: c,
	}, nil
}
