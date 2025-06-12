package storage

import (
	"io"
	"mime/multipart"
)

type Storage interface {
	Save(fileHeader *multipart.FileHeader) (string, error)

	Get(filename string) (io.ReadCloser, string, error)
}
