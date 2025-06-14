package service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/emday4prez/fs-go/internal/storage"
)

type FileService struct {
	storage storage.Storage
}

func NewFileService(s storage.Storage) *FileService {
	log.Println("new FileService created with storage")
	return &FileService{
		storage: s,
	}
}

func (s *FileService) SaveFile(fileHeader *multipart.FileHeader) error {
	_, err := s.storage.Save(fileHeader)
	if err != nil {
		return fmt.Errorf("service failed to save file: %w", err)
	}
	return nil
}

func (s *FileService) ListFiles() ([]string, error) {
	//returns a sorted list of entries in a dir
	filenames, err := s.storage.List()
	if err != nil {

		return nil, fmt.Errorf("service failed to list files: %w", err)
	}

	return filenames, nil
}

func (s *FileService) GetFile(filename string) (io.ReadCloser, string, error) {
	file, path, err := s.storage.Get(filename)
	if err != nil {
		return nil, "", fmt.Errorf("service failed to get file: %w", err)
	}
	return file, path, nil
}
