package filestorage

import (
	"bytes"
	"errors"
	"log"

	"github.com/minio/minio-go"
)

type FileStorage struct {
	client        *minio.Client
	MaxUploadSize int64
}

func NewFileStorage(endpoint, accessKeyID, secretAccessKey string, secure bool, maxUploadSize int64) (*FileStorage, error) {
	c, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
	if err != nil {
		log.Printf("NewFileStorage: minio.New failed: %v\n", err)
		return nil, err
	}

	return &FileStorage{
		client: c,
	}, nil
}

func (fs *FileStorage) StoreFile(bucketName string, file *bytes.Buffer, filename string) error {
	bucketExists, err := fs.CheckIfBucketExists(bucketName)
	if err != nil {
		return err
	}

	if !bucketExists {
		return errors.New("Bucket does not exist")
	}

	_, err = fs.client.PutObject(bucketName, filename, file, int64(file.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) CheckIfBucketExists(name string) (bool, error) {
	exists, err := fs.client.BucketExists(name)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (fs *FileStorage) CreateBucketIfDoesNotExist(name string) error {
	exists, err := fs.CheckIfBucketExists(name)
	if err != nil {
		return err
	}

	if exists {
		return nil
	} else {
		err := fs.client.MakeBucket(name, "default")
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *FileStorage) GetFile(bucketName, filename string) ([]byte, error) {
	obj, err := fs.client.GetObject(bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	byteData := []byte{}
	_, err = obj.Read(byteData)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return byteData, nil
}
