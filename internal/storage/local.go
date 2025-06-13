package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	uploadDir string
}

func NewLocalStorage(uploadDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload folder: %w", err)
	}

	return &LocalStorage{
		uploadDir: uploadDir,
	}, nil
}

func (ls *LocalStorage) Save(fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dstPath := filepath.Join(ls.uploadDir, fileHeader.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	return dstPath, nil
}

func (ls *LocalStorage) List() ([]string, error) {
	entries, err := os.ReadDir(ls.uploadDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read upload directory: %w", err)
	}

	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filenames = append(filenames, entry.Name())
		}
	}
	return filenames, nil
}

func (ls *LocalStorage) Get(filename string) (io.ReadCloser, string, error) {
	// construct a safe file path to prevent directory traversal.
	filePath := filepath.Join(ls.uploadDir, filename)

	// check if the file actually exists before trying to open it.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("file not found: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}

	// return the opened file (which is an io.ReadCloser) and its path.
	return file, filePath, nil
}
