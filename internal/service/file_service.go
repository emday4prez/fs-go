package service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
)

type FileService struct {
	uploadDir string
}

func NewFileService(uploadDir string) *FileService {
	log.Println("new FileService created")
	return &FileService{
		uploadDir: uploadDir,
	}
}

func (s *FileService) SaveFile(fileHeader *multipart.FileHeader) error {
	//open file from request
	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	//create destination dir, mkdirAll does not error if dir exists

	if err := os.MkdirAll(s.uploadDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create upload dir: %w", err)
	}

	//create destination file on server
	dstPath := fmt.Sprintf("%s/%s", s.uploadDir, fileHeader.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)

	}
	defer dst.Close()

	//copy from uploaded file to destination file
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file content %w", err)
	}

	log.Printf("saved file to %s", dstPath)
	return nil
}

func (s *FileService) ListFiles() ([]string, error) {

	//returns a sorted list of entries in a dir
	entries, err := os.ReadDir(s.uploadDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read uploads dir: %w", err)
	}

	var filenames []string

	for _, entry := range entries {
		if !entry.IsDir() {
			filenames = append(filenames, entry.Name())
		}
	}

	return filenames, nil
}
