package service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
)

type FileService struct {
}

func NewFileService() *FileService {
	log.Println("new FileService created")
	return &FileService{}
}

func (s *FileService) SaveFile(fileHeader *multipart.FileHeader) error {
	//open file from request
	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	//create destination dir, mkdirAll does not error if dir exists
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create upload dir: %w", err)
	}

	//create destination file on server
	dstPath := fmt.Sprintf("%s/%s", uploadDir, fileHeader.Filename)
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
